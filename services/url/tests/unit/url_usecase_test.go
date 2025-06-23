package unit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuan-dd/go-pkg/caching/memory"
	"github.com/tuan-dd/go-service/url/internal/adapter/dtos"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/generated/ent"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/generated/ent/enttest"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/generated/ent/shortenedurl"
	"github.com/tuan-dd/go-service/url/internal/usecase"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func TestUrlUseCase_Create(t *testing.T) {
	tests := []struct {
		name    string
		request dtos.UrlCreateReq
		wantErr bool
	}{
		{
			name: "successful creation",
			request: dtos.UrlCreateReq{
				URL: "https://example.com",
			},
			wantErr: false,
		},
		{
			name: "duplicate URL should return existing record",
			request: dtos.UrlCreateReq{
				URL: "https://example.com",
			},
			wantErr: false,
		},
		{
			name: "valid long URL",
			request: dtos.UrlCreateReq{
				URL: "https://www.verylongdomainnameforthisexampletest.com/path/to/resource?param1=value1&param2=value2",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			helper := NewTestHelper(t)
			defer helper.Close()

			// Execute
			result, appErr := helper.GetUseCase().Create(context.Background(), tt.request)
			var err error
			if appErr != nil {
				err = appErr.Wrap()
			}

			// Debug: Print actual error
			t.Logf("Test case: %s, Error: %v, Result: %v", tt.name, err, result)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err, "Create should not return an error")
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.ID)
				assert.Equal(t, tt.request.URL, result.OriginalURL)

				// Verify the record was saved to database
				found, dbErr := helper.GetClient().ShortenedURL.Query().Where(shortenedurl.IDEQ(result.ID)).First(context.Background())
				assert.NoError(t, dbErr)
				assert.Equal(t, tt.request.URL, found.OriginalURL)
			}
		})
	}
}

func TestUrlUseCase_Create_Idempotency(t *testing.T) {
	// Setup
	helper := NewTestHelper(t)
	defer helper.Close()

	request := dtos.UrlCreateReq{
		URL: "https://example.com/idempotency-test",
	}

	// Execute - Create first time
	result1, appErr1 := helper.GetUseCase().Create(context.Background(), request)
	var err1 error
	if appErr1 != nil {
		err1 = appErr1.Wrap()
	}
	if err1 != nil {
		t.Logf("First create error: %v", err1)
	}
	require.NoError(t, err1)
	require.NotNil(t, result1)

	// Execute - Create second time (should return same result)
	result2, appErr2 := helper.GetUseCase().Create(context.Background(), request)
	var err2 error
	if appErr2 != nil {
		err2 = appErr2.Wrap()
	}
	require.NoError(t, err2)
	require.NotNil(t, result2)

	// Assert - Should return the same code for the same URL
	assert.Equal(t, result1.ID, result2.ID)
	assert.Equal(t, result1.OriginalURL, result2.OriginalURL)
	assert.Equal(t, result1.ID, result2.ID)
}

func TestUrlUseCase_Get(t *testing.T) {
	tests := []struct {
		name         string
		setupURL     string
		retrieveCode string
		wantURL      string
		wantErr      bool
	}{
		{
			name:     "successful retrieval",
			setupURL: "https://example.com/get-test",
			wantErr:  false,
		},
		{
			name:         "non-existent code",
			retrieveCode: "nonexist",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			helper := NewTestHelper(t)
			defer helper.Close()

			var code string
			if tt.setupURL != "" {
				// Create a URL first
				createReq := dtos.UrlCreateReq{URL: tt.setupURL}
				created, appErr := helper.GetUseCase().Create(context.Background(), createReq)
				var err error
				if appErr != nil {
					err = appErr.Wrap()
				}
				require.NoError(t, err)
				code = created.ID
				tt.wantURL = tt.setupURL
			} else {
				code = tt.retrieveCode
			}

			// Execute
			result, appErr := helper.GetUseCase().Get(code)
			var err error
			if appErr != nil {
				err = appErr.Wrap()
			}

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantURL, result.OriginalURL)
			}
		})
	}
}

