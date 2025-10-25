# 动画形象AI对话系统

基于Vue3 + Golang的智能对话系统，集成通义百炼智能体和阿里CosyVoice语音合成，实现动画形象口型同步对话体验。

## 功能特性

- ✅ AI智能对话(基于通义百炼智能体)
- ✅ 实时语音合成(基于阿里CosyVoice)
- ✅ 动画口型同步(5帧动画状态切换)
- ✅ WebSocket实时通信
- ✅ 流式文本输出
- ✅ 对话历史查看
- ✅ 语音播放控制(停止功能)
- ✅ 响应式H5界面
- ✅ 欢迎语自动播放

## 技术栈

### 后端
- Golang 1.21+
- Gorilla WebSocket
- GORM (MySQL ORM)
- Viper (配置管理)
- Zap (日志系统)

### 前端
- Vue 3
- Pinia (状态管理)
- Vite (构建工具)
- Web Audio API (音频处理)
- Font Awesome (图标)

### 数据库
- MySQL 5.7+

### 第三方服务
- 通义百炼智能体API
- 阿里智能语音CosyVoice API

## 项目结构

```
.
├── backend/                    # 后端代码
│   ├── cmd/                    # 主程序入口
│   │   └── main.go
│   ├── internal/               # 内部代码
│   │   ├── config/             # 配置管理
│   │   ├── dao/                # 数据访问层
│   │   ├── handler/            # 处理器
│   │   ├── model/              # 数据模型
│   │   ├── service/            # 业务服务
│   │   └── websocket/          # WebSocket服务
│   ├── pkg/                    # 公共包
│   │   ├── cosyvoice/          # CosyVoice客户端
│   │   ├── qianwen/            # 通义百炼客户端
│   │   └── utils/              # 工具类
│   └── go.mod                  # Go模块定义
├── frontend/                   # 前端代码
│   ├── src/
│   │   ├── api/                # API接口
│   │   ├── assets/             # 静态资源
│   │   ├── components/         # Vue组件
│   │   ├── store/              # 状态管理
│   │   ├── utils/              # 工具类
│   │   ├── App.vue             # 根组件
│   │   └── main.js             # 入口文件
│   ├── public/                 # 公共资源
│   │   └── images/             # 动画图片(需放置5张口型图片)
│   ├── index.html
│   ├── package.json
│   └── vite.config.js
├── config/                     # 配置文件
│   └── app.example.json        # 配置示例
├── database/                   # 数据库脚本
│   └── init.sql                # 初始化脚本
└── README.md
```

## 快速开始

### 1. 环境准备

**必需软件:**
- Go 1.21+
- Node.js 16+
- MySQL 5.7+

**第三方服务:**
- 通义百炼智能体账号和API密钥
- 阿里智能语音CosyVoice账号和API密钥

### 2. 数据库初始化

```bash
# 登录MySQL
mysql -u root -p

# 执行初始化脚本
source database/init.sql
```

### 3. 后端配置

```bash
cd backend

# 复制配置文件
cp ../config/app.example.json ../config/app.json

# 编辑配置文件,填入真实的API密钥和数据库信息
# vim ../config/app.json
```

**配置文件说明 (config/app.json):**

```json
{
  "qianwen": {
    "app_id": "你的智能体应用ID",
    "api_key": "你的通义百炼API密钥",
    "api_url": "https://dashscope.aliyuncs.com/api/v1/apps/{app_id}/completion"
  },
  "cosy_voice": {
    "voice_id": "longxiaochun",
    "api_key": "你的CosyVoice API密钥",
    "api_url": "https://nls-gateway.cn-shanghai.aliyuncs.com/stream/v1/tts",
    "sample_rate": 24000,
    "volume": 50,
    "speech_rate": 0,
    "pitch_rate": 0
  },
  "database": {
    "host": "127.0.0.1",
    "port": 3306,
    "user": "root",
    "password": "你的数据库密码",
    "dbname": "ai_chat_system",
    "max_connections": 100,
    "max_idle_connections": 10
  }
}
```

### 4. 后端启动

```bash
cd backend

# 安装依赖
go mod tidy

# 运行服务
go run cmd/main.go -config ../config/app.json -log logs/app.log
```

服务将在以下端口启动:
- WebSocket服务: `ws://localhost:8080/ws`
- 健康检查: `http://localhost:8080/health`

### 5. 前端配置

**准备动画图片:**

在 `frontend/public/images/` 目录下放置5张口型图片:
- `mouth-0.jpg` - 闭口
- `mouth-1.jpg` - 微张
- `mouth-2.jpg` - 半张
- `mouth-3.jpg` - 大张
- `mouth-4.jpg` - 完全张口

**图片规格要求:**
- 尺寸: 720×1440 (竖版)
- 格式: JPG/PNG (推荐WebP格式以优化体积)
- 大小: 建议每张<500KB

### 6. 前端启动

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

访问: `http://localhost:3000`

### 7. 生产构建

**前端构建:**
```bash
cd frontend
npm run build
# 构建产物在 dist/ 目录
```

**后端编译:**
```bash
cd backend
go build -o ai-chat-server cmd/main.go
```

