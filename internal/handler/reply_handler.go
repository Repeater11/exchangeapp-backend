package handler

import (
	"errors"
	"exchangeapp/internal/dto"
	"exchangeapp/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReplyHandler struct {
	svc *service.ReplyService
}

func NewReplyHandler(svc *service.ReplyService) *ReplyHandler {
	return &ReplyHandler{svc: svc}
}

func (h *ReplyHandler) Create(ctx *gin.Context) {
	var req dto.CreateReplyReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "帖子 ID 无效"})
		return
	}
	threadID := uint(id64)

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

	resp, err := h.svc.Create(userID, threadID, req)
	if err != nil {
		if errors.Is(err, service.ErrThreadNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "回复失败"})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

func (h *ReplyHandler) ListByThreadID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "帖子 ID 无效"})
		return
	}
	threadID := uint(id64)

	resp, err := h.svc.ListByThreadID(threadID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取回复失败"})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
