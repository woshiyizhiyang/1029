package model

import "time"

// Session 用户会话模型
type Session struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	SessionID      string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"session_id"`
	UserIdentifier string    `gorm:"type:varchar(128);index" json:"user_identifier"`
	StartTime      time.Time `gorm:"not null" json:"start_time"`
	LastActiveTime time.Time `gorm:"not null;index:idx_status_last_active" json:"last_active_time"`
	Status         int8      `gorm:"type:tinyint;default:1;index:idx_status_last_active" json:"status"` // 1-活跃,2-已结束
	CreatedAt      time.Time `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (Session) TableName() string {
	return "sessions"
}

// ConversationHistory 对话历史模型
type ConversationHistory struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	SessionID    string    `gorm:"type:varchar(64);index;not null" json:"session_id"`
	Role         string    `gorm:"type:varchar(20);not null" json:"role"` // user/assistant/system
	Content      string    `gorm:"type:text;not null" json:"content"`
	AudioURL     *string   `gorm:"type:varchar(512)" json:"audio_url"`
	Duration     *int      `gorm:"type:int" json:"duration"`           // 音频时长(毫秒)
	AIModel      *string   `gorm:"type:varchar(50)" json:"ai_model"`   // AI模型标识
	ResponseTime *int      `gorm:"type:int" json:"response_time"`      // 响应耗时(毫秒)
	CreatedAt    time.Time `gorm:"autoCreateTime;index:idx_session_created" json:"created_at"`
}

// TableName 指定表名
func (ConversationHistory) TableName() string {
	return "conversation_history"
}

// SystemConfig 系统配置模型
type SystemConfig struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ConfigKey   string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"config_key"`
	ConfigValue string    `gorm:"type:text;not null" json:"config_value"`
	Description *string   `gorm:"type:varchar(255)" json:"description"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (SystemConfig) TableName() string {
	return "system_config"
}
