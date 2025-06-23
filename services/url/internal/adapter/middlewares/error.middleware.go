package middlewares

import (
	"reflect"

	"github.com/tuan-dd/go-pkg/common/constants"
	"github.com/tuan-dd/go-pkg/common/response"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/helper"

	"github.com/gofiber/fiber/v3"
)

// TODO if Error status 500 to send notify discord
func ErrorHandler(c fiber.Ctx, err error) error {
	if err != nil && reflect.ValueOf(err).Kind() == reflect.Pointer {
		internalError := response.NewAppError(err.Error(), constants.InternalServerErr)
		if appErr, ok := err.(*response.AppError); ok {
			internalError = appErr
		}

		return c.Status(constants.HttpCode[internalError.Code]).JSON(response.ErrorResponse(helper.GetHttpReqCtx(c), internalError))
	}

	return nil
}
