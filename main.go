package main

import (
	"docvault/config"
	"docvault/factory"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	f, err := factory.New(cfg)
	if err != nil {
		log.Fatal("Failed to init factory: ", err)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		if err := f.DB.Ping(); err != nil {
			c.JSON(500, gin.H{"error": "database unavailable"})
			return
		}
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	r.Run(":" + cfg.Port)
}
