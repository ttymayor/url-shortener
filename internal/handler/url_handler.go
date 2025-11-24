package handler

import (
	"errors"
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
	URL  string `json:"url" binding:"required"`
	Code string `json:"code"` // Optional custom code
}

func (h *URLHandler) Shorten(c *gin.Context) {
	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, err := h.service.ShortenURL(req.URL, req.Code)
	if err != nil {
		if errors.Is(err, service.ErrCodeInUse) {
			c.JSON(http.StatusConflict, gin.H{"error": "Custom short code is already in use"})
			return
		}
		if errors.Is(err, service.ErrInvalidURL) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format"})
			return
		}
		if errors.Is(err, service.ErrInvalidCode) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid short code format or reserved word"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to shorten URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_code":   url.ShortCode,
		"original_url": url.OriginalURL,
	})
}

func (h *URLHandler) GetAllURLs(c *gin.Context) {
	urls, err := h.service.GetAllURLs()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URLs"})
		return
	}

	c.JSON(http.StatusOK, urls)
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
