package bootstrap

import (
	"github.com/gin-gonic/gin"

	"github.com/mahditd/url-shortener/internal/application/usecase"
	"github.com/mahditd/url-shortener/internal/domain/ports"
	"github.com/mahditd/url-shortener/internal/presentation/controller"
	"github.com/mahditd/url-shortener/internal/presentation/middleware"
	routes "github.com/mahditd/url-shortener/internal/presentation/routes/http"
)

func NewApp(baseURL string, repository ports.LinkRepository, rateLimiter *middleware.RateLimiter) *gin.Engine {

	router := gin.Default()

	linkUsecase := usecase.NewLinkUsecase(repository, baseURL)

	linkHandler := controller.NewLinkHandler(linkUsecase)

	routes.SetUpRoutes(router, linkHandler, rateLimiter)

	return router
}
