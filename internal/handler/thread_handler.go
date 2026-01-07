package handler

import (
	"exchangeapp/internal/dto"
	"exchangeapp/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ThreadHandler struct {
	svc *service.ThreadService
}

func NewThreadHandler(svc *service.ThreadService) *ThreadHandler {
	return &ThreadHandler{svc: svc}
}

func (h *ThreadHandler) Create(ctx *gin.Context) {
	var req dto.CreateThreadReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	userIDVal, ok := ctx.Get("userID")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "用户信息异常"})
		return
	}

	resp, err := h.svc.Create(userID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "发帖失败"})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

func (h *ThreadHandler) List(ctx *gin.Context) {
	resp, err := h.svc.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取帖子失败"})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
