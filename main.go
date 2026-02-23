package main

import (
	"context"
	"docvault/config"
	"docvault/factory"
	"docvault/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	f, err := factory.New(cfg)
	if err != nil {
		log.Fatal("Failed to init factory: ", err)
	}

	r := gin.Default()

	r.Use(middleware.LoggingMiddleware())

	r.GET("/health", func(c *gin.Context) {
		if err := f.DB.Ping(); err != nil {
			c.JSON(500, gin.H{"error": "database unavailable"})
			return
		}
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	r.POST("/api/documents/upload", f.DocumentHandler.Upload)
	r.GET("/api/documents", f.DocumentHandler.List)
	r.GET("/api/documents/:id", f.DocumentHandler.GetMetadata)
	r.GET("/api/documents/:id/download", f.DocumentHandler.Download)
	r.DELETE("/api/documents/:id", f.DocumentHandler.Delete)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("Server error:", err)
		}
		wg.Done()
	}()
	wg.Add(1)

	go func() {
		f.NotificationWorker.Start(ctx)
		wg.Done()
	}()
	go func() {
		f.SchedulerWorker.Start(ctx)
		wg.Done()
	}()

	<-quit

	log.Println("Shutting down...")
	cancel()
	server.Shutdown(context.Background())
	wg.Wait()
	log.Println("Server stopped gracefully")
}
