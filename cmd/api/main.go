package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohammadanang/logistics-api/internal/config"
	delivHttp "github.com/mohammadanang/logistics-api/internal/delivery/http"
	"github.com/mohammadanang/logistics-api/internal/middleware"
	"github.com/mohammadanang/logistics-api/internal/repository"
	"github.com/mohammadanang/logistics-api/internal/usecase"
	"github.com/mohammadanang/logistics-api/pkg/cache"
	"github.com/mohammadanang/logistics-api/pkg/database"
	"github.com/mohammadanang/logistics-api/pkg/paseto"
)

// @title Logistics Management API
// @version 1.0
// @description API for managing last-mile delivery and tracking.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Set Gin Mode
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 3. Initialize Database & Cache
	_ = database.NewPostgresConn(cfg.PostgresDSN)
	_ = cache.NewRedisClient(cfg.RedisAddr)

	db := database.NewPostgresConn(cfg.PostgresDSN)
	rdb := cache.NewRedisClient(cfg.RedisAddr)

	// 4. Wiring Clean Architecture (Setup Dependencies)
	packageRepo := repository.NewPostgresRepository(db)
	cacheRepo := repository.NewRedisRepository(rdb)
	packageUsecase := usecase.NewPackageUsecase(packageRepo, cacheRepo)
	paymentUsecase := usecase.NewPaymentUsecase(cfg.XenditAPIKey)

	// 4.5 Inisialisasi Paseto Maker
	tokenMaker, err := paseto.NewTokenMaker(cfg.PasetoSecretKey)
	if err != nil {
		log.Fatalf("Cannot create token maker: %v", err)
	}

	// 5. Setup Router & Handler
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "API is running smoothly"})
	})

	// Grouping API
	api := router.Group("/api/v1")

	// Rute Publik (Tanpa Middleware)
	delivHttp.NewAuthHandler(api, tokenMaker)                 // Rute: POST /api/v1/login
	delivHttp.NewWebhookHandler(api, cfg.XenditCallbackToken) // Rute: POST /api/v1/webhooks/xendit
	delivHttp.NewPackageHandler(router, packageUsecase)       // (Package Handler sebelumnya)

	// Rute Terproteksi (Menggunakan Middleware Auth)
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(tokenMaker))
	{
		// Contoh rute yang mewajibkan login kurir/admin
		protected.POST("/packages/batch-update", func(c *gin.Context) {
			userID := c.GetString("x-user-id")
			role := c.GetString("x-user-role")
			c.JSON(http.StatusOK, gin.H{
				"message": "Welcome to protected route!",
				"user_id": userID,
				"role":    role,
			})
		})
		delivHttp.NewPaymentHandler(protected, paymentUsecase)
	}

	// 5. Konfigurasi HTTP Server
	srv := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// 6. Jalankan Server di dalam Goroutine (Asinkron)
	go func() {
		log.Printf("Starting server on port %s...\n", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	// 7. Implementasi Graceful Shutdown
	quit := make(chan os.Signal, 1)
	// Menerima sinyal interrupt dari OS (Ctrl+C atau Docker stop)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server gracefully...")

	// Memberikan waktu 5 detik untuk menyelesaikan request yang sedang berjalan
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
