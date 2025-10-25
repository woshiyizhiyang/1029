package main

import (
	"ai-chat-system/internal/config"
	"ai-chat-system/internal/handler"
	ws "ai-chat-system/internal/websocket"
	"ai-chat-system/pkg/utils"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config/app.json", "配置文件路径")
	logPath := flag.String("log", "logs/app.log", "日志文件路径")
	flag.Parse()

	// 初始化日志
	if err := utils.InitLogger(*logPath); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer utils.GetLogger().Sync()

	utils.Info("Starting AI Chat System...")

	// 加载配置
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		utils.Fatal("Failed to load config", zap.Error(err))
	}

	utils.Info("Config loaded successfully")

	// 初始化数据库
	if err := utils.InitDB(&cfg.Database); err != nil {
		utils.Fatal("Failed to initialize database", zap.Error(err))
	}

	utils.Info("Database initialized successfully")

	// 创建WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()

	// 创建WebSocket处理器
	wsHandler := handler.NewWSHandler(hub, cfg)

	// 注册路由
	http.HandleFunc("/ws", wsHandler.HandleWebSocket)

	// 健康检查接口
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 启动HTTP服务器
	addr := fmt.Sprintf(":%d", cfg.Server.WSPort)
	utils.Info("Server starting", zap.String("address", addr))

	// 优雅关闭
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			utils.Fatal("Server error", zap.Error(err))
		}
	}()

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	utils.Info("Shutting down server...")
}
