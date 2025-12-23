package main

import (
	"exchangeapp/internal/app"
	"exchangeapp/internal/config"
	"log"
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
	defer func() {
		if err := s.Close(); err != nil {
			log.Printf("关闭资源失败：%v", err)
		}
	}()

	if err := s.Run(cfg.App.Port); err != nil {
		log.Fatalf("服务启动失败：%v", err)
	}
}
