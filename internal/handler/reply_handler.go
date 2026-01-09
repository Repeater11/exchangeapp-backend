package handler

import (
	"errors"
	"exchangeapp/internal/dto"
	"exchangeapp/internal/service"
	"net/http"

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
	if !bindJSON(ctx, &req) {
		return
	}

	threadID, ok := parseUintParam(ctx, "id", "帖子 ID 无效")
	if !ok {
		return
	}

	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	resp, err := h.svc.Create(userID, threadID, req)
	if err != nil {
		if errors.Is(err, service.ErrThreadNotFound) {
			jsonError(ctx, http.StatusNotFound, "帖子不存在")
			return
		}
		jsonError(ctx, http.StatusInternalServerError, "回复失败")
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

func (h *ReplyHandler) ListByThreadID(ctx *gin.Context) {
	threadID, ok := parseUintParam(ctx, "id", "帖子 ID 无效")
	if !ok {
		return
	}

	page, size := parsePageSize(ctx.Query("page"), ctx.Query("size"))
	resp, err := h.svc.ListByThreadID(threadID, page, size)
	if err != nil {
		if errors.Is(err, service.ErrThreadNotFound) {
			jsonError(ctx, http.StatusNotFound, "帖子不存在")
			return
		}
		jsonError(ctx, http.StatusInternalServerError, "获取回复失败")
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *ReplyHandler) ListMine(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	page, size := parsePageSize(ctx.Query("page"), ctx.Query("size"))
	resp, err := h.svc.ListByUserID(userID, page, size)
	if err != nil {
		jsonError(ctx, http.StatusInternalServerError, "获取回复失败")
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ReplyHandler) Update(ctx *gin.Context) {
	var req dto.UpdateReplyReq
	if !bindJSON(ctx, &req) {
		return
	}

	replyID, ok := parseUintParam(ctx, "id", "评论 ID 无效")
	if !ok {
		return
	}

	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	resp, err := h.svc.Update(userID, replyID, req)
	if err != nil {
		if errors.Is(err, service.ErrReplyNotFound) {
			jsonError(ctx, http.StatusNotFound, "评论不存在")
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			jsonError(ctx, http.StatusForbidden, "没有修改权限")
			return
		}
		jsonError(ctx, http.StatusInternalServerError, "修改失败")
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ReplyHandler) Delete(ctx *gin.Context) {
	replyID, ok := parseUintParam(ctx, "id", "评论 ID 无效")
	if !ok {
		return
	}

	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	if err := h.svc.Delete(userID, replyID); err != nil {
		if errors.Is(err, service.ErrReplyNotFound) {
			jsonError(ctx, http.StatusNotFound, "评论不存在")
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			jsonError(ctx, http.StatusForbidden, "没有删除权限")
			return
		}
		jsonError(ctx, http.StatusInternalServerError, "删除失败")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
