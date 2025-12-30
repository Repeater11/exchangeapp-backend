package app

import (
	"exchangeapp/internal/config"
	"exchangeapp/internal/db"
	"exchangeapp/internal/handler"
	"exchangeapp/internal/models"
	"exchangeapp/internal/repository"
	"exchangeapp/internal/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	e  *gin.Engine
	db *gorm.DB
}

func NewServer(cfg *config.Config) (*Server, error) {
	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery())

	gormDB, err := db.NewMySQL(&cfg.Database)
	if err != nil {
		return nil, err
	}
	if err := runMigrations(gormDB); err != nil {
		if sqlDB, err := gormDB.DB(); err == nil {
			sqlDB.Close()
		}
		return nil, fmt.Errorf("数据库迁移失败：%w", err)
	}

	userRepo := repository.NewUserRepository(gormDB)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	e.POST("/register", userHandler.Register)

	e.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return &Server{e: e, db: gormDB}, nil
}

func (s *Server) Run(addr int) error {
	return s.e.Run(fmt.Sprintf(":%d", addr))
}

func (s *Server) Close() error {
	if s.db == nil {
		return nil
	}

	sqlDB, err := s.db.DB()

	if err != nil {
		return fmt.Errorf("获取底层连接失败：%w", err)
	}

	return sqlDB.Close()
}

func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{})
}
