package api

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/config"
	v1 "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/api/v1"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
	customMiddleware "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/middleware"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/utility"
)

func NewRouter(handlers *handler.Handlers, utilities *utility.Utilities, appConfig config.AppConfig) *echo.Echo {
	router := echo.New()
	AddEchoMiddlewares(router)
	AddEchoPrometheus(router)
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

func AddEchoMiddlewares(router *echo.Echo) {
	router.Use(middleware.Recover())
	router.Use(middleware.RequestID())
	router.Use(middleware.RequestLogger())
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	router.Use(middleware.Secure())
	router.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
}

func AddEchoPrometheus(router *echo.Echo) {
	router.Use(echoprometheus.NewMiddleware("myapp"))
	router.GET("/metrics", echoprometheus.NewHandler())
}
