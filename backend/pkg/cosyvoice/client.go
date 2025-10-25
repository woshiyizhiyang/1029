package cosyvoice

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client CosyVoice语音合成客户端
type Client struct {
	apiKey     string
	apiURL     string
	voiceID    string
	sampleRate int
	volume     int
	speechRate int
	pitchRate  int
	client     *http.Client
}

// NewClient 创建CosyVoice客户端
func NewClient(apiKey, apiURL, voiceID string, sampleRate, volume, speechRate, pitchRate int) *Client {
	return &Client{
		apiKey:     apiKey,
		apiURL:     apiURL,
		voiceID:    voiceID,
		sampleRate: sampleRate,
		volume:     volume,
		speechRate: speechRate,
		pitchRate:  pitchRate,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// TTSRequest 语音合成请求
type TTSRequest struct {
	Text       string `json:"text"`
	Voice      string `json:"voice"`
	Format     string `json:"format"`
	SampleRate int    `json:"sample_rate"`
	Volume     int    `json:"volume"`
	SpeechRate int    `json:"speech_rate"`
	PitchRate  int    `json:"pitch_rate"`
}

// TTSResponse 语音合成响应
type TTSResponse struct {
	RequestID string `json:"request_id"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	Data      struct {
		AudioData string `json:"audio_data"` // Base64编码的音频数据
		Duration  int    `json:"duration"`   // 音频时长(毫秒)
	} `json:"data"`
}

// Synthesize 合成语音
func (c *Client) Synthesize(text string) ([]byte, int, error) {
	// 文本预处理
	text = preprocessText(text)
	if text == "" {
		return nil, 0, fmt.Errorf("text is empty after preprocessing")
	}

	req := TTSRequest{
		Text:       text,
		Voice:      c.voiceID,
		Format:     "mp3",
		SampleRate: c.sampleRate,
		Volume:     c.volume,
		SpeechRate: c.speechRate,
		PitchRate:  c.pitchRate,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var ttsResp TTSResponse
	if err := json.Unmarshal(body, &ttsResp); err != nil {
		return nil, 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if ttsResp.Code != "" && ttsResp.Code != "Success" {
		return nil, 0, fmt.Errorf("API error: %s - %s", ttsResp.Code, ttsResp.Message)
	}

	// 解码Base64音频数据
	audioData, err := base64.StdEncoding.DecodeString(ttsResp.Data.AudioData)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to decode audio data: %w", err)
	}

	return audioData, ttsResp.Data.Duration, nil
}

// SynthesizeToBase64 合成语音并返回Base64编码
func (c *Client) SynthesizeToBase64(text string) (string, int, error) {
	audioData, duration, err := c.Synthesize(text)
	if err != nil {
		return "", 0, err
	}

	base64Audio := base64.StdEncoding.EncodeToString(audioData)
	return base64Audio, duration, nil
}

// preprocessText 文本预处理
func preprocessText(text string) string {
	// 过滤特殊字符,保留基本标点
	// 这里可以根据实际需求扩展
	return text
}
