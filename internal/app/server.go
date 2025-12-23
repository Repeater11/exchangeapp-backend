package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	e *gin.Engine
}

func NewServer() (*Server, error) {
	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery())

	e.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return &Server{e: e}, nil
}

func (s *Server) Run(addr int) error {
	return s.e.Run(fmt.Sprintf(":%d", addr))
}
