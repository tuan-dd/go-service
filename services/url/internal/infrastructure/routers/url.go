package routers

import (
	"github.com/tuan-dd/go-service/url/internal/adapter/controllers"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/core"
)

func UrlRouters(server *core.HttpServer, urlController *controllers.UrlController) {
	router := server.App.Group("shortURLs/")
	router.Get("/:code", urlController.Get)
	router.Post("", urlController.Create)
	router.Get("/:code/redirect", urlController.Redirect)
}
