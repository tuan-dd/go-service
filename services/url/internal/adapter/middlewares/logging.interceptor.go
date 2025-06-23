package middlewares

import (
	"github.com/tuan-dd/go-pkg/appLogger"
	"github.com/tuan-dd/go-pkg/common/constants"
	"github.com/tuan-dd/go-pkg/common/request"
	"github.com/tuan-dd/go-pkg/common/response"

	"github.com/gofiber/fiber/v3"
)

func LoggingInterceptor(log *appLogger.Logger) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		requestContext := c.Locals(constants.REQUEST_CONTEXT_KEY).(*request.ReqContext)
		log.ReqClientLog(requestContext, c.Method(), c.Path())

		err := c.Next()
		errApp := response.ConvertError(err)
		code := constants.Success
		if errApp != nil && errApp.Code >= 9000 {
			code = errApp.Code
			log.ResClientLog(requestContext, uint(code), errApp)
		} else {
			log.ResClientLog(requestContext, uint(c.Response().StatusCode()), nil)
		}

		return err
	}
}
