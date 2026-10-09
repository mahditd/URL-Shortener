package controller_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mahditd/url-shortener/internal/application/usecase"
	"github.com/mahditd/url-shortener/internal/domain/ports"
	"github.com/mahditd/url-shortener/internal/presentation/controller"
	routes "github.com/mahditd/url-shortener/internal/presentation/routes/http"
)

func setupTestRouterWithRepository(repository ports.LinkRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	linkUsecase := usecase.NewLinkUsecase(repository, "http://localhost:8080")

	linkHandler := controller.NewLinkHandler(linkUsecase)

	routes.SetUpRoutes(router, linkHandler, nil)

	return router
}

func TestShortenRepositoryError(t *testing.T) {

	repository := FakeRepository{
		FindByURLError: errors.New("unknown persistence error"),
	}

	router := setupTestRouterWithRepository(&repository)

	body := `{"url":"https://go.dev/doc/"}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/shorten",
		bytes.NewBufferString(body),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}

}

func TestRedirectRepositoryError(t *testing.T) {

	repository := FakeRepository{
		FindByCodeError: errors.New("unknown persistence error"),
	}
	router := setupTestRouterWithRepository(&repository)

	req := httptest.NewRequest(
		http.MethodGet,
		"/something",
		nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}

}

func TestMetadataRepositoryError(t *testing.T) {
	repository := FakeRepository{
		FindByCodeError: errors.New("unknown persistence error"),
	}
	router := setupTestRouterWithRepository(&repository)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/links/something",
		nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}
