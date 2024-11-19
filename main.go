package main

import (
	"context"
	"go-demo/config"
	"go-demo/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化路由
	r := router.InitRouter()

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    config.AppConfig.App.Port,
		Handler: r,
	}

	// 在独立的 goroutine 中启动服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("启动服务器失败: %v\n", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	// 监听 SIGINT (Ctrl+C) 和 SIGTERM 信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	// 设置 5 秒的超时时间来处理关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 在这里添加清理资源的代码
	// 例如：关闭数据库连接
	// if err := db.Close(); err != nil {
	//     log.Printf("关闭数据库连接出错: %v\n", err)
	// }

	// 优雅地关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("服务器强制关闭:", err)
	}

	log.Println("服务器已优雅退出")
}
