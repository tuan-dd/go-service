package unit

import (
	"time"

	"entgo.io/ent/dialect"
	"github.com/tuan-dd/go-pkg/caching/memory"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/generated/ent"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/generated/ent/enttest"
	"github.com/tuan-dd/go-service/url/internal/usecase"
)

// TestHelper provides common test utilities
type TestHelper struct {
	client  *ent.Client
	cache   *memory.MemoryCache
	useCase *usecase.UrlUseCase
}

// NewTestHelper creates a new test helper with initialized dependencies
func NewTestHelper(t TestingT) *TestHelper {
	client := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	cache := memory.NewMemoryCache(100, 5*time.Minute)
	useCase := usecase.NewUrlUseCase(client, cache)

	return &TestHelper{
		client:  client,
		cache:   cache,
		useCase: useCase,
	}
}

// Close cleans up test resources
func (h *TestHelper) Close() {
	if h.client != nil {
		h.client.Close()
	}
}

// GetClient returns the Ent client
func (h *TestHelper) GetClient() *ent.Client {
	return h.client
}

// GetCache returns the memory cache
func (h *TestHelper) GetCache() *memory.MemoryCache {
	return h.cache
}

// GetUseCase returns the URL use case
func (h *TestHelper) GetUseCase() *usecase.UrlUseCase {
	return h.useCase
}

// TestingT represents the interface used by testing.T and testing.B
type TestingT interface {
	Error(args ...any)
	FailNow()
}

// Common test data
var (
	TestURLs = []string{
		"https://example.com",
		"https://google.com",
		"https://github.com",
		"https://stackoverflow.com",
		"https://developer.mozilla.org/en-US/docs/Web/JavaScript",
	}

	InvalidURLs = []string{
		"",
		"not-a-url",
		"ftp://example.com",
		"http://",
		"https://",
		"http://a",
		"https://.com",
		"https://example.",
	}

	ValidCodes = []string{
		"abc123",
		"XYZ789",
		"Mix3d4",
		"123456",
		"ABCDEF",
		"abcdef",
	}

	InvalidCodes = []string{
		"",
		"12345",         // too short
		"1234567890123", // too long
		"abc!23",        // invalid chars
		"abc 123",       // spaces
		"abc-123",       // dash
		"abc_123",       // underscore
	}
)
