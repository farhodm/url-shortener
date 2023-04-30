package handlers

import (
	"github.com/farhodm/url-shortener/pkg/models"
	"github.com/farhodm/url-shortener/pkg/shortener"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h *handler) CreateShortUrl(c *gin.Context) {
	var link models.UrlShorter
	err := c.ShouldBindJSON(&link)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ch := link.LongUrl
	link.LongUrl = ch
	if err := h.DB.Where("long_url =?", link.LongUrl).First(&link).Error; err != nil {
		link.ShortUrl = shortener.RandStr(6)
		if h.DB.Create(&link).Error != nil {
			log.Println("inserting link data to DB:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error",
			})
			return
		}
	}
	host := "http://localhost:9999/"
	c.JSON(200, gin.H{
		"message":   "short url created successfully",
		"short_url": host + link.ShortUrl,
	})
}

func (h *handler) HandlerShortUrlRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	var link models.UrlShorter
	if err := h.DB.Where("short_url =?", shortUrl).First(&link).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "I can't find such links!",
		})
		return
	}

	c.Redirect(http.StatusMovedPermanently, link.LongUrl)
}
