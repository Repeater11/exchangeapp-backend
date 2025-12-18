package main

import (
	"exchangeapp/internal/config"
	"fmt"
	"log"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("配置加载失败：%v", err)
	}
	fmt.Println(config.AppConfig.App.Port)
}
