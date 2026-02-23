package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if recoveredValue := recover(); recoveredValue != nil {
				log.Printf("Panic: %v", recoveredValue)
				ctx.AbortWithStatusJSON(500, gin.H{"message": "Server Error"})
			}
		}()

		ctx.Next()
	}
}
