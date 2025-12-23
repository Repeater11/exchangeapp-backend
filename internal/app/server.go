package app

import "github.com/gin-gonic/gin"

type Server struct {
	e *gin.Engine
}

func NewServer() (*Server, error) {
	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery())

	return &Server{e: e}, nil
}
