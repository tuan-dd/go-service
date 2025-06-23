package middlewares

import (
	"github.com/tuan-dd/go-pkg/common/constants"
	"github.com/tuan-dd/go-pkg/common/request"

	"github.com/gofiber/fiber/v3"
)

func ReqContextHandler(c fiber.Ctx) error {
	cid := fiber.Locals[string](c, constants.CORRELATION_ID_KEY)

	authorization := c.Get(string(constants.AUTHORIZATION_KEY))
	if authorization != "" {
		authorization = authorization[7:]
	}

	requestContext := request.BuildRequestContext(&cid, &authorization, nil, &request.UserInfo[any]{})
	c.Locals(constants.REQUEST_CONTEXT_KEY, requestContext)
	c.Locals(constants.CORRELATION_ID_KEY, requestContext.CID)

	return c.Next()
}
