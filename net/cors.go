package net

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 设置跨域
func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		origin := ctx.Request.Header.Get("Origin") // 请求头部
		// 接收客户端发送的origin
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		// 服务器支持的所有跨域请求的方法
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		// 允许跨域设置可以返回其他字段，可以自定义字段
		ctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token, session, timeStamp, randomStr, signature, content-type")
		//  允许浏览器可以解析的头部
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		// 允许客户端传递校验信息比如
		ctx.Header("Access-Control-Allow-Credentials", "true")
		// 允许类型校验
		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
		ctx.Next()
	}
}
