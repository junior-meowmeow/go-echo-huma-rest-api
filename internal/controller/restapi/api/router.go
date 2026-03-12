package api

import (
	v1 "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/api/v1"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

func NewRouter(handlers *handler.Handlers, apiBasePath string) *echo.Echo {
	router := echo.New()
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	RegisterDocumentations(router, apiBasePath)

	humaConfig := CreateHumaConfig(apiBasePath)
	api := humaecho.New(router, humaConfig)

	RegisterRoutes(api, handlers)
	v1.RegisterGroup(api, handlers)

	return router
}

func RegisterDocumentations(router *echo.Echo, apiBasePath string) {
	router.GET("/docs", StoplightElements(apiBasePath))
	router.GET("/docs/scalar", ScalarDocs(apiBasePath))
	router.GET("/docs/swagger", SwaggerUI(apiBasePath))
}

func CreateHumaConfig(apiBasePath string) huma.Config {
	humaConfig := huma.DefaultConfig("API Reference Documentation", "1.0.0")
	humaConfig.DocsPath = ""
	humaConfig.OpenAPI.Servers = []*huma.Server{
		{
			URL:         apiBasePath,
			Description: "Base Server",
		},
	}
	// disable the $schema property
	humaConfig.CreateHooks = nil
	return humaConfig
}
