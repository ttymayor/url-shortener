package app

import (
	"github.com/ttymayor/url-shortener/internal/config"
	"github.com/ttymayor/url-shortener/internal/handler"
	"github.com/ttymayor/url-shortener/internal/repository"
	"github.com/ttymayor/url-shortener/internal/router"
	"github.com/ttymayor/url-shortener/internal/service"

	"github.com/gin-gonic/gin"
)

type App struct {
	engine *gin.Engine
}

func NewApp() *App {
	cfg := config.Load()

	// Initialize connections
	redisClient := repository.NewRedisClient(cfg)
	pgClient := repository.NewPostgresClient(cfg)

	// Initialize layers
	urlRepo := repository.NewURLRepository(pgClient, redisClient)
	urlService := service.NewURLService(urlRepo)
	urlHandler := handler.NewURLHandler(urlService)

	// Setup router
	r := gin.Default()
	router.RegisterRoutes(r, urlHandler)

	return &App{engine: r}
}

func (a *App) Run() error {
	return a.engine.Run(":8081")
}
