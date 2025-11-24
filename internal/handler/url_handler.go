package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ttymayor/url-shortener/internal/service"
)

type URLHandler struct {
	service service.URLService
}

func NewURLHandler(service service.URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}

func (h *URLHandler) Shorten(c *gin.Context) {
	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, err := h.service.ShortenURL(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to shorten URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_code":   url.ShortCode,
		"original_url": url.OriginalURL,
	})
}

func (h *URLHandler) Redirect(c *gin.Context) {
	shortCode := c.Param("code")

	originalURL, err := h.service.GetOriginalURL(shortCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if originalURL == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}
