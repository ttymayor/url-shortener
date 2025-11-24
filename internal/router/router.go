package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ttymayor/url-shortener/internal/handler"
)

func RegisterRoutes(r *gin.Engine, h *handler.URLHandler, auth *handler.AuthHandler) {
	// Public API endpoints
	api := r.Group("/api")
	{
		api.POST("/login", auth.Login)
		api.POST("/logout", auth.Logout)
		api.GET("/check-session", auth.CheckSession)
	}

	// Protected API endpoints
	protected := r.Group("/api")
	protected.Use(handler.AuthMiddleware())
	{
		protected.POST("/shorten", h.Shorten)
	}

	// Redirect endpoint (Public)
	r.GET("/:code", h.Redirect)
}
