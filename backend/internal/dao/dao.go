package dao

import (
	"ai-chat-system/internal/model"
	"ai-chat-system/pkg/utils"
	"gorm.io/gorm"
	"time"
)

// SessionDAO 会话数据访问层
type SessionDAO struct {
	db *gorm.DB
}

// NewSessionDAO 创建会话DAO实例
func NewSessionDAO() *SessionDAO {
	return &SessionDAO{
		db: utils.GetDB(),
	}
}

// Create 创建新会话
func (d *SessionDAO) Create(session *model.Session) error {
	return d.db.Create(session).Error
}

// GetBySessionID 根据会话ID查询会话
func (d *SessionDAO) GetBySessionID(sessionID string) (*model.Session, error) {
	var session model.Session
	err := d.db.Where("session_id = ?", sessionID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// UpdateLastActiveTime 更新最后活跃时间
func (d *SessionDAO) UpdateLastActiveTime(sessionID string) error {
	return d.db.Model(&model.Session{}).
		Where("session_id = ?", sessionID).
		Update("last_active_time", time.Now()).Error
}

// UpdateStatus 更新会话状态
func (d *SessionDAO) UpdateStatus(sessionID string, status int8) error {
	return d.db.Model(&model.Session{}).
		Where("session_id = ?", sessionID).
		Update("status", status).Error
}

// GetExpiredSessions 获取过期会话
func (d *SessionDAO) GetExpiredSessions(timeoutMinutes int) ([]model.Session, error) {
	var sessions []model.Session
	expireTime := time.Now().Add(-time.Duration(timeoutMinutes) * time.Minute)
	err := d.db.Where("status = ? AND last_active_time < ?", 1, expireTime).Find(&sessions).Error
	return sessions, err
}

// ConversationHistoryDAO 对话历史数据访问层
type ConversationHistoryDAO struct {
	db *gorm.DB
}

// NewConversationHistoryDAO 创建对话历史DAO实例
func NewConversationHistoryDAO() *ConversationHistoryDAO {
	return &ConversationHistoryDAO{
		db: utils.GetDB(),
	}
}

// Create 创建对话记录
func (d *ConversationHistoryDAO) Create(history *model.ConversationHistory) error {
	return d.db.Create(history).Error
}

// GetBySessionID 根据会话ID查询对话历史
func (d *ConversationHistoryDAO) GetBySessionID(sessionID string) ([]model.ConversationHistory, error) {
	var histories []model.ConversationHistory
	err := d.db.Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Find(&histories).Error
	return histories, err
}

// DeleteOldRecords 删除过期记录
func (d *ConversationHistoryDAO) DeleteOldRecords(days int) error {
	expireTime := time.Now().AddDate(0, 0, -days)
	return d.db.Where("created_at < ?", expireTime).
		Delete(&model.ConversationHistory{}).Error
}

// SystemConfigDAO 系统配置数据访问层
type SystemConfigDAO struct {
	db *gorm.DB
}

// NewSystemConfigDAO 创建系统配置DAO实例
func NewSystemConfigDAO() *SystemConfigDAO {
	return &SystemConfigDAO{
		db: utils.GetDB(),
	}
}

// GetByKey 根据配置键查询配置值
func (d *SystemConfigDAO) GetByKey(key string) (*model.SystemConfig, error) {
	var config model.SystemConfig
	err := d.db.Where("config_key = ?", key).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Upsert 更新或插入配置
func (d *SystemConfigDAO) Upsert(config *model.SystemConfig) error {
	return d.db.Save(config).Error
}
