package config

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

// Config 全局配置结构
type Config struct {
	Qianwen   QianwenConfig   `json:"qianwen"`
	CosyVoice CosyVoiceConfig `json:"cosy_voice"`
	Database  DatabaseConfig  `json:"database"`
	Server    ServerConfig    `json:"server"`
	WebSocket WebSocketConfig `json:"websocket"`
	Session   SessionConfig   `json:"session"`
	History   HistoryConfig   `json:"history"`
	WelcomeMessage string     `json:"welcome_message"`
}

type QianwenConfig struct {
	AppID  string `json:"app_id"`
	APIKey string `json:"api_key"`
	APIURL string `json:"api_url"`
}

type CosyVoiceConfig struct {
	VoiceID    string `json:"voice_id"`
	APIKey     string `json:"api_key"`
	APIURL     string `json:"api_url"`
	SampleRate int    `json:"sample_rate"`
	Volume     int    `json:"volume"`
	SpeechRate int    `json:"speech_rate"`
	PitchRate  int    `json:"pitch_rate"`
}

type DatabaseConfig struct {
	Host               string `json:"host"`
	Port               int    `json:"port"`
	User               string `json:"user"`
	Password           string `json:"password"`
	DBName             string `json:"dbname"`
	MaxConnections     int    `json:"max_connections"`
	MaxIdleConnections int    `json:"max_idle_connections"`
}

type ServerConfig struct {
	WSPort   int `json:"ws_port"`
	HTTPPort int `json:"http_port"`
}

type WebSocketConfig struct {
	HeartbeatInterval int `json:"heartbeat_interval"`
	Timeout           int `json:"timeout"`
	MaxReconnect      int `json:"max_reconnect"`
}

type SessionConfig struct {
	Timeout int `json:"timeout"`
}

type HistoryConfig struct {
	RetentionDays int `json:"retention_days"`
}

var (
	globalConfig *Config
	once         sync.Once
)

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	var err error
	once.Do(func() {
		viper.SetConfigFile(configPath)
		viper.SetConfigType("json")

		if readErr := viper.ReadInConfig(); readErr != nil {
			err = fmt.Errorf("failed to read config file: %w", readErr)
			return
		}

		globalConfig = &Config{}
		if unmarshalErr := viper.Unmarshal(globalConfig); unmarshalErr != nil {
			err = fmt.Errorf("failed to unmarshal config: %w", unmarshalErr)
			return
		}

		// 验证必要配置
		if err = validateConfig(globalConfig); err != nil {
			return
		}
	})

	return globalConfig, err
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	return globalConfig
}

// validateConfig 验证配置
func validateConfig(cfg *Config) error {
	if cfg.Qianwen.AppID == "" || cfg.Qianwen.APIKey == "" {
		return fmt.Errorf("qianwen config is incomplete")
	}
	if cfg.CosyVoice.APIKey == "" || cfg.CosyVoice.VoiceID == "" {
		return fmt.Errorf("cosyvoice config is incomplete")
	}
	if cfg.Database.Host == "" || cfg.Database.DBName == "" {
		return fmt.Errorf("database config is incomplete")
	}
	return nil
}
