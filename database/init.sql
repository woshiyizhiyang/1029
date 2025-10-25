-- 创建数据库
CREATE DATABASE IF NOT EXISTS ai_chat_system DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE ai_chat_system;

-- 用户会话表
CREATE TABLE IF NOT EXISTS sessions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    session_id VARCHAR(64) NOT NULL UNIQUE COMMENT '会话唯一标识',
    user_identifier VARCHAR(128) DEFAULT NULL COMMENT '用户标识(可选)',
    start_time DATETIME NOT NULL COMMENT '会话开始时间',
    last_active_time DATETIME NOT NULL COMMENT '最后活跃时间',
    status TINYINT DEFAULT 1 COMMENT '会话状态:1-活跃,2-已结束',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_session_id (session_id),
    INDEX idx_status_last_active (status, last_active_time),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户会话表';

-- 对话历史表
CREATE TABLE IF NOT EXISTS conversation_history (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    session_id VARCHAR(64) NOT NULL COMMENT '关联会话ID',
    role VARCHAR(20) NOT NULL COMMENT '角色:user/assistant/system',
    content TEXT NOT NULL COMMENT '对话内容文本',
    audio_url VARCHAR(512) DEFAULT NULL COMMENT '音频文件URL(保留NULL,音频不存储)',
    duration INT DEFAULT NULL COMMENT '音频时长(毫秒,可记录用于统计)',
    ai_model VARCHAR(50) DEFAULT NULL COMMENT 'AI模型标识',
    response_time INT DEFAULT NULL COMMENT '响应耗时(毫秒)',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_session_id (session_id),
    INDEX idx_created_at (created_at),
    INDEX idx_session_created (session_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='对话历史表';

-- 系统配置表
CREATE TABLE IF NOT EXISTS system_config (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    config_key VARCHAR(100) NOT NULL UNIQUE COMMENT '配置键',
    config_value TEXT NOT NULL COMMENT '配置值(JSON格式)',
    description VARCHAR(255) DEFAULT NULL COMMENT '配置说明',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_config_key (config_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';
