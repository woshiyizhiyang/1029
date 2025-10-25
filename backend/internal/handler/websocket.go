package handler

import (
	"ai-chat-system/internal/config"
	"ai-chat-system/internal/service"
	ws "ai-chat-system/internal/websocket"
	"ai-chat-system/pkg/utils"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有跨域请求
	},
}

// WSHandler WebSocket处理器
type WSHandler struct {
	hub                 *ws.Hub
	sessionService      *service.SessionService
	conversationService *service.ConversationService
	cfg                 *config.Config
}

// NewWSHandler 创建WebSocket处理器
func NewWSHandler(hub *ws.Hub, cfg *config.Config) *WSHandler {
	return &WSHandler{
		hub:                 hub,
		sessionService:      service.NewSessionService(cfg),
		conversationService: service.NewConversationService(cfg),
		cfg:                 cfg,
	}
}

// HandleWebSocket 处理WebSocket连接
func (h *WSHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.Error("Failed to upgrade connection", zap.Error(err))
		return
	}

	// 生成客户端ID
	clientID := service.GenerateSessionID()
	
	client := &ws.Client{
		ID:   clientID,
		Conn: conn,
		Send: make(chan *ws.Message, 256),
		Hub:  h.hub,
	}

	h.hub.register <- client

	utils.Info("WebSocket client connected", zap.String("client_id", clientID))

	// 启动读写协程
	go client.WritePump()
	go client.ReadPump(h.handleMessage)
}

// handleMessage 处理WebSocket消息
func (h *WSHandler) handleMessage(client *ws.Client, msg *ws.Message) {
	utils.Info("Received message",
		zap.String("client_id", client.ID),
		zap.String("type", msg.Type),
		zap.String("session_id", msg.SessionID),
	)

	switch msg.Type {
	case "user_message":
		h.handleUserMessage(client, msg)
	case "get_welcome":
		h.handleWelcome(client, msg)
	case "get_history":
		h.handleGetHistory(client, msg)
	default:
		utils.Warn("Unknown message type", zap.String("type", msg.Type))
	}
}

// handleWelcome 处理欢迎语请求
func (h *WSHandler) handleWelcome(client *ws.Client, msg *ws.Message) {
	sessionID := msg.SessionID
	if sessionID == "" {
		sessionID = service.GenerateSessionID()
	}

	client.SessionID = sessionID

	// 创建会话
	_, err := h.sessionService.CreateSession(sessionID)
	if err != nil {
		utils.Error("Failed to create session", zap.Error(err))
		h.sendError(client, "SESSION_ERROR", "创建会话失败")
		return
	}

	// 生成欢迎语
	welcomeText, audioBase64, duration, err := h.conversationService.GenerateWelcomeMessage(sessionID)
	if err != nil {
		utils.Error("Failed to generate welcome message", zap.Error(err))
		h.sendError(client, "WELCOME_ERROR", "生成欢迎语失败")
		return
	}

	// 发送欢迎语
	client.Send <- &ws.Message{
		Type:        "welcome",
		SessionID:   sessionID,
		Content:     welcomeText,
		AudioBase64: audioBase64,
		Duration:    duration,
	}

	utils.Info("Welcome message sent", zap.String("session_id", sessionID))
}

// handleUserMessage 处理用户消息
func (h *WSHandler) handleUserMessage(client *ws.Client, msg *ws.Message) {
	// 重置停止标记
	client.ResetStopFlag()

	sessionID := msg.SessionID
	userMessage := strings.TrimSpace(msg.Content)

	if sessionID == "" || userMessage == "" {
		h.sendError(client, "INVALID_INPUT", "会话ID或消息内容为空")
		return
	}

	// 更新会话活跃时间
	if err := h.sessionService.UpdateActivity(sessionID); err != nil {
		utils.Error("Failed to update session activity", zap.Error(err))
	}

	// 流式处理AI回复
	var fullText strings.Builder
	streamCallback := func(text string, isDone bool) error {
		// 检查是否收到停止指令
		if client.IsStopped() {
			return fmt.Errorf("stopped by user")
		}

		if isDone {
			// 发送完成标记
			client.Send <- &ws.Message{
				Type:     "ai_text_complete",
				FullText: fullText.String(),
			}
			return nil
		}

		if text != "" {
			fullText.WriteString(text)
			// 推送文本片段
			client.Send <- &ws.Message{
				Type:    "ai_text_chunk",
				Content: text,
			}
		}
		return nil
	}

	// 调用对话服务
	aiResponse, err := h.conversationService.ProcessUserMessage(sessionID, client.ID, userMessage, streamCallback)
	if err != nil {
		// 检查是否是用户停止
		if client.IsStopped() {
			utils.Info("Conversation stopped by user", zap.String("session_id", sessionID))
			return
		}

		utils.Error("Failed to process user message", zap.Error(err))
		h.sendError(client, "AI_ERROR", "AI服务异常,请稍后再试")
		return
	}

	// 检查是否在AI回复过程中被停止
	if client.IsStopped() {
		utils.Info("Audio generation skipped due to stop", zap.String("session_id", sessionID))
		return
	}

	// 生成语音
	audioBase64, duration, err := h.conversationService.GenerateAudio(aiResponse)
	if err != nil {
		utils.Error("Failed to generate audio", zap.Error(err))
		// 降级:仅显示文字
		h.sendError(client, "TTS_ERROR", "语音合成失败,已显示文字")
		return
	}

	// 检查是否在语音合成过程中被停止
	if client.IsStopped() {
		utils.Info("Audio sending skipped due to stop", zap.String("session_id", sessionID))
		return
	}

	// 发送音频
	client.Send <- &ws.Message{
		Type:        "audio_data",
		AudioBase64: audioBase64,
		Duration:    duration,
	}

	utils.Info("Audio sent", zap.String("session_id", sessionID), zap.Int("duration", duration))
}

// handleGetHistory 处理获取历史记录请求
func (h *WSHandler) handleGetHistory(client *ws.Client, msg *ws.Message) {
	sessionID := msg.SessionID
	if sessionID == "" {
		h.sendError(client, "INVALID_INPUT", "会话ID为空")
		return
	}

	history, err := h.conversationService.GetHistory(sessionID)
	if err != nil {
		utils.Error("Failed to get history", zap.Error(err))
		h.sendError(client, "HISTORY_ERROR", "获取历史记录失败")
		return
	}

	// 发送历史记录
	client.Send <- &ws.Message{
		Type:    "history_data",
		Content: mustMarshalJSON(history),
	}
}

// sendError 发送错误消息
func (h *WSHandler) sendError(client *ws.Client, code, msg string) {
	client.Send <- &ws.Message{
		Type:      "error",
		ErrorCode: code,
		ErrorMsg:  msg,
	}
}

// mustMarshalJSON JSON序列化
func mustMarshalJSON(v interface{}) string {
	data, _ := utils.GetLogger().Sugar().Desugar().Core().Enabled(zap.DebugLevel)
	if data {
		return fmt.Sprintf("%+v", v)
	}
	return ""
}
