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
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Server struct {
	e       *gin.Engine
	db      *gorm.DB
	rdb     *redis.Client
	httpSrv *http.Server

	workerCancel context.CancelFunc
	workerDone   chan struct{}
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

	rdb, err := db.NewRedis(&cfg.Redis)
	if err != nil {
		if sqlDB, dbErr := gormDB.DB(); dbErr == nil {
			sqlDB.Close()
		}
		return nil, err
	}

	userRepo := repository.NewUserRepository(gormDB)
	userSvc := service.NewUserService(userRepo, cfg.JWT)
	userHandler := handler.NewUserHandler(userSvc)

	redisCounter := repository.NewRedisLikeCounter(rdb)

	dbthreadRepo := repository.NewThreadRepository(gormDB)
	threadRepo := repository.NewCachedThreadRepository(dbthreadRepo, rdb)
	threadLikeRepo := repository.NewThreadLikeRepository(gormDB)
	likeCounter := repository.NewCachedThreadLikeCounter(threadRepo, redisCounter)
	threadSvc := service.NewThreadService(threadRepo, threadLikeRepo, likeCounter)
	threadLikeSvc := service.NewThreadLikeService(threadRepo, threadLikeRepo, likeCounter)
	threadHandler := handler.NewThreadHandler(threadSvc)
	threadLikeHandler := handler.NewThreadLikeHandler(threadLikeSvc)

	replyRepo := repository.NewReplyRepository(gormDB)
	replySvc := service.NewReplyService(replyRepo, threadRepo)
	replyHandler := handler.NewReplyHandler(replySvc)

	writer, ok := dbthreadRepo.(repository.ThreadLikeCountWriter)
	if !ok {
		return nil, fmt.Errorf("线程仓库不支持 SetLikeCount")
	}
	batch := cfg.LikeWorker.Batch
	interval := time.Duration(cfg.LikeWorker.IntervalSeconds) * time.Second
	flusher := NewLikeCountFlusher(redisCounter, writer, batch, interval)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		flusher.Run(ctx)
	}()

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
	authGroup.GET("/threads/:id/like", threadLikeHandler.Status)

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
		e:            e,
		db:           gormDB,
		rdb:          rdb,
		httpSrv:      httpSrv,
		workerCancel: cancel,
		workerDone:   done,
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

	if s.workerCancel != nil {
		s.workerCancel()
	}
	if s.workerDone != nil {
		select {
		case <-s.workerDone:
		case <-time.After(2 * time.Second):
		}
	}

	closeErr := s.Close()
	if closeErr != nil {
		closeErr = fmt.Errorf("关闭资源失败：%w", closeErr)
	}

	return errors.Join(httpErr, closeErr)
}

func (s *Server) Close() error {
	var dbErr error
	var redisErr error

	if s.db != nil {
		sqlDB, err := s.db.DB()
		if err != nil {
			dbErr = fmt.Errorf("获取底层连接失败：%w", err)
		} else if err := sqlDB.Close(); err != nil {
			dbErr = fmt.Errorf("关闭 MySQL 失败：%w", err)
		}
	}

	if s.rdb != nil {
		if err := s.rdb.Close(); err != nil {
			redisErr = fmt.Errorf("关闭 Redis 失败：%w", err)
		}
	}

	return errors.Join(dbErr, redisErr)
}

func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.Thread{}, &models.Reply{}, &models.ThreadLike{})
}
