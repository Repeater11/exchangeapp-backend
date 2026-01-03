package handler

import (
	"errors"
	"exchangeapp/internal/dto"
	"exchangeapp/internal/repository"
	"exchangeapp/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var req dto.RegisterReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	resp, err := h.svc.Register(req)
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}
	ctx.JSON(http.StatusCreated, resp)
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var req dto.LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	resp, err := h.svc.Login(req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败"})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Me(ctx *gin.Context) {
	userID, ok := ctx.Get("userID")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	username, _ := ctx.Get("username")
	ctx.JSON(http.StatusOK, gin.H{
		"id":       userID,
		"username": username,
	})
}
