package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"

	"github.com/tuan-dd/go-pkg/caching/memory"
	"github.com/tuan-dd/go-pkg/common/response"
	"github.com/tuan-dd/go-service/url/internal/adapter/dtos"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/generated/ent"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/generated/ent/shortenedurl"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/persistence"
)

const (
	charset     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	prefix      = "url-shortener-"
)

type UrlUseCase struct {
	client *ent.ShortenedURLClient
	cache  *memory.MemoryCache
}

func NewUrlUseCase(client *ent.Client, cache *memory.MemoryCache) *UrlUseCase {
	return &UrlUseCase{client: client.ShortenedURL, cache: cache}
}

func (u *UrlUseCase) buildKey(code string) string {
	return prefix + code
}

func (u *UrlUseCase) Create(ctx context.Context, dto dtos.UrlCreateReq) (*ent.ShortenedURL, *response.AppError) {
	maxRetries := 5
	var code string
	for i := range maxRetries {
		hashInput := dto.URL
		if i > 0 {
			hashInput = fmt.Sprintf("%s_%d", dto.URL, i)
		}

		code = hashCode(hashInput)
		existing, checkErr := u.client.Query().Where(shortenedurl.IDEQ(code)).Only(ctx)
		if checkErr != nil {
			if ent.IsNotFound(checkErr) {
				break
			}
			return nil, persistence.ConvertEntErr(checkErr, "check existing short URL")
		}

		if existing.OriginalURL == dto.URL {
			return existing, nil
		}

		if i == maxRetries-1 {
			return nil, response.ServerError("failed to generate unique short URL after multiple attempts")
		}
	}

	if record, err := u.client.Create().SetID(code).SetOriginalURL(dto.URL).SetNillableExpiredAt(dto.ExpiredAt).Save(ctx); err != nil {
		return nil, persistence.ConvertEntErr(err, "create short URL")
	} else {
		u.cache.Set(u.buildKey(code), dto.URL)
		return record, nil
	}
}

func (u *UrlUseCase) Redirect(code string) (string, *response.AppError) {
	if val, ok := u.cache.Get(u.buildKey(code)); ok {
		return val, nil
	}

	urlStr, err := u.client.Query().Where(shortenedurl.IDEQ(code)).Select(shortenedurl.FieldOriginalURL).Only(context.Background())
	if err != nil {
		return "", persistence.ConvertEntErr(err, "get original URL")
	}

	if time.Now().After(urlStr.ExpiredAt) && !urlStr.ExpiredAt.IsZero() {
		return "", response.QueryInvalid("URL has expired")
	}

	u.cache.Set(u.buildKey(code), urlStr.OriginalURL)
	return urlStr.OriginalURL, nil
}

func (u *UrlUseCase) Get(code string) (*ent.ShortenedURL, *response.AppError) {
	record, err := u.client.Query().Where(shortenedurl.IDEQ(code)).First(context.Background())
	if err != nil {
		return nil, persistence.ConvertEntErr(err, "get original URL")
	}

	if time.Now().After(record.ExpiredAt) && !record.ExpiredAt.IsZero() {
		return nil, response.QueryInvalid("URL has expired")
	}

	return record, nil
}

func hashCode(original string) string {
	sum := sha256.Sum256([]byte(original))
	enc := base62Encode(sum[:])
	if len(enc) > 8 {
		enc = enc[:8]
	}
	return enc
}

func base62Encode(input []byte) string {
	n := new(big.Int).SetBytes(input)
	if n.Sign() == 0 {
		return string(base62Chars[0])
	}
	base := big.NewInt(62)
	var result []byte
	for n.Sign() > 0 {
		mod := new(big.Int)
		n.DivMod(n, base, mod)
		result = append(result, base62Chars[mod.Int64()])
	}
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return string(result)
}