## 使用说明

### 基本使用流程

1. **打开页面** - 自动播放欢迎语"欢迎使用小易助手,你有什么问题我都可以帮您。"
2. **输入问题** - 在底部输入框输入文字(最多10000字)
3. **发送消息** - 点击发送按钮(飞机图标)
4. **查看回复** - AI流式返回文字,并自动播放语音,动画形象同步口型
5. **停止播放** - 播放期间可点击停止按钮(停止图标)中断
6. **查看历史** - 点击右上角历史按钮查看对话记录

### 功能说明

**字数统计:**
- 实时显示输入字数(X/10000)
- 超过9000字时变红色警告
- 超过10000字禁止发送

**停止功能:**
- 发送消息后,发送按钮自动变为停止按钮
- 点击停止立即中断AI回复和语音播放
- 停止后可立即发送新消息

**历史对话:**
- 有对话后右上角显示历史按钮
- 点击查看当前会话的所有对话
- 历史记录仅显示文字,不重新播放语音

**会话管理:**
- 会话ID自动保存到浏览器LocalStorage
- 刷新页面不恢复对话内容,但保留会话ID
- 每次打开页面重新播放欢迎语

## API接口说明

### WebSocket消息协议

**客户端发送:**

```json
{
  "type": "user_message",      // 消息类型
  "session_id": "xxx",          // 会话ID
  "content": "用户输入的问题",  // 消息内容
  "timestamp": 1234567890       // 时间戳
}
```

**服务端推送:**

| 消息类型 | 说明 | 字段 |
|---------|------|------|
| welcome | 欢迎语 | content, audioBase64, duration |
| ai_text_chunk | AI文本片段 | content |
| ai_text_complete | AI文本完成 | fullText |
| audio_data | 语音数据 | audioBase64, duration |
| error | 错误信息 | errorCode, errorMsg |
| stop_ack | 停止确认 | - |
| heartbeat_ack | 心跳响应 | - |

## 常见问题

### 1. WebSocket连接失败
- 检查后端服务是否启动
- 确认端口8080未被占用
- 检查防火墙设置

### 2. AI回复失败
- 确认通义百炼API密钥配置正确
- 检查app_id是否有效
- 查看后端日志排查错误

### 3. 语音合成失败
- 确认CosyVoice API密钥配置正确
- 检查voice_id是否正确
- 查看后端日志排查错误

### 4. 动画不显示
- 确认5张口型图片已放置在`frontend/public/images/`目录
- 检查图片命名是否正确(mouth-0.jpg ~ mouth-4.jpg)
- 查看浏览器控制台是否有图片加载错误

### 5. 音频无法播放
- 检查浏览器是否允许自动播放音频
- 尝试用户交互后再播放
- 确认音频格式为MP3

## 性能优化建议

1. **图片优化** - 使用WebP格式压缩图片体积
2. **音频缓存** - 可选:对相同文本的音频进行Redis缓存
3. **CDN加速** - 生产环境使用CDN加速静态资源
4. **数据库索引** - 已创建必要索引,无需额外优化
5. **连接池** - 数据库连接池已配置(最大100连接)

## 部署说明

### 单机部署

**推荐配置:**
- CPU: 2核
- 内存: 4GB
- 磁盘: 50GB SSD
- 系统: Linux (Ubuntu 20.04 / CentOS 7+)

**部署步骤:**

1. 安装MySQL并导入数据库
2. 编译后端程序并配置systemd服务
3. 构建前端并配置Nginx反向代理
4. 配置域名和SSL证书(可选)

### Nginx配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    location / {
        root /var/www/ai-chat-frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    # WebSocket代理
    location /ws {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # 健康检查
    location /health {
        proxy_pass http://localhost:8080;
    }
}
```

## 维护说明

### 日志管理
- 日志路径: `backend/logs/app.log`
- 日志轮转: 建议配置logrotate
- 日志级别: 可在代码中调整(INFO/DEBUG/WARN/ERROR)

### 数据清理
- 会话过期: 30分钟无活动自动标记过期
- 历史清理: 30天前的对话记录自动删除
- 定时任务: 建议配置cron执行清理脚本

### 监控指标
- WebSocket连接数
- API响应时间
- 错误率统计
- 数据库性能

## 开发说明

### 后端开发

**添加新的API接口:**
1. 在`internal/handler/`创建处理器
2. 在`cmd/main.go`注册路由
3. 在`internal/service/`实现业务逻辑

**修改配置项:**
1. 更新`internal/config/config.go`结构体
2. 更新`config/app.example.json`示例配置
3. 更新README配置说明

### 前端开发

**添加新组件:**
1. 在`src/components/`创建Vue组件
2. 在需要的地方引入使用
3. 更新状态管理(如需要)

**修改样式:**
- 全局样式: `src/assets/styles/global.css`
- 组件样式: 各组件的`<style scoped>`

## 许可证

本项目仅供学习和研究使用。

## 联系方式

如有问题,请通过以下方式联系:
- 提交Issue
- 发送邮件

---

**注意:** 使用前请确保已获取通义百炼和CosyVoice的API使用权限,并遵守相关服务条款。
