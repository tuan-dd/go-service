package request

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tuan-dd/go-pkg/common/constants"
)

type (
	ReqContext struct {
		CID              string `json:"cid"`
		IP               string
		RequestTimestamp int64 `json:"request_timestamp"`
		UserInfo         *UserInfo[any]
		AccessToken      *string `json:"access_token"`
		RefreshToken     *string `json:"refresh_token"`
	}

	UserInfo[T any] struct {
		ID    string
		Name  string
		Email string
		Role  []string
		Data  T
	}
)

func BuildRequestContext(cid *string, accessToken *string, refreshToken *string, userInfo *UserInfo[any]) *ReqContext {
	return &ReqContext{
		CID:              GetCid(cid),
		RequestTimestamp: time.Now().UnixMilli(),
		AccessToken:      accessToken,
		UserInfo:         userInfo,
		RefreshToken:     refreshToken,
	}
}

func GetCid(cid *string) string {
	if cid != nil && *cid != "" {
		return *cid
	}

	return generateUUIDV7()
}

func CalculateDuration(timeStamp int64) int64 {
	return time.Now().UnixMilli() - timeStamp
}

func FormatMilliseconds(timeStamp int64) string {
	return fmt.Sprintf("%dms", CalculateDuration(timeStamp))
}

func GetUserCtx[T any](ctx context.Context) *UserInfo[T] {
	reqCtx := GetReqCtx(ctx)
	res, _ := structToStruct[UserInfo[any], UserInfo[T]](reqCtx.UserInfo)

	return res
}

func GetReqCtx(ctx context.Context) *ReqContext {
	reqCtx, ok := ctx.Value(constants.REQUEST_CONTEXT_KEY).(*ReqContext)
	if !ok || reqCtx == nil {
		return BuildRequestContext(nil, nil, nil, &UserInfo[any]{})
	}

	return reqCtx
}

func SetUserCtx(ctx context.Context, user *UserInfo[any]) context.Context {
	reqCtx := ctx.Value(constants.REQUEST_CONTEXT_KEY).(*ReqContext)

	if reqCtx == nil {
		reqCtx = BuildRequestContext(nil, nil, nil, user)
		ctx = context.WithValue(ctx, constants.REQUEST_CONTEXT_KEY, reqCtx)
	} else {
		reqCtx.UserInfo = user
	}

	return ctx
}

func structToStruct[T, U any](t *T) (*U, error) {
	var u U
	b, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func generateUUIDV7() string {
	newCid, _ := uuid.NewV7()

	return newCid.String()
}
