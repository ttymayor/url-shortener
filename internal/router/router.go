package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ttymayor/url-shortener/internal/handler"
)

func RegisterRoutes(r *gin.Engine, h *handler.URLHandler) {
	// API endpoints
	api := r.Group("/api")
	{
		api.POST("/shorten", h.Shorten)
		api.GET("/urls", h.GetAllURLs)
	}

	// Redirect endpoint
	r.GET("/:code", h.Redirect)
}
