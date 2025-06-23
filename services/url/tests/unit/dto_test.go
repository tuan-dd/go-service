package unit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tuan-dd/go-service/url/internal/adapter/dtos"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/core"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/generated/ent"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/global"
)

func TestUrlCreateReq_IsValidURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "valid HTTP URL",
			url:  "http://example.com",
			want: true,
		},
		{
			name: "valid HTTPS URL",
			url:  "https://example.com",
			want: true,
		},
		{
			name: "valid URL with path",
			url:  "https://example.com/path/to/resource",
			want: true,
		},
		{
			name: "valid URL with query params",
			url:  "https://example.com/search?q=test&type=web",
			want: true,
		},
		{
			name: "valid URL with subdomain",
			url:  "https://api.example.com",
			want: true,
		},
		{
			name: "valid URL with port",
			url:  "https://example.com:8080",
			want: true,
		},
		{
			name: "empty URL",
			url:  "",
			want: false,
		},
		{
			name: "URL without protocol",
			url:  "example.com",
			want: false,
		},
		{
			name: "invalid protocol",
			url:  "ftp://example.com",
			want: false,
		},
		{
			name: "URL too short",
			url:  "http://a",
			want: false,
		},
		{
			name: "URL without domain",
			url:  "https://",
			want: false,
		},
		{
			name: "URL with invalid domain (no dot)",
			url:  "https://localhost",
			want: false,
		},
		{
			name: "URL starting with dot",
			url:  "https://.example.com",
			want: false,
		},
		{
			name: "URL ending with dot",
			url:  "https://example.com.",
			want: false,
		},
		{
			name: "very long URL (valid)",
			url:  "https://example.com/" + generateLongPath(2000),
			want: true,
		},
		{
			name: "too long URL (invalid)",
			url:  "https://example.com/" + generateLongPath(2100),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := dtos.UrlCreateReq{URL: tt.url}
			result := req.IsValidReq()
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestValidCode(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "valid code with numbers",
			code: "123456",
			want: true,
		},
		{
			name: "valid code with lowercase letters",
			code: "abcdef",
			want: true,
		},
		{
			name: "valid code with uppercase letters",
			code: "ABCDEF",
			want: true,
		},
		{
			name: "valid code mixed case and numbers",
			code: "Ab3Cd7",
			want: true,
		},
		{
			name: "valid maximum length code",
			code: "123456789012",
			want: true,
		},
		{
			name: "valid minimum length code",
			code: "123456",
			want: true,
		},
		{
			name: "too short code",
			code: "12345",
			want: false,
		},
		{
			name: "too long code",
			code: "1234567890123",
			want: false,
		},
		{
			name: "empty code",
			code: "",
			want: false,
		},
		{
			name: "code with special characters",
			code: "abc123!",
			want: false,
		},
		{
			name: "code with spaces",
			code: "abc 123",
			want: false,
		},
		{
			name: "code with dash",
			code: "abc-123",
			want: false,
		},
		{
			name: "code with underscore",
			code: "abc_123",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dtos.ValidCode(tt.code)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestToUrlCreateRes(t *testing.T) {
	// Setup global.App for testing
	global.App = &core.HttpServer{
		Host: "localhost",
		Port: 8080,
	}

	// Create a sample ShortenedURL entity
	entity := &ent.ShortenedURL{
		ID:          "abc123",
		OriginalURL: "https://example.com",
		CreatedAt:   time.Now(),
	}

	// Execute
	result := dtos.ToUrlCreateRes(entity)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, "abc123", result.ID)
	assert.NotEmpty(t, result.ShortURL)
	assert.Contains(t, result.ShortURL, "abc123")
}

// Helper function to generate long paths for testing
func generateLongPath(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[i%len(charset)]
	}
	return string(result)
}

func BenchmarkUrlCreateReq_IsValidURL(b *testing.B) {
	req := dtos.UrlCreateReq{URL: "https://example.com/path/to/resource?param=value"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.IsValidReq()
	}
}

func BenchmarkValidCode(b *testing.B) {
	code := "abc123DEF"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dtos.ValidCode(code)
	}
}
