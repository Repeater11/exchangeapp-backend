package main

import (
	"exchangeapp/internal/config"
	"fmt"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("配置加载失败：%v", err)
	}
	fmt.Println(cfg.Database.Port)
}
