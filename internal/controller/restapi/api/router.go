package api

import (
	"context"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/config"
	v1 "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/api/v1"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
	customMiddleware "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/middleware"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/utility"
)

func NewRouter(handlers *handler.Handlers, utilities *utility.Utilities, appConfig config.AppConfig) *echo.Echo {
	router := echo.New()
	AddEchoMiddlewares(router, appConfig)
	AddEchoPrometheus(router, appConfig)
	RegisterDocumentations(router, appConfig.APIBasePath)

	humaConfig := CreateHumaConfig(appConfig.APIBasePath)
	AddOpenAPITags(humaConfig.OpenAPI)
	api := humaecho.New(router, humaConfig)

	public := huma.NewGroup(api, "")

	protected := huma.NewGroup(api, "")
	protected.UseMiddleware(
		customMiddleware.RequireToken(protected, utilities.Token),
	)
	protected.UseSimpleModifier(func(op *huma.Operation) {
		op.Security = []map[string][]string{
			{"BearerAuth": {}},
		}
	})

	RegisterRoutes(public, protected, handlers)
	v1.RegisterGroup(public, protected, handlers)

	return router
}

func AddEchoMiddlewares(router *echo.Echo, appConfig config.AppConfig) {
	router.Use(middleware.Recover())
	router.Use(customRequestID())
	router.Use(customRequestLogger())
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	router.Use(middleware.Secure())
	if appConfig.RequestPerSecLimit != 0 {
		router.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(appConfig.RequestPerSecLimit))))
	}
}

func AddEchoPrometheus(router *echo.Echo, appConfig config.AppConfig) {
	if !appConfig.EnablePrometheus {
		return
	}
	router.Use(echoprometheus.NewMiddleware("myapp"))
	router.GET("/metrics", echoprometheus.NewHandler())
}

func customRequestID() echo.MiddlewareFunc {
	return middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		RequestIDHandler: func(c echo.Context, requestID string) {
			req := c.Request()
			ctx := context.WithValue(req.Context(), handler.RequestIDKey, requestID)
			c.SetRequest(req.WithContext(ctx))
		},
	})
}

func customRequestLogger() echo.MiddlewareFunc {
	skipper := func(c echo.Context) bool {
		path := c.Request().URL.Path
		return path == "/health" || path == "/metrics"
	}
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogLatency:       true,
		LogProtocol:      false,
		LogRemoteIP:      true,
		LogHost:          true,
		LogMethod:        true,
		LogURI:           true,
		LogURIPath:       false,
		LogRoutePath:     false,
		LogRequestID:     true,
		LogReferer:       false,
		LogUserAgent:     true,
		LogStatus:        true,
		LogError:         true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogHeaders:       nil,
		LogQueryParams:   nil,
		LogFormValues:    nil,
		HandleError:      true, // forwards error to the global error handler, so it can decide appropriate status code
		Skipper:          skipper,
		//revive:disable-next-line:unused-parameter
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				slog.LogAttrs(context.Background(), slog.LevelInfo, "HTTP_REQUEST",
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.Duration("latency", v.Latency),
					slog.String("host", v.Host),
					slog.String("bytes_in", v.ContentLength),
					slog.Int64("bytes_out", v.ResponseSize),
					slog.String("user_agent", v.UserAgent),
					slog.String("remote_ip", v.RemoteIP),
					slog.String("request_id", v.RequestID),
				)
			} else {
				slog.LogAttrs(context.Background(), slog.LevelError, "HTTP_REQUEST_ERROR",
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.Duration("latency", v.Latency),
					slog.String("host", v.Host),
					slog.String("bytes_in", v.ContentLength),
					slog.Int64("bytes_out", v.ResponseSize),
					slog.String("user_agent", v.UserAgent),
					slog.String("remote_ip", v.RemoteIP),
					slog.String("request_id", v.RequestID),

					slog.String("error", v.Error.Error()),
				)
			}
			return nil
		},
	})
}
