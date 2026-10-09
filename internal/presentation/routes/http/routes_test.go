package http

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mahditd/url-shortener/internal/application/usecase"
	"github.com/mahditd/url-shortener/internal/infrastructure/persistence/memory"
	"github.com/mahditd/url-shortener/internal/presentation/controller"
)

func TestSetUpRoutes(t *testing.T) {

	gin.SetMode(gin.TestMode)
	repository := memory.NewLinkRepository()

	linkUsecase := usecase.NewLinkUsecase(repository, "http://localhost:8080")

	handler := controller.NewLinkHandler(linkUsecase)

	router := gin.New()

	SetUpRoutes(router, handler, nil)

	routes := router.Routes()

	expectedRoutes := map[string]bool{
		"POST /api/shorten":       false,
		"GET /api/v1/links/:code": false,
		"GET /:code":              false,
	}

	for _, route := range routes {
		key := route.Method + " " + route.Path

		if _, exists := expectedRoutes[key]; exists {
			expectedRoutes[key] = true
		}
	}

	for route, found := range expectedRoutes {
		if !found {
			t.Fatalf("expected route %s to be registered", route)
		}
	}

}
