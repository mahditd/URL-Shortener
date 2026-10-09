package http

import (
	"github.com/gin-gonic/gin"
	"github.com/mahditd/url-shortener/internal/presentation/controller"
	"github.com/mahditd/url-shortener/internal/presentation/middleware"
)

func SetUpRoutes(router *gin.Engine, linkHandler *controller.LinkHandler, rateLimiter *middleware.RateLimiter) {

	api := router.Group("/api")
	{
		if rateLimiter != nil {
			api.POST(
				"/shorten",
				middleware.RateLimitMiddleware(rateLimiter),
				linkHandler.Shorten,
			)
		} else {
			api.POST(
				"/shorten",
				linkHandler.Shorten,
			)
		}

		v1 := api.Group("/v1")
		{
			v1.GET("/links/:code", linkHandler.Metadata)
		}

	}

	router.GET("/:code", linkHandler.Redirect)
}
