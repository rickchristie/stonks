package main

import (
	"net"

	"github.com/gin-gonic/gin"

	accsNote "app-template/accessor/note"
	"app-template/lib/lb"
	psNote "app-template/pservice/app/note"
	"app-template/pservice/middleware"
	svcNote "app-template/service/note"
)

var Version = "dev"

func main() {
	cfg := MustLoadConfig()

	corsConfig := middleware.MustParseCORSConfig(cfg.CorsAllowedOrigins)
	r := gin.Default()
	r.Use(corsConfig.Middleware())

	r.GET("/health-check", func(c *gin.Context) {
		c.String(200, "Gesundheit!")
	})

	r.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{"version": Version})
	})

	connSelector := lb.Randomized([]string{cfg.ConnStr()})
	noteStorage := accsNote.NewPgStorage(connSelector, connSelector)
	noteSvc := svcNote.NewAppService(noteStorage)
	noteHandler := psNote.NewHandler(noteSvc)
	noteHandler.RegisterRoutes(r)

	if err := r.Run(net.JoinHostPort(cfg.BackendHost, cfg.BackendPort)); err != nil {
		panic("failed to run backend: " + err.Error())
	}
}
