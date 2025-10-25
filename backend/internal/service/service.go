package service

import (
	"ai-chat-system/internal/config"
	"ai-chat-system/internal/dao"
	"ai-chat-system/internal/model"
	"ai-chat-system/pkg/cosyvoice"
	"ai-chat-system/pkg/qianwen"
	"ai-chat-system/pkg/utils"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

// SessionService 会话服务
type SessionService struct {
	sessionDAO *dao.SessionDAO
	cfg        *config.Config
}

// NewSessionService 创建会话服务
func NewSessionService(cfg *config.Config) *SessionService {
	return &SessionService{
		sessionDAO: dao.NewSessionDAO(),
		cfg:        cfg,
	}
}

// CreateSession 创建新会话
func (s *SessionService) CreateSession(sessionID string) (*model.Session, error) {
	now := time.Now()
	session := &model.Session{
		SessionID:      sessionID,
		StartTime:      now,
		LastActiveTime: now,
		Status:         1, // 活跃
	}

	if err := s.sessionDAO.Create(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	utils.Info("Session created", zap.String("session_id", sessionID))
	return session, nil
}

// UpdateActivity 更新会话活跃时间
func (s *SessionService) UpdateActivity(sessionID string) error {
	return s.sessionDAO.UpdateLastActiveTime(sessionID)
}

// GetSession 获取会话信息
func (s *SessionService) GetSession(sessionID string) (*model.Session, error) {
	return s.sessionDAO.GetBySessionID(sessionID)
}

// ConversationService 对话服务
type ConversationService struct {
	historyDAO    *dao.ConversationHistoryDAO
	qianwenClient *qianwen.Client
	ttsClient     *cosyvoice.Client
	cfg           *config.Config
}

// NewConversationService 创建对话服务
func NewConversationService(cfg *config.Config) *ConversationService {
	qwClient := qianwen.NewClient(
		cfg.Qianwen.AppID,
		cfg.Qianwen.APIKey,
		cfg.Qianwen.APIURL,
	)

	ttsClient := cosyvoice.NewClient(
		cfg.CosyVoice.APIKey,
		cfg.CosyVoice.APIURL,
		cfg.CosyVoice.VoiceID,
		cfg.CosyVoice.SampleRate,
		cfg.CosyVoice.Volume,
		cfg.CosyVoice.SpeechRate,
		cfg.CosyVoice.PitchRate,
	)

	return &ConversationService{
		historyDAO:    dao.NewConversationHistoryDAO(),
		qianwenClient: qwClient,
		ttsClient:     ttsClient,
		cfg:           cfg,
	}
}

// SaveMessage 保存对话消息
func (s *ConversationService) SaveMessage(sessionID, role, content string, duration *int, responseTime *int) error {
	history := &model.ConversationHistory{
		SessionID:    sessionID,
		Role:         role,
		Content:      content,
		Duration:     duration,
		ResponseTime: responseTime,
	}

	if err := s.historyDAO.Create(history); err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	return nil
}

// GetHistory 获取对话历史
func (s *ConversationService) GetHistory(sessionID string) ([]model.ConversationHistory, error) {
	return s.historyDAO.GetBySessionID(sessionID)
}

// ProcessUserMessage 处理用户消息(流式)
func (s *ConversationService) ProcessUserMessage(sessionID, userID, userMessage string, streamCallback qianwen.StreamCallback) (string, error) {
	startTime := time.Now()

	// 保存用户消息
	if err := s.SaveMessage(sessionID, "user", userMessage, nil, nil); err != nil {
		utils.Error("Failed to save user message", zap.Error(err))
	}

	// 调用通义百炼API
	req := qianwen.ChatRequest{
		SessionID: sessionID,
		UserID:    userID,
		Prompt:    userMessage,
		Stream:    true,
	}

	fullText, err := s.qianwenClient.ChatStream(req, streamCallback)
	if err != nil {
		return "", fmt.Errorf("qianwen API error: %w", err)
	}

	// 计算响应时间
	responseTime := int(time.Since(startTime).Milliseconds())

	// 保存AI回复
	if err := s.SaveMessage(sessionID, "assistant", fullText, nil, &responseTime); err != nil {
		utils.Error("Failed to save assistant message", zap.Error(err))
	}

	return fullText, nil
}

// GenerateAudio 生成语音
func (s *ConversationService) GenerateAudio(text string) (string, int, error) {
	startTime := time.Now()

	audioBase64, duration, err := s.ttsClient.SynthesizeToBase64(text)
	if err != nil {
		return "", 0, fmt.Errorf("TTS error: %w", err)
	}

	utils.Info("Audio generated",
		zap.Int("duration_ms", duration),
		zap.Int("generation_time_ms", int(time.Since(startTime).Milliseconds())),
	)

	return audioBase64, duration, nil
}

// GenerateWelcomeMessage 生成欢迎语
func (s *ConversationService) GenerateWelcomeMessage(sessionID string) (string, string, int, error) {
	welcomeText := s.cfg.WelcomeMessage
	if welcomeText == "" {
		welcomeText = "欢迎使用小易助手,你有什么问题我都可以帮您。"
	}

	// 保存欢迎语到历史
	if err := s.SaveMessage(sessionID, "system", welcomeText, nil, nil); err != nil {
		utils.Error("Failed to save welcome message", zap.Error(err))
	}

	// 生成语音
	audioBase64, duration, err := s.GenerateAudio(welcomeText)
	if err != nil {
		return welcomeText, "", 0, fmt.Errorf("failed to generate welcome audio: %w", err)
	}

	return welcomeText, audioBase64, duration, nil
}

// GenerateSessionID 生成会话ID
func GenerateSessionID() string {
	return uuid.New().String()
}
