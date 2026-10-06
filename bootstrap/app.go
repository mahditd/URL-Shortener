package bootstrap

import (
	"github.com/gin-gonic/gin"

	"github.com/mahditd/url-shortener/internal/application/usecase"
	"github.com/mahditd/url-shortener/internal/infrastructure/persistence/memory"
	"github.com/mahditd/url-shortener/internal/presentation/controller"
	routes "github.com/mahditd/url-shortener/internal/presentation/routes/http"
)

func NewApp(baseURL string) *gin.Engine {

	router := gin.Default()

	repository := memory.NewLinkRepository()

	linkUsecase := usecase.NewLinkUsecase(repository, baseURL)

	linkHandler := controller.NewLinkHandler(linkUsecase)

	routes.SetUpRoutes(router, linkHandler)

	return router
}
