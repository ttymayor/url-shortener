package app

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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
	authHandler := handler.NewAuthHandler(cfg)

	// Setup router
	r := gin.Default()

	// Setup CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	// Setup session middleware
	// Use API Key as session secret for simplicity, or add a dedicated SESSION_SECRET env var back if preferred.
	// Here we'll use a hardcoded fallback or a dedicated secret if available,
	// but reusing API Key as a secret is NOT recommended for production security.
	// Let's assume we use a string "super-secret-session-key" for now or read from env.
	store := cookie.NewStore([]byte(cfg.Auth.SessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions("mysession", store))
	router.RegisterRoutes(r, urlHandler, authHandler)

	return &App{engine: r}
}

func (a *App) Run() error {
	return a.engine.Run(":8080")
}
