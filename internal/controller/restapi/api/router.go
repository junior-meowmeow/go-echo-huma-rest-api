package api

import (
	v1 "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/api/v1"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
	customMiddleware "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/middleware"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/util"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

func NewRouter(handlers *handler.Handlers, utilities *util.Utilities, apiBasePath string) *echo.Echo {
	router := echo.New()
	AddEchoMiddlewares(router)
	RegisterDocumentations(router, apiBasePath)

	humaConfig := CreateHumaConfig(apiBasePath)
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
	router.Use(middleware.Logger())
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	router.Use(middleware.Secure())
	router.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
}
