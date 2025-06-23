package dtos

import (
	"fmt"
	"strings"
	"time"

	"github.com/tuan-dd/go-service/url/internal/infrastructure/generated/ent"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/global"
)

type (
	UrlCreateReq struct {
		URL       string     `json:"url"`
		ExpiredAt *time.Time `json:"expired_at,omitempty"`
	}

	UrlCreateRes struct {
		ID       string `json:"short_code"`
		ShortURL string `json:"short_url"`
	}

	UrlGetRes struct {
		ID          string    `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		OriginalURL string    `json:"original_url"`
	}
)

func (url *UrlCreateReq) IsValidReq() bool {
	if url.URL == "" {
		return false
	}

	if !strings.HasPrefix(url.URL, "http://") && !strings.HasPrefix(url.URL, "https://") {
		return false
	}

	if len(url.URL) < 10 || len(url.URL) > 2048 {
		return false
	}

	protocolEnd := strings.Index(url.URL, "://") + 3
	if protocolEnd >= len(url.URL) {
		return false
	}

	domain := url.URL[protocolEnd:]
	isValid := strings.Contains(domain, ".") && !strings.HasPrefix(domain, ".") && !strings.HasSuffix(domain, ".")

	if !isValid {
		return false
	}

	if url.ExpiredAt != nil {
		expiredAtUTC := url.ExpiredAt.UTC()
		if isValid = time.Now().UTC().Before(expiredAtUTC); !isValid {
			return false
		}
		url.ExpiredAt = &expiredAtUTC
	}

	return true
}

func ShortURL(code string) string {
	return strings.TrimRight(fmt.Sprintf("http://%s:%d/shortURLs/%s/redirect", global.App.Host, global.App.Port, code), "/")
}

func ToUrlCreateRes(record *ent.ShortenedURL) *UrlCreateRes {
	return &UrlCreateRes{
		ID:       record.ID,
		ShortURL: ShortURL(record.ID),
	}
}

func ToUrlGetRes(record *ent.ShortenedURL) *UrlGetRes {
	return &UrlGetRes{
		ID:          record.ID,
		CreatedAt:   record.CreatedAt,
		OriginalURL: record.OriginalURL,
	}
}

func ValidCode(code string) bool {
	if len(code) < 6 || len(code) > 12 {
		return false
	}

	for _, char := range code {
		if !strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", char) {
			return false
		}
	}

	return true
}
