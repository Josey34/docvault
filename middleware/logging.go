package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		path := ctx.Request.URL.Path

		startTime := time.Now()

		ctx.Next()

		status := ctx.Writer.Status()
		duration := time.Since(startTime)

		log.Printf("[%s] %s %s | %d | %.2fms",
			time.Now().Format("2006-01-02 15:04:05"),
			method,
			path,
			status,
			duration.Seconds()*1000)
	}
}
