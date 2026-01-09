package main

import (
	"context"
	"errors"
	"exchangeapp/internal/app"
	"exchangeapp/internal/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("配置加载失败：%v", err)
	}

	s, err := app.NewServer(cfg)
	if err != nil {
		log.Fatalf("新建服务失败：%v", err)
	}

	go func() {
		if err := s.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("服务启动失败：%v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	<-ctx.Done()
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(shutdownCtx); err != nil {
		log.Printf("服务关闭失败：%v", err)
	}
}
