package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"projectwebcurhat/config"
	"projectwebcurhat/config/middleware"
	"projectwebcurhat/controller"
	"projectwebcurhat/repository"
	"projectwebcurhat/service"
)

func Run() {
	log.Println("Starting application...")

	cfg := config.Get()
	if cfg == nil {
		log.Fatal("Failed to load configuration")
		return
	}

	startServer(cfg)
}

func startServer(cfg *config.AppConfig) {
	repo := repository.New()
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
