# 部署指南

## 生产环境部署

### 1. 服务器准备

**系统要求:**
- Linux (Ubuntu 20.04 / CentOS 7+)
- CPU: 2核心
- 内存: 4GB
- 磁盘: 50GB SSD

**安装必要软件:**

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y mysql-server nginx

# CentOS/RHEL
sudo yum install -y mysql-server nginx
```

### 2. 数据库部署

```bash
# 启动MySQL
sudo systemctl start mysql
sudo systemctl enable mysql

# 安全配置
sudo mysql_secure_installation

# 创建数据库
mysql -u root -p
```

执行SQL:
```sql
CREATE DATABASE ai_chat_system DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'aichat'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON ai_chat_system.* TO 'aichat'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

导入数据库:
```bash
mysql -u aichat -p ai_chat_system < database/init.sql
```

### 3. 后端部署

**编译程序:**
```bash
cd backend
go build -o ai-chat-server cmd/main.go
```

**配置文件:**
```bash
sudo mkdir -p /etc/ai-chat-system
sudo cp ../config/app.json /etc/ai-chat-system/
sudo chown -R www-data:www-data /etc/ai-chat-system
```

**创建systemd服务:**

创建文件 `/etc/systemd/system/ai-chat.service`:

```ini
[Unit]
Description=AI Chat System Backend
After=network.target mysql.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/var/www/ai-chat-backend
ExecStart=/var/www/ai-chat-backend/ai-chat-server \
    -config /etc/ai-chat-system/app.json \
    -log /var/log/ai-chat/app.log
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

**部署后端:**
```bash
sudo mkdir -p /var/www/ai-chat-backend
sudo mkdir -p /var/log/ai-chat
sudo cp ai-chat-server /var/www/ai-chat-backend/
sudo chown -R www-data:www-data /var/www/ai-chat-backend
sudo chown -R www-data:www-data /var/log/ai-chat

# 启动服务
sudo systemctl daemon-reload
sudo systemctl start ai-chat
sudo systemctl enable ai-chat

