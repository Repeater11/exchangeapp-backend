package handler

import (
	"errors"
	"exchangeapp/internal/repository"
	"exchangeapp/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ThreadLikeHandler struct {
	svc *service.ThreadLikeService
}

func NewThreadLikeHandler(svc *service.ThreadLikeService) *ThreadLikeHandler {
	return &ThreadLikeHandler{svc: svc}
}

func (h *ThreadLikeHandler) Like(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	threadID, ok := parseUintParam(ctx, "id", "帖子 ID 无效")
	if !ok {
		return
	}

	if err := h.svc.Like(userID, threadID); err != nil {
		if errors.Is(err, service.ErrThreadNotFound) {
			jsonError(ctx, http.StatusNotFound, "帖子不存在")
			return
		}
		if errors.Is(err, repository.ErrAlreadyLiked) {
			jsonError(ctx, http.StatusConflict, "已点赞")
			return
		}
		jsonError(ctx, http.StatusInternalServerError, "点赞帖子失败")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "点赞帖子成功"})
}

func (h *ThreadLikeHandler) Unlike(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	threadID, ok := parseUintParam(ctx, "id", "帖子 ID 无效")
	if !ok {
		return
	}

	if err := h.svc.Unlike(userID, threadID); err != nil {
		if errors.Is(err, service.ErrThreadNotFound) {
			jsonError(ctx, http.StatusNotFound, "帖子不存在")
			return
		}
		if errors.Is(err, repository.ErrLikeNotFound) {
			jsonError(ctx, http.StatusConflict, "未点赞")
			return
		}
		jsonError(ctx, http.StatusInternalServerError, "取消点赞帖子失败")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "取消点赞帖子成功"})
}
