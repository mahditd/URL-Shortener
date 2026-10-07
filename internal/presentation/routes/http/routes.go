package http

import (
	"github.com/gin-gonic/gin"
	"github.com/mahditd/url-shortener/internal/presentation/controller"
)

func SetUpRoutes(router *gin.Engine, linkHandler *controller.LinkHandler) {

	api := router.Group("/api")
	{
		api.POST("/shorten", linkHandler.Shorten)

		v1 := api.Group("/v1")
		{
			v1.GET("/links/:code", linkHandler.Metadata)
		}

	}

	router.GET("/:code", linkHandler.Redirect)
}
