package handler

import (
	"errors"
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
	if !bindJSON(ctx, &req) {
		return
	}

	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	resp, err := h.svc.Create(userID, req)
	if err != nil {
		jsonError(ctx, http.StatusInternalServerError, "发帖失败")
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

func (h *ThreadHandler) List(ctx *gin.Context) {
	page, size := parsePageSize(ctx.Query("page"), ctx.Query("size"))
	resp, err := h.svc.List(page, size)
	if err != nil {
		jsonError(ctx, http.StatusInternalServerError, "获取帖子失败")
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ThreadHandler) ListMine(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	page, size := parsePageSize(ctx.Query("page"), ctx.Query("size"))
	resp, err := h.svc.ListByUserID(userID, page, size)
	if err != nil {
		jsonError(ctx, http.StatusInternalServerError, "获取帖子失败")
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ThreadHandler) Detail(ctx *gin.Context) {
	threadID, ok := parseUintParam(ctx, "id", "帖子 ID 无效")
	if !ok {
		return
	}

	resp, err := h.svc.GetByID(threadID)
	if err != nil {
		if errors.Is(err, service.ErrThreadNotFound) {
			jsonError(ctx, http.StatusNotFound, "帖子不存在")
			return
		}
		jsonError(ctx, http.StatusInternalServerError, "获取帖子信息失败")
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *ThreadHandler) Update(ctx *gin.Context) {
	var req dto.UpdateThreadReq
	if !bindJSON(ctx, &req) {
		return
	}

	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	threadID, ok := parseUintParam(ctx, "id", "帖子 ID 无效")
	if !ok {
		return
	}

	resp, err := h.svc.Update(userID, threadID, req)
	if err != nil {
		if errors.Is(err, service.ErrThreadNotFound) {
			jsonError(ctx, http.StatusNotFound, "帖子不存在")
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

func (h *ThreadHandler) Delete(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	threadID, ok := parseUintParam(ctx, "id", "帖子 ID 无效")
	if !ok {
		return
	}

	if err := h.svc.Delete(userID, threadID); err != nil {
		if errors.Is(err, service.ErrThreadNotFound) {
			jsonError(ctx, http.StatusNotFound, "帖子不存在")
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
