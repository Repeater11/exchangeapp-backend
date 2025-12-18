package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name string
		Port string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
}

var AppConfig *Config

func InitConfig() error {
	useFile := false

	if path, ok := os.LookupEnv("EXCHANGEAPP_CONFIG_PATH"); ok && path != "" {
		viper.SetConfigFile(path)
		useFile = true
	} else {
		if exePath, err := os.Executable(); err == nil {
			exeDir := filepath.Dir(exePath)
			defaultPath := filepath.Join(exeDir, "configs", "config.yml")
			if _, err := os.Stat(defaultPath); err == nil {
				viper.SetConfigFile(defaultPath)
				useFile = true
			}
		}
	}

	if !useFile {
		viper.SetConfigName("config")
		viper.SetConfigType("yml")
		viper.AddConfigPath("./configs")
	}

	viper.SetEnvPrefix("EXCHANGEAPP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置失败：%w", err)
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("解码失败：%w", err)
	}

	AppConfig = cfg
	return nil
}
