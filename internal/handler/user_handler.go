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
