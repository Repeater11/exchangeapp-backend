package app

import (
	"context"
	"errors"
	"exchangeapp/internal/config"
	"exchangeapp/internal/db"
	"exchangeapp/internal/handler"
	"exchangeapp/internal/middleware"
	"exchangeapp/internal/models"
	"exchangeapp/internal/repository"
	"exchangeapp/internal/service"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	e       *gin.Engine
	db      *gorm.DB
	httpSrv *http.Server
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
	userSvc := service.NewUserService(userRepo, cfg.JWT)
	userHandler := handler.NewUserHandler(userSvc)

	threadRepo := repository.NewThreadRepository(gormDB)
	threadLikeRepo := repository.NewThreadLikeRepository(gormDB)
	threadSvc := service.NewThreadService(threadRepo, threadLikeRepo)
	threadLikeSvc := service.NewThreadLikeService(threadRepo, threadLikeRepo)
	threadHandler := handler.NewThreadHandler(threadSvc)
	threadLikeHandler := handler.NewThreadLikeHandler(threadLikeSvc)

	replyRepo := repository.NewReplyRepository(gormDB)
	replySvc := service.NewReplyService(replyRepo, threadRepo)
	replyHandler := handler.NewReplyHandler(replySvc)

	e.POST("/register", userHandler.Register)
	e.POST("/login", userHandler.Login)
	e.GET("/threads", threadHandler.List)
	e.GET("/threads/:id/replies", replyHandler.ListByThreadID)
	e.GET("/threads/:id", threadHandler.Detail)

	authGroup := e.Group("/api")
	authGroup.Use(middleware.Auth(cfg.JWT.Secret))
	authGroup.GET("/me", userHandler.Me)
	authGroup.GET("/me/threads", threadHandler.ListMine)
	authGroup.GET("/me/replies", replyHandler.ListMine)
	authGroup.POST("/threads", threadHandler.Create)
	authGroup.POST("/threads/:id/replies", replyHandler.Create)
	authGroup.PUT("/threads/:id", threadHandler.Update)
	authGroup.DELETE("/threads/:id", threadHandler.Delete)
	authGroup.PUT("/replies/:id", replyHandler.Update)
	authGroup.DELETE("/replies/:id", replyHandler.Delete)
	authGroup.POST("/threads/:id/like", threadLikeHandler.Like)
	authGroup.DELETE("/threads/:id/like", threadLikeHandler.Unlike)

	e.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	httpSrv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      e,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		e:       e,
		db:      gormDB,
		httpSrv: httpSrv,
	}, nil
}

func (s *Server) Run() error {
	return s.httpSrv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	var httpErr error
	if s.httpSrv != nil {
		if err := s.httpSrv.Shutdown(ctx); err != nil {
			httpErr = fmt.Errorf("关闭 HTTP 服务失败：%w", err)
		}
	}

	dbErr := s.Close()
	if dbErr != nil {
		dbErr = fmt.Errorf("关闭数据库失败：%w", dbErr)
	}

	return errors.Join(httpErr, dbErr)
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
	return db.AutoMigrate(&models.User{}, &models.Thread{}, &models.Reply{}, &models.ThreadLike{})
}
