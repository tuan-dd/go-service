package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/tuan-dd/go-pkg/common/response"
	"github.com/tuan-dd/go-service/url/internal/adapter/dtos"
	"github.com/tuan-dd/go-service/url/internal/usecase"
)

// UrlController wires HTTP handlers with the usecase logic.
type UrlController struct {
	useCase *usecase.UrlUseCase
}

func NewUrlController(useCase *usecase.UrlUseCase) *UrlController {
	return &UrlController{useCase: useCase}
}

func (uc *UrlController) Redirect(c fiber.Ctx) error {
	code := c.Params("code")
	urlStr, err := uc.useCase.Redirect(code)
	if err != nil {
		return err
	}
	return c.Redirect().Status(http.StatusFound).To(urlStr)
}

func (uc *UrlController) Create(c fiber.Ctx) error {
	var req dtos.UrlCreateReq

	if err := c.Bind().Body(&req); err != nil {
		return response.QueryInvalid("invalid request body")
	}
	isUrlValid := req.IsValidReq()

	if !isUrlValid {
		return response.QueryInvalid("invalid request body")
	}

	record, err := uc.useCase.Create(c.Context(), req)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(dtos.ToUrlCreateRes(record))
}

func (uc *UrlController) Get(c fiber.Ctx) error {
	code := c.Params("code")
	isValidCode := dtos.ValidCode(code)
	if code == "" || !isValidCode {
		return response.QueryInvalid("invalid code")
	}
	record, err := uc.useCase.Get(code)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(dtos.ToUrlGetRes(record))
}
