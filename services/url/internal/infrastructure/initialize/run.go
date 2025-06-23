package initialize

import (
	"github.com/tuan-dd/go-pkg/appLogger"
	"github.com/tuan-dd/go-pkg/caching/memory"
	"github.com/tuan-dd/go-pkg/common"
	"github.com/tuan-dd/go-pkg/common/response"
	"github.com/tuan-dd/go-pkg/database/mysql"
	"github.com/tuan-dd/go-pkg/extractor"
	"github.com/tuan-dd/go-service/url/internal/adapter/controllers"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/core"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/global"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/routers"
	"github.com/tuan-dd/go-service/url/internal/usecase"
)

func must[T any](value T, err *response.AppError) T {
	if err != nil {
		panic(err)
	}
	return value
}

func Run() {
	global.AppConfig = must(common.LoadConfig[global.Config](""))
	global.Log = must(appLogger.NewLogger(global.AppConfig.Logger, global.AppConfig.Server))
	global.App = must(core.NewHttpServer(global.AppConfig.AllowOrigins, global.Log))
	global.SQLDB = must(mysql.NewConnection(global.AppConfig.SQL))

	global.Orm = must(core.NewEntClient(global.SQLDB, global.AppConfig.SQL, global.Log))

	global.Ext = extractor.New()

	cacheInst := memory.NewMemoryCache(1000, 0)
	useCase := usecase.NewUrlUseCase(global.Orm.Client(), cacheInst)
	urlController := controllers.NewUrlController(useCase)
	routers.UrlRouters(global.App, urlController)

	global.App.Start(global.Log, shutdown)
}

func shutdown() {
	// if global.InternalBroker != nil {
	// 	if err := global.InternalBroker.Shutdown(); err != nil {
	// 		global.Log.Error("failed to shutdown internal broker", err)
	// 	}
	// }

	if global.Orm != nil {
		if err := global.Orm.Shutdown(); err != nil {
			global.Log.Error("failed to shutdown ent client", err)
		}
	}
	if global.SQLDB != nil {
		if err := global.SQLDB.Shutdown(); err != nil {
			global.Log.Error("failed to shutdown sql connection", err)
		}
	}
	if global.Cache != nil {
		if err := global.Cache.Shutdown(); err != nil {
			global.Log.Error("failed to shutdown cache", err)
		}
	}
}
