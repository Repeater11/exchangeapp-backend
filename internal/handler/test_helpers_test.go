package handler

import "github.com/gin-gonic/gin"

func testAuthMiddleware(userID uint) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if userID != 0 {
			ctx.Set("userID", userID)
			ctx.Set("username", "tester")
		}
		ctx.Next()
	}
}
