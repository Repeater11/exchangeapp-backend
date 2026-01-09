package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func jsonError(ctx *gin.Context, code int, msg string) {
	ctx.JSON(code, gin.H{"error": msg})
}

func bindJSON(ctx *gin.Context, v any) bool {
	if err := ctx.ShouldBindJSON(v); err != nil {
		jsonError(ctx, http.StatusBadRequest, "参数错误")
		return false
	}
	return true
}

func getUserID(ctx *gin.Context) (uint, bool) {
	userIDVal, ok := ctx.Get("userID")
	if !ok {
		jsonError(ctx, http.StatusUnauthorized, "未登录")
		return 0, false
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		jsonError(ctx, http.StatusInternalServerError, "用户信息异常")
		return 0, false
	}
	return userID, true
}

func parseUintParam(ctx *gin.Context, name, errMsg string) (uint, bool) {
	raw := ctx.Param(name)
	val, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		jsonError(ctx, http.StatusBadRequest, errMsg)
		return 0, false
	}
	return uint(val), true
}
