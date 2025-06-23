package unit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuan-dd/go-pkg/caching/memory"
	"github.com/tuan-dd/go-pkg/common/constants"
	"github.com/tuan-dd/go-pkg/common/request"
	"github.com/tuan-dd/go-service/url/internal/adapter/controllers"
	"github.com/tuan-dd/go-service/url/internal/adapter/dtos"
	"github.com/tuan-dd/go-service/url/internal/adapter/middlewares"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/core"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/generated/ent"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/generated/ent/enttest"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/global"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/routers"
	"github.com/tuan-dd/go-service/url/internal/usecase"
)

func createTestApp() *core.HttpServer {
	// Setup global for testing
	global.App = &core.HttpServer{
		Host: "localhost",
		Port: 8080,
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: middlewares.ErrorHandler,
	})

	// Add middleware to set request context
	app.Use(func(c fiber.Ctx) error {
		requestContext := &request.ReqContext{
			CID:              "test-correlation-id",
			RequestTimestamp: time.Now().Unix(),
			UserInfo:         &request.UserInfo[any]{},
		}
		c.Locals(constants.REQUEST_CONTEXT_KEY, requestContext)
		return c.Next()
	})

	return &core.HttpServer{
		App:  app,
		Host: "localhost",
		Port: 8080,
		Name: "MyApp",
	}
}

