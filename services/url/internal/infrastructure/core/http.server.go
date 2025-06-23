package core

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/tuan-dd/go-pkg/appLogger"
	"github.com/tuan-dd/go-pkg/common/response"

	"github.com/tuan-dd/go-service/url/internal/adapter/middlewares"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/helper"
	"go.uber.org/zap"
	"golang.org/x/sys/unix"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"

	"github.com/gofiber/fiber/v3/middleware/recover"
)

type HttpServer struct {
	App  *fiber.App
	Host string
	Port int
	Name string
}

func NewHttpServer(allowOrigins string, log *appLogger.Logger) (*HttpServer, *response.AppError) {
	// Configure Fiber app

	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		ErrorHandler: middlewares.ErrorHandler,
	})

	app.Use(
		helmet.New(
			helmet.Config{
				CrossOriginEmbedderPolicy: "false",
				ContentSecurityPolicy:     "false",
				CrossOriginOpenerPolicy:   "cross-origin",
				CrossOriginResourcePolicy: "cross-origin",
			},
		),
		cors.New(
			cors.Config{
				AllowOrigins:     strings.Split(allowOrigins, ","),
				AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
				AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
				AllowCredentials: false,
			},
		),
		recover.New(),
		// update new version
		// fiberzap.New(
		// 	fiberzap.Config{
		// 		Logger: log.Logger(),
		// 	},
		// ),
		middlewares.ReqContextHandler,
		middlewares.LoggingInterceptor(log),
	)

	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("I'm good!")
	})

	return &HttpServer{
		App:  app,
		Host: "localhost",
		Port: 8080,
		Name: "MyApp",
	}, nil
}

func (s *HttpServer) Start(log *appLogger.Logger, shutDown func()) *response.AppError {
	defer func() {
		_ = log.Logger().Sync()
	}()
	// Graceful shutdown setup
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, unix.SIGTERM, unix.SIGINT, unix.SIGTSTP)

	// Start server in goroutine
	addr := fmt.Sprintf(":%d", s.Port)

	serverErr := make(chan error, 1)
	go func() {
		if err := s.App.Listen(addr); err != nil {
			serverErr <- err
		}
	}()

	s.App.Use(func() fiber.Handler {
		return func(c fiber.Ctx) error {
			return c.Status(http.StatusNotFound).JSON(response.ErrorResponse(helper.GetHttpReqCtx(c), response.NotFound("route not found")))
		}
	}())

	// Wait for interrupt signal or server error
	log.Info(fmt.Sprintf("Http server is running on %s", addr))
	select {
	case err := <-serverErr:
		return response.ServerError(fmt.Sprintf("server error: %v", err.Error()))
	case sig := <-sigChan:
		log.Info("received signal", zap.String("signal", sig.String()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shutDown()

	done := make(chan struct{})
	go func() {
		s.App.ShutdownWithContext(ctx)
		close(done)
	}()

	select {
	case <-done:
		log.Info("graceful shutdown completed", nil)
	case <-ctx.Done():
		log.Info("graceful shutdown timeout", nil)
	}

	return nil
}
