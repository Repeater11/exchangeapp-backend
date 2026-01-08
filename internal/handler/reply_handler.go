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

	page, size := parsePageSize(ctx.Query("page"), ctx.Query("size"))
	resp, err := h.svc.ListByThreadID(threadID, page, size)
	if err != nil {
		if errors.Is(err, service.ErrThreadNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取回复失败"})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *ReplyHandler) Update(ctx *gin.Context) {
	var req dto.UpdateReplyReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "评论 ID 无效"})
		return
	}
	replyID := uint(id64)

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

	resp, err := h.svc.Update(userID, replyID, req)
	if err != nil {
		if errors.Is(err, service.ErrReplyNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "没有修改权限"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "修改失败"})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ReplyHandler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "评论 ID 无效"})
		return
	}
	replyID := uint(id64)

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

	if err := h.svc.Delete(userID, replyID); err != nil {
		if errors.Is(err, service.ErrReplyNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "没有删除权限"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
