package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

// Message WebSocket消息结构
type Message struct {
	Type      string `json:"type"`       // 消息类型
	SessionID string `json:"session_id"` // 会话ID
	Content   string `json:"content"`    // 消息内容
	Timestamp int64  `json:"timestamp"`  // 时间戳
	
	// 额外字段(用于服务端推送)
	AudioBase64 string `json:"audioBase64,omitempty"` // 音频Base64
	Duration    int    `json:"duration,omitempty"`    // 音频时长
	FullText    string `json:"fullText,omitempty"`    // 完整文本
	ErrorCode   string `json:"errorCode,omitempty"`   // 错误码
	ErrorMsg    string `json:"errorMsg,omitempty"`    // 错误信息
}

// Client WebSocket客户端连接
type Client struct {
	ID         string
	SessionID  string
	Conn       *websocket.Conn
	Send       chan *Message
	Hub        *Hub
	mu         sync.Mutex
	isStopped  bool // 标记是否收到停止指令
}

// Hub WebSocket连接管理器
type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// NewHub 创建Hub实例
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run 运行Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.Send)
			}
			h.mu.Unlock()
		}
	}
}

// SendToClient 发送消息给指定客户端
func (h *Hub) SendToClient(clientID string, message *Message) error {
	h.mu.RLock()
	client, ok := h.clients[clientID]
	h.mu.RUnlock()

	if !ok {
		return fmt.Errorf("client not found: %s", clientID)
	}

	select {
	case client.Send <- message:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("send timeout")
	}
}

// GetClient 获取客户端
func (h *Hub) GetClient(clientID string) (*Client, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	client, ok := h.clients[clientID]
	return client, ok
}

// ReadPump 读取客户端消息
func (c *Client) ReadPump(handler func(*Client, *Message)) {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, messageData, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("WebSocket error: %v\n", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(messageData, &msg); err != nil {
			fmt.Printf("Failed to unmarshal message: %v\n", err)
			continue
		}

		// 处理停止指令
		if msg.Type == "stop" {
			c.mu.Lock()
			c.isStopped = true
			c.mu.Unlock()
			
			// 发送停止确认
			c.Send <- &Message{
				Type: "stop_ack",
			}
			continue
		}

		// 心跳响应
		if msg.Type == "heartbeat" {
			c.Send <- &Message{
				Type: "heartbeat_ack",
			}
			continue
		}

		// 调用处理器
		if handler != nil {
			handler(c, &msg)
		}
	}
}

// WritePump 写入消息到客户端
func (c *Client) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				fmt.Printf("Failed to marshal message: %v\n", err)
				continue
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// IsStopped 检查是否已停止
func (c *Client) IsStopped() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.isStopped
}

// ResetStopFlag 重置停止标记
func (c *Client) ResetStopFlag() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.isStopped = false
}
