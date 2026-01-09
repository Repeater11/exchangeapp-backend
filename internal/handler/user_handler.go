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
	if !bindJSON(ctx, &req) {
		return
	}

	resp, err := h.svc.Register(req)
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			jsonError(ctx, http.StatusConflict, "用户名已存在")
			return
		}
		jsonError(ctx, http.StatusInternalServerError, "注册失败")
		return
	}
	ctx.JSON(http.StatusCreated, resp)
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var req dto.LoginReq
	if !bindJSON(ctx, &req) {
		return
	}

	resp, err := h.svc.Login(req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			jsonError(ctx, http.StatusUnauthorized, "用户名或密码错误")
			return
		}
		jsonError(ctx, http.StatusInternalServerError, "登录失败")
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Me(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}
	username, _ := ctx.Get("username")
	ctx.JSON(http.StatusOK, gin.H{
		"id":       userID,
		"username": username,
	})
}
