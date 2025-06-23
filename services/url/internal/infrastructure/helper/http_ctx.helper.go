package helper

import (
	"github.com/tuan-dd/go-pkg/common/constants"
	"github.com/tuan-dd/go-pkg/common/request"

	"github.com/gofiber/fiber/v3"
)

func GetHttpUserCtx(c fiber.Ctx) *request.UserInfo[any] {
	return c.Locals(constants.REQUEST_CONTEXT_KEY).(*request.ReqContext).UserInfo
}

func GetHttpReqCtx(c fiber.Ctx) *request.ReqContext {
	ctc, oke := c.Locals(constants.REQUEST_CONTEXT_KEY).(*request.ReqContext)
	if !oke {
		return request.BuildRequestContext(nil, nil, nil, &request.UserInfo[any]{})
	}
	return ctc
}

func SetHttpUserCtx(c fiber.Ctx, user *request.UserInfo[any]) {
	c.Locals(constants.REQUEST_CONTEXT_KEY).(*request.ReqContext).UserInfo = user
}
