package unit

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuan-dd/go-service/url/internal/adapter/dtos"
)

func TestHashFunctionConsistency(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Close()

	tests := []struct {
		name string
		url  string
	}{
		{
			name: "simple URL",
			url:  "https://example.com",
		},
		{
			name: "URL with path",
			url:  "https://example.com/path/to/resource",
		},
		{
			name: "URL with query params",
			url:  "https://example.com/search?q=test&type=web",
		},
		{
			name: "URL with fragment",
			url:  "https://example.com/page#section1",
		},
		{
			name: "complex URL",
			url:  "https://api.example.com:8080/v1/users/123?include=profile&format=json#metadata",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var codes []string

			// Generate the same hash multiple times
			for range 10 {
				req := dtos.UrlCreateReq{URL: tt.url}
				result, appErr := helper.GetUseCase().Create(context.Background(), req)
				var err error
				if appErr != nil {
					err = appErr.Wrap()
				}
				require.NoError(t, err)
				codes = append(codes, result.ID)
			}

			// All codes should be identical
			firstCode := codes[0]
			for i, code := range codes {
				assert.Equal(t, firstCode, code, "Hash %d should match the first hash", i)
			}

			// Verify the code meets requirements
			assert.True(t, len(firstCode) <= 8, "Code should not exceed 8 characters")
			assert.True(t, len(firstCode) >= 6, "Code should be at least 6 characters")

			// Verify the code contains only valid characters
			for _, char := range firstCode {
				assert.True(t, strings.ContainsRune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", char),
					"Code should contain only valid base62 characters")
			}
		})
	}
}

func TestHashCollisionResistance(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Close()

	// Generate URLs that are similar but should produce different hashes
	similarUrls := []string{
		"https://example.com/path1",
		"https://example.com/path2",
		"https://example.com/path3",
		"https://example.org/path1",
		"https://examples.com/path1",
		"https://example.com/path1?param=1",
		"https://example.com/path1?param=2",
		"https://example.com/Path1", // Different case
	}

	generatedCodes := make(map[string]string) // code -> URL mapping

	for _, url := range similarUrls {
		req := dtos.UrlCreateReq{URL: url}
		result, appErr := helper.GetUseCase().Create(context.Background(), req)
		var err error
		if appErr != nil {
			err = appErr.Wrap()
		}
		require.NoError(t, err)

		if existingURL, exists := generatedCodes[result.ID]; exists {
			t.Logf("Potential collision detected: '%s' and '%s' both generated code '%s'",
				url, existingURL, result.ID)
		} else {
			generatedCodes[result.ID] = url
		}
	}

	assert.True(t, len(generatedCodes) > 1, "Should generate different codes for different URLs")
}

func TestHashDeterminism(t *testing.T) {
	// Test that the same URL always produces the same hash across different instances
	helper1 := NewTestHelper(t)
	defer helper1.Close()

	helper2 := NewTestHelper(t)
	defer helper2.Close()

	testURL := "https://determinism-test.com/path"

	// Create URL with first instance
	req1 := dtos.UrlCreateReq{URL: testURL}
	result1, appErr1 := helper1.GetUseCase().Create(context.Background(), req1)
	var err1 error
	if appErr1 != nil {
		err1 = appErr1.Wrap()
	}
	require.NoError(t, err1)

	// Create URL with second instance
	req2 := dtos.UrlCreateReq{URL: testURL}
	result2, appErr2 := helper2.GetUseCase().Create(context.Background(), req2)
	var err2 error
	if appErr2 != nil {
		err2 = appErr2.Wrap()
	}
	require.NoError(t, err2)

	// Both should produce the same code
	assert.Equal(t, result1.ID, result2.ID,
		"Same URL should produce same hash across different instances")
}

func TestHashDistribution(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Close()

	// Generate many URLs and check distribution of first character
	const numUrls = 100
	charCount := make(map[rune]int)

	for i := 0; i < numUrls; i++ {
		url := fmt.Sprintf("https://test-%d.com", i)
		req := dtos.UrlCreateReq{URL: url}
		result, appErr := helper.GetUseCase().Create(context.Background(), req)
		var err error
		if appErr != nil {
			err = appErr.Wrap()
		}
		require.NoError(t, err)

		if len(result.ID) > 0 {
			firstChar := rune(result.ID[0])
			charCount[firstChar]++
		}
	}

	// Should have reasonable distribution (not all codes starting with same character)
	assert.True(t, len(charCount) > 1, "Hash should have reasonable distribution")

	// Log distribution for analysis
	for char, count := range charCount {
		t.Logf("Character '%c' appears %d times as first character", char, count)
	}
}
