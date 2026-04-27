package api

import (
	"github.com/danielgtaylor/huma/v2"
)

func CreateHumaConfig(apiBasePath string) huma.Config {
	humaConfig := huma.DefaultConfig("API Reference Documentation", "1.0.3")
	humaConfig.DocsPath = ""
	humaConfig.OpenAPI.Servers = []*huma.Server{
		{
			URL:         apiBasePath,
			Description: "Base Server",
		},
	}
	humaConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"BearerAuth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}
	return humaConfig
}

func AddOpenAPITags(openapi *huma.OpenAPI) {
	tags := []*huma.Tag{
		{
			Name:        "Users",
			Description: "Operations related to user authentication.",
		},
		{
			Name:        "Reviews",
			Description: "Operations related to reviews.",
		},
		{
			Name:        "Files",
			Description: "File upload/download services.",
		},
		{
			Name:        "Books",
			Description: "Operations related to books.",
		},
		{
			Name:        "Book Pages",
			Description: "Operations related to book pages.",
		},
		{
			Name:        "Pets",
			Description: "Operations related to pets from Petstore service.",
		},
	}

	openapi.Tags = append(openapi.Tags, tags...)
}
