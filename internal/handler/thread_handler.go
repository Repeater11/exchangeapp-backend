package handler

import (
	"errors"
	"exchangeapp/internal/dto"
	"exchangeapp/internal/service"
	"net/http"
	"strconv"

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
	page, size := parsePageSize(ctx.Query("page"), ctx.Query("size"))
	resp, err := h.svc.List(page, size)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取帖子失败"})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ThreadHandler) Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "帖子 ID 无效"})
		return
	}
	threadID := uint(id64)

	resp, err := h.svc.GetByID(threadID)
	if err != nil {
		if errors.Is(err, service.ErrThreadNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取帖子信息失败"})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
