package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mahditd/url-shortener/internal/application/usecase"
	"github.com/mahditd/url-shortener/internal/infrastructure/persistence/memory"
	"github.com/mahditd/url-shortener/internal/presentation/controller"
	routes "github.com/mahditd/url-shortener/internal/presentation/routes/http"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	repository := memory.NewLinkRepository()

	linkUsecase := usecase.NewLinkUsecase(repository, "http://localhost:8080")

	linkHandler := controller.NewLinkHandler(linkUsecase)

	routes.SetUpRoutes(router, linkHandler)

	return router
}

func TestShorten(t *testing.T) {

	type shortenResponse struct {
		Code     string `json:"code"`
		ShortURL string `json:"short_url"`
	}
	router := setupTestRouter()

	body := `{"url":"https://go.dev/doc/"}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/shorten",
		bytes.NewBufferString(body),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var response shortenResponse

	err := json.NewDecoder(w.Body).Decode(&response)

	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Code == "" {
		t.Fatal("expected code but got empty")
	}

	if response.ShortURL == "" {
		t.Fatal("expected short_url but got empty")
	}

	if response.ShortURL != "http://localhost:8080/"+response.Code {
		t.Fatalf("unexpected short url: %s", response.ShortURL)
	}

	if len(response.Code) < 6 || len(response.Code) > 8 {
		t.Fatalf("invalid code length: %d", len(response.Code))
	}

}

func TestShortenIdempotency(t *testing.T) {
	type shortenResponse struct {
		Code     string `json:"code"`
		ShortURL string `json:"short_url"`
	}
	router := setupTestRouter()

	body := `{"url":"https://go.dev/doc/"}`

	req1 := httptest.NewRequest(
		http.MethodPost,
		"/api/shorten",
		bytes.NewBufferString(body),
	)
	req1.Header.Set("Content-Type", "application/json")

	req2 := httptest.NewRequest(
		http.MethodPost,
		"/api/shorten",
		bytes.NewBufferString(body),
	)
	req2.Header.Set("Content-Type", "application/json")

	w1 := httptest.NewRecorder()

	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w1.Code)
	}

	var response1 shortenResponse

	err := json.NewDecoder(w1.Body).Decode(&response1)

	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	w2 := httptest.NewRecorder()

	router.ServeHTTP(w2, req2)

	if w1.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w2.Code)
	}

	var response2 shortenResponse

	err = json.NewDecoder(w2.Body).Decode(&response2)

	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response1.ShortURL != response2.ShortURL {
		t.Fatal("expected same short URL")
	}

	if response1.Code != response2.Code {
		t.Fatal("expected to be equal")
	}

}

func TestRedirect(t *testing.T) {

	type shortenResponse struct {
		Code     string `json:"code"`
		ShortURL string `json:"short_url"`
	}
	router := setupTestRouter()

	body := `{"url":"https://go.dev/doc/"}`

	req1 := httptest.NewRequest(
		http.MethodPost,
		"/api/shorten",
		bytes.NewBufferString(body),
	)
	req1.Header.Set("Content-Type", "application/json")

	w1 := httptest.NewRecorder()

	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w1.Code)
	}

	var response shortenResponse

	err := json.NewDecoder(w1.Body).Decode(&response)

	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	req2 := httptest.NewRequest(
		http.MethodGet,
		"/"+response.Code,
		nil,
	)
	req2.Header.Set("Content-Type", "application/json")

	w2 := httptest.NewRecorder()

	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusFound {
		t.Fatalf("expected status 302, got %d", w2.Code)
	}

	location := w2.Header().Get("Location")

	if location != "https://go.dev/doc/" {
		t.Fatalf("unexpected location: %s", location)
	}

}

func TestInvalidURL(t *testing.T) {

	tests := []struct {
		name string
		body string
	}{
		{
			name: "missing url",
			body: `{}`,
		},
		{
			name: "empty url",
			body: `{"url":""}`,
		},
		{
			name: "invalid scheme",
			body: `{"url":"ftp://example.com"}`,
		},
		{
			name: "random string",
			body: `{"url":"hello"}`,
		},
	}



	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter()

			req := httptest.NewRequest(
				http.MethodPost,
				"/api/shorten",
				bytes.NewBufferString(tt.body),
			)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400, got %d", w.Code)
			}

		})
	}

}

func TestUnknownCode(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(
		http.MethodGet,
		"/randomThing",
		nil,
	)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}