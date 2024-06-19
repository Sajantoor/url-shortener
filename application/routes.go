package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func createRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	return router
}

func loadRoutes() *gin.Engine {
	router := createRouter()
	router.GET("health", healthCheck)

	publicRoutes := router.Group("v1/")
	shortenRoutes := publicRoutes.Group("shorten")

	shortenRoutes.POST("/", shortenURL)
	shortenRoutes.GET("/:shortURL", redirectToOriginalURL)

	return router
}

func healthCheck(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "Up"})
}

func shortenURL(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "Shorten URL"})
}

func redirectToOriginalURL(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "Redirect to original URL"})
}