# 查看状态
sudo systemctl status ai-chat
```

### 4. 前端部署

**构建前端:**
```bash
cd frontend
npm install
npm run build
```

**部署到Nginx:**
```bash
sudo mkdir -p /var/www/ai-chat-frontend
sudo cp -r dist/* /var/www/ai-chat-frontend/
sudo chown -R www-data:www-data /var/www/ai-chat-frontend
```

### 5. Nginx配置

创建文件 `/etc/nginx/sites-available/ai-chat`:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    location / {
        root /var/www/ai-chat-frontend;
        index index.html;
        try_files $uri $uri/ /index.html;
        
        # 缓存控制
        expires 1d;
        add_header Cache-Control "public, immutable";
    }

    # WebSocket代理
    location /ws {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # 健康检查
    location /health {
        proxy_pass http://localhost:8080;
        access_log off;
    }

    # 静态资源缓存
    location ~* \.(jpg|jpeg|png|gif|ico|css|js|svg|woff|woff2|ttf|eot)$ {
        root /var/www/ai-chat-frontend;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
}
```

启用站点:
```bash
sudo ln -s /etc/nginx/sites-available/ai-chat /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 6. SSL证书配置(可选)

使用Let's Encrypt:

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

### 7. 防火墙配置

```bash
# Ubuntu UFW
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# CentOS Firewalld
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload
```

## 维护操作

### 日志管理

**配置logrotate:**

创建文件 `/etc/logrotate.d/ai-chat`:

```
/var/log/ai-chat/*.log {
    daily
    rotate 30
    compress
    delaycompress
    notifempty
    create 0640 www-data www-data
    sharedscripts
    postrotate
        systemctl reload ai-chat > /dev/null 2>&1 || true
    endscript
}
```

### 数据清理

创建清理脚本 `/usr/local/bin/ai-chat-cleanup.sh`:

```bash
#!/bin/bash
# 清理30天前的对话历史

mysql -u aichat -p'your_password' ai_chat_system <<EOF
DELETE FROM conversation_history WHERE created_at < DATE_SUB(NOW(), INTERVAL 30 DAY);
UPDATE sessions SET status = 2 WHERE status = 1 AND last_active_time < DATE_SUB(NOW(), INTERVAL 30 MINUTE);
EOF

echo "Cleanup completed at $(date)"
```

添加到crontab:
```bash
sudo chmod +x /usr/local/bin/ai-chat-cleanup.sh
sudo crontab -e
```

添加:
```
0 2 * * * /usr/local/bin/ai-chat-cleanup.sh >> /var/log/ai-chat/cleanup.log 2>&1
```

### 备份策略

**数据库备份脚本:**

```bash
#!/bin/bash
BACKUP_DIR="/var/backups/ai-chat"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR

# 备份数据库
mysqldump -u aichat -p'your_password' ai_chat_system > $BACKUP_DIR/db_$DATE.sql

# 压缩
gzip $BACKUP_DIR/db_$DATE.sql

# 删除7天前的备份
find $BACKUP_DIR -name "db_*.sql.gz" -mtime +7 -delete

echo "Backup completed: db_$DATE.sql.gz"
```

### 监控配置

**服务状态监控:**

```bash
# 检查服务状态
sudo systemctl status ai-chat
sudo systemctl status nginx
sudo systemctl status mysql

# 查看日志
sudo journalctl -u ai-chat -f
sudo tail -f /var/log/ai-chat/app.log
sudo tail -f /var/log/nginx/access.log
```

### 更新部署

**后端更新:**
```bash
# 编译新版本
cd backend
go build -o ai-chat-server cmd/main.go

# 停止服务
sudo systemctl stop ai-chat

# 备份旧版本
sudo cp /var/www/ai-chat-backend/ai-chat-server /var/www/ai-chat-backend/ai-chat-server.bak

# 部署新版本
sudo cp ai-chat-server /var/www/ai-chat-backend/

# 启动服务
sudo systemctl start ai-chat
```

**前端更新:**
```bash
# 构建新版本
cd frontend
npm run build

# 备份旧版本
sudo cp -r /var/www/ai-chat-frontend /var/www/ai-chat-frontend.bak

# 部署新版本
sudo rm -rf /var/www/ai-chat-frontend/*
sudo cp -r dist/* /var/www/ai-chat-frontend/
sudo chown -R www-data:www-data /var/www/ai-chat-frontend

# 清理Nginx缓存(如有)
sudo systemctl reload nginx
```

## 故障排查

### 后端服务无法启动

```bash
# 查看服务状态
sudo systemctl status ai-chat

# 查看日志
sudo journalctl -u ai-chat -n 100

# 检查配置文件
cat /etc/ai-chat-system/app.json

# 检查端口占用
sudo netstat -tlnp | grep 8080
```

### WebSocket连接失败

```bash
# 测试WebSocket连接
wscat -c ws://localhost:8080/ws

# 检查Nginx配置
sudo nginx -t

# 查看Nginx错误日志
sudo tail -f /var/log/nginx/error.log
```

### 数据库连接失败

```bash
# 测试数据库连接
mysql -u aichat -p ai_chat_system

# 检查MySQL状态
sudo systemctl status mysql

# 查看MySQL错误日志
sudo tail -f /var/log/mysql/error.log
```

## 性能优化

### 数据库优化

```sql
-- 分析表
ANALYZE TABLE sessions;
ANALYZE TABLE conversation_history;

-- 查看慢查询
SHOW VARIABLES LIKE 'slow_query%';
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 2;
```

### Nginx优化

添加到nginx配置:
```nginx
# Gzip压缩
gzip on;
gzip_vary on;
gzip_min_length 1024;
gzip_types text/plain text/css text/xml text/javascript application/json application/javascript application/xml+rss;

# 连接优化
keepalive_timeout 65;
client_max_body_size 10M;
```

### 系统优化

```bash
# 增加文件描述符限制
echo "* soft nofile 65535" | sudo tee -a /etc/security/limits.conf
echo "* hard nofile 65535" | sudo tee -a /etc/security/limits.conf
```

---

以上为完整的生产环境部署指南。
