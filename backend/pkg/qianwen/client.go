package qianwen

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client 通义百炼智能体客户端
type Client struct {
	appID  string
	apiKey string
	apiURL string
	client *http.Client
}

// NewClient 创建通义百炼客户端
func NewClient(appID, apiKey, apiURL string) *Client {
	return &Client{
		appID:  appID,
		apiKey: apiKey,
		apiURL: strings.Replace(apiURL, "{app_id}", appID, 1),
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// ChatRequest 对话请求
type ChatRequest struct {
	SessionID string                 `json:"session_id,omitempty"`
	UserID    string                 `json:"user_id,omitempty"`
	Prompt    string                 `json:"prompt"`
	Stream    bool                   `json:"stream"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// ChatResponse SSE响应片段
type ChatResponse struct {
	StatusCode int    `json:"status_code"`
	RequestID  string `json:"request_id"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Output     struct {
		Text         string `json:"text"`
		FinishReason string `json:"finish_reason"`
	} `json:"output"`
	Usage struct {
		Models []struct {
			ModelID      string `json:"model_id"`
			InputTokens  int    `json:"input_tokens"`
			OutputTokens int    `json:"output_tokens"`
		} `json:"models"`
	} `json:"usage"`
}

// StreamCallback 流式回调函数
type StreamCallback func(text string, isDone bool) error

// ChatStream 流式对话
func (c *Client) ChatStream(req ChatRequest, callback StreamCallback) (string, error) {
	req.Stream = true
	
	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var fullText strings.Builder
	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fullText.String(), fmt.Errorf("error reading stream: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "data:") {
			continue
		}

		// 提取 data: 后的内容
		data := strings.TrimPrefix(line, "data:")
		data = strings.TrimSpace(data)

		// 检查是否为结束标记
		if data == "[DONE]" {
			if callback != nil {
				_ = callback("", true)
			}
			break
		}

		// 解析JSON
		var chatResp ChatResponse
		if err := json.Unmarshal([]byte(data), &chatResp); err != nil {
			continue
		}

		// 检查错误
		if chatResp.Code != "" && chatResp.Code != "Success" {
			return fullText.String(), fmt.Errorf("API error: %s - %s", chatResp.Code, chatResp.Message)
		}

		// 获取文本片段
		text := chatResp.Output.Text
		if text != "" {
			fullText.WriteString(text)
			if callback != nil {
				if err := callback(text, false); err != nil {
					return fullText.String(), err
				}
			}
		}

		// 检查是否完成
		if chatResp.Output.FinishReason != "" {
			if callback != nil {
				_ = callback("", true)
			}
			break
		}
	}

	return fullText.String(), nil
}
