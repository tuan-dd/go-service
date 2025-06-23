package core

import (
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/tuan-dd/go-pkg/appLogger"
	"github.com/tuan-dd/go-pkg/common/response"
	"github.com/tuan-dd/go-pkg/database/mysql"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/generated/ent"
	_ "github.com/tuan-dd/go-service/url/internal/infrastructure/generated/ent/runtime"
)

type EntClient struct {
	client *ent.Client
}

func NewEntClient(sqlDb *mysql.Connection, cgf *mysql.SQLConfig, log *appLogger.Logger) (*EntClient, *response.AppError) {
	var err *response.AppError
	if sqlDb == nil {
		err = response.ServerError("sqlDb is nil")
		return nil, err
	}
	drv := entsql.OpenDB(cgf.RDBMS, sqlDb.DB())

	if cgf.LogEnabled {
		drvDebug := dialect.DebugWithContext(drv, log.DBLog)

		return &EntClient{client: ent.NewClient(ent.Driver(drvDebug))}, nil
	}

	return &EntClient{client: ent.NewClient(ent.Driver(drv))}, nil
}

func (c *EntClient) Client() *ent.Client {
	return c.client
}

func (c *EntClient) Shutdown() *response.AppError {
	if c.client == nil {
		return response.ServerError("ent client is nil")
	}
	if err := c.client.Close(); err != nil {
		return response.ServerError("failed to close ent client: " + err.Error())
	}
	return nil
}