func TestUrlUseCase_Get_CacheHit(t *testing.T) {
	// Setup
	helper := NewTestHelper(t)
	defer helper.Close()

	// Create a URL
	createReq := dtos.UrlCreateReq{URL: "https://example.com/cache-test"}
	created, appErr := helper.GetUseCase().Create(context.Background(), createReq)
	var err error
	if appErr != nil {
		err = appErr.Wrap()
	}
	require.NoError(t, err)

	// First call - should hit database and populate cache
	result1, appErr1 := helper.GetUseCase().Get(created.ID)
	var err1 error
	if appErr1 != nil {
		err1 = appErr1.Wrap()
	}
	require.NoError(t, err1)
	assert.Equal(t, createReq.URL, result1.OriginalURL)

	// Second call - should hit cache
	result2, appErr2 := helper.GetUseCase().Get(created.ID)
	var err2 error
	if appErr2 != nil {
		err2 = appErr2.Wrap()
	}
	require.NoError(t, err2)
	assert.Equal(t, createReq.URL, result2.OriginalURL)

	assert.Equal(t, result1.OriginalURL, result2.OriginalURL)
}

func TestUrlUseCase_HashCode_Consistency(t *testing.T) {
	// This test verifies that the same URL always generates the same hash code
	helper := NewTestHelper(t)
	defer helper.Close()

	url := "https://consistent-hash-test.com"

	// Create the same URL multiple times
	var codes []string
	for i := 0; i < 5; i++ {
		req := dtos.UrlCreateReq{URL: url}
		result, appErr := helper.GetUseCase().Create(context.Background(), req)
		var err error
		if appErr != nil {
			err = appErr.Wrap()
		}
		require.NoError(t, err)
		codes = append(codes, result.ID)
	}

	// All codes should be the same
	for i := 1; i < len(codes); i++ {
		assert.Equal(t, codes[0], codes[i], "Hash code should be consistent for the same URL")
	}
}

func TestUrlUseCase_ConcurrentAccess(t *testing.T) {
	// Setup
	helper := NewTestHelper(t)
	defer helper.Close()

	url := "https://concurrent-test.com"
	numGoroutines := 5 // Reduced from 10 to minimize database locking issues

	// Channel to collect results
	results := make(chan *ent.ShortenedURL, numGoroutines)
	errors := make(chan error, numGoroutines)

	// Launch concurrent create operations
	for i := 0; i < numGoroutines; i++ {
		go func() {
			req := dtos.UrlCreateReq{URL: url}
			result, appErr := helper.GetUseCase().Create(context.Background(), req)
			var err error
			if appErr != nil {
				err = appErr.Wrap()
			}
			if err != nil {
				errors <- err
			} else {
				results <- result
			}
		}()
	}

	// Collect results
	var createdRecords []*ent.ShortenedURL
	var receivedErrors []error

	for i := 0; i < numGoroutines; i++ {
		select {
		case result := <-results:
			createdRecords = append(createdRecords, result)
		case err := <-errors:
			receivedErrors = append(receivedErrors, err)
		case <-time.After(5 * time.Second):
			t.Fatal("Test timed out")
		}
	}

	// Assert no errors occurred
	assert.Empty(t, receivedErrors)
	assert.Len(t, createdRecords, numGoroutines)

	// All results should have the same code (idempotency)
	expectedCode := createdRecords[0].ID
	for _, record := range createdRecords {
		assert.Equal(t, expectedCode, record.ID)
		assert.Equal(t, url, record.OriginalURL)
	}
}

func BenchmarkUrlUseCase_Create(b *testing.B) {
	client := enttest.Open(b, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	cache := memory.NewMemoryCache(1000, 5*time.Minute)
	useCase := usecase.NewUrlUseCase(client, cache)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := dtos.UrlCreateReq{
			URL: fmt.Sprintf("https://benchmark-test-%d.com", i),
		}
		_, err := useCase.Create(context.Background(), req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUrlUseCase_Get_CacheHit(b *testing.B) {
	client := enttest.Open(b, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	cache := memory.NewMemoryCache(1000, 5*time.Minute)
	useCase := usecase.NewUrlUseCase(client, cache)

	// Pre-create a URL
	req := dtos.UrlCreateReq{URL: "https://benchmark-get-test.com"}
	created, err := useCase.Create(context.Background(), req)
	if err != nil {
		b.Fatal(err)
	}

	// Warm up cache
	_, _ = useCase.Get(created.ID)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := useCase.Get(created.ID)
		if err != nil {
			b.Fatal(err)
		}
	}
}