func TestUrlController_Create(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    dtos.UrlCreateReq
		expectedStatus int
		expectError    bool
	}{
		{
			name: "successful creation",
			requestBody: dtos.UrlCreateReq{
				URL: "https://example.com",
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "valid URL with path",
			requestBody: dtos.UrlCreateReq{
				URL: "https://example.com/path/to/resource",
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "invalid URL - no protocol",
			requestBody: dtos.UrlCreateReq{
				URL: "example.com",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "invalid URL - empty",
			requestBody: dtos.UrlCreateReq{
				URL: "",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "invalid URL - too short",
			requestBody: dtos.UrlCreateReq{
				URL: "http://a",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			client := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
			defer client.Close()

			cache := memory.NewMemoryCache(100, 5*time.Minute)
			useCase := usecase.NewUrlUseCase(client, cache)
			controller := controllers.NewUrlController(useCase)

			httpServer := createTestApp()
			routers.UrlRouters(httpServer, controller)

			// Prepare request
			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/shortURLs", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Execute
			resp, err := httpServer.App.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Assert
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			if !tt.expectError {
				var response dtos.UrlCreateRes
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)

				assert.NotEmpty(t, response.ID)
				assert.NotEmpty(t, response.ShortURL)
				assert.Contains(t, response.ShortURL, response.ID)
			} else {
				// Should contain error response
				assert.Contains(t, string(body), "error")
			}
		})
	}
}

func TestUrlController_Create_InvalidJSON(t *testing.T) {
	// Setup
	client := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	cache := memory.NewMemoryCache(100, 5*time.Minute)
	useCase := usecase.NewUrlUseCase(client, cache)
	controller := controllers.NewUrlController(useCase)

	httpServer := createTestApp()
	httpServer.App.Post("/shortURLs", controller.Create)

	// Prepare invalid JSON request
	req := httptest.NewRequest(http.MethodPost, "/shortURLs", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, err := httpServer.App.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUrlController_Get(t *testing.T) {
	tests := []struct {
		name           string
		setupURL       string
		requestCode    string
		expectedStatus int
		expectRedirect bool
	}{
		{
			name:           "successful redirect",
			setupURL:       "https://example.com",
			expectedStatus: http.StatusFound,
			expectRedirect: true,
		},
		{
			name:           "successful redirect with path",
			setupURL:       "https://google.com",
			expectedStatus: http.StatusFound,
			expectRedirect: true,
		},
		{
			name:           "non-existent code",
			requestCode:    "nonexist",
			expectedStatus: http.StatusNotFound,
			expectRedirect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			client := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
			defer client.Close()

			cache := memory.NewMemoryCache(100, 5*time.Minute)
			useCase := usecase.NewUrlUseCase(client, cache)
			controller := controllers.NewUrlController(useCase)

			httpServer := createTestApp()
			routers.UrlRouters(httpServer, controller)

			code := &ent.ShortenedURL{}
			if tt.setupURL != "" {
				// Create a URL first
				createReq := dtos.UrlCreateReq{URL: tt.setupURL}
				created, err := useCase.Create(context.Background(), createReq)
				require.NoError(t, err.Wrap())
				code = created
			} else {
				code.ID = tt.requestCode
			}

			// Prepare request
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/shortURLs/%s/redirect", code.ID), nil)

			// Execute
			resp, err := httpServer.App.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Assert

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectRedirect {
				location := resp.Header.Get("Location")
				assert.Equal(t, tt.setupURL, location)
			}
		})
	}
}

func TestUrlController_Get_InvalidCode(t *testing.T) {
	// Setup
	client := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	cache := memory.NewMemoryCache(100, 5*time.Minute)
	useCase := usecase.NewUrlUseCase(client, cache)
	controller := controllers.NewUrlController(useCase)

	httpServer := createTestApp()
	routers.UrlRouters(httpServer, controller)

	// Test with empty code
	req := httptest.NewRequest(http.MethodGet, "/shortURLs/1BlRRSRM", nil)

	resp, err := httpServer.App.Test(req)
	fmt.Println("Response:", resp.Status)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestUrlController_Integration_CreateAndGet(t *testing.T) {
	// Setup
	client := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	cache := memory.NewMemoryCache(100, 5*time.Minute)
	useCase := usecase.NewUrlUseCase(client, cache)
	controller := controllers.NewUrlController(useCase)

	httpServer := createTestApp()
	routers.UrlRouters(httpServer, controller)

	originalURL := "https://example.com/integration-test"

	// Step 1: Create URL
	createReq := dtos.UrlCreateReq{URL: originalURL}
	jsonBody, err := json.Marshal(createReq)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/shortURLs", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpServer.App.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var createResponse dtos.UrlCreateRes
	err = json.Unmarshal(body, &createResponse)
	require.NoError(t, err)

	// Step 2: Get URL using the generated code
	getReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/shortURLs/%s/redirect", createResponse.ID), nil)

	getResp, err := httpServer.App.Test(getReq)
	require.NoError(t, err)
	defer getResp.Body.Close()

	assert.Equal(t, http.StatusFound, getResp.StatusCode)
	assert.Equal(t, originalURL, getResp.Header.Get("Location"))
}

func BenchmarkUrlController_Create(b *testing.B) {
	// Setup
	client := enttest.Open(b, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	cache := memory.NewMemoryCache(1000, 5*time.Minute)
	useCase := usecase.NewUrlUseCase(client, cache)
	controller := controllers.NewUrlController(useCase)

	httpServer := createTestApp()
	httpServer.App.Post("/shortURLs", controller.Create)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		createReq := dtos.UrlCreateReq{
			URL: fmt.Sprintf("https://benchmark-test-%d.com", i),
		}
		jsonBody, _ := json.Marshal(createReq)

		req := httptest.NewRequest(http.MethodPost, "/shortURLs", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpServer.App.Test(req)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}
	}
}

func BenchmarkUrlController_Get(b *testing.B) {
	// Setup
	client := enttest.Open(b, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	cache := memory.NewMemoryCache(1000, 5*time.Minute)
	useCase := usecase.NewUrlUseCase(client, cache)
	controller := controllers.NewUrlController(useCase)

	httpServer := createTestApp()
	httpServer.App.Get("/shortURLs/:code/redirect", controller.Redirect)

	// Pre-create a URL
	createReq := dtos.UrlCreateReq{URL: "https://benchmark-get-test.com"}
	created, err := useCase.Create(context.Background(), createReq)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/shortURLs/%s/redirect", created.ID), nil)

		resp, err := httpServer.App.Test(req)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusFound {
			b.Fatalf("Expected status 302, got %d", resp.StatusCode)
		}
	}
}
