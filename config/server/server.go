package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"projectwebcurhat/config"
	dbConfig "projectwebcurhat/config/database"
	"projectwebcurhat/config/middleware"
	"projectwebcurhat/controller"
	dbMigration "projectwebcurhat/database"
	"projectwebcurhat/repository"
	"projectwebcurhat/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Run() {
	log.Println("Starting application...")

	cfg := config.Get()
	if cfg == nil {
		log.Fatal("Failed to load configuration")
		return
	}

	// Connect to database
	db, _, err := dbConfig.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
		return
	}

	// Run migrations
	if err := dbMigration.RunMigration(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
		return
	}

	startServer(cfg, db)
}

func startServer(cfg *config.AppConfig, db *gorm.DB) {
	repo := repository.New(db)
	serv := service.New(repo)

	if cfg.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.CORSMiddleware())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "WebRTC Signaling Server - ProjectWebCurhat")
	})

	r.StaticFile("/test-client", "./test-client.html")

	controller.New(r, serv)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Server is running on port %d", cfg.Port)
	log.Printf("WebSocket endpoint: ws://localhost:%d/ws", cfg.Port)
	log.Fatal(srv.ListenAndServe())
}
