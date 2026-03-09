package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func StoplightElements(apiBasePath string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.HTML(http.StatusOK, `<!doctype html>
			<html lang="en">
			<head>
				<meta charset="utf-8">
				<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
				<title>API Reference</title>
				<!-- Embed elements Elements via Web Component -->
				<script src="https://unpkg.com/@stoplight/elements/web-components.min.js"></script>
				<link rel="stylesheet" href="https://unpkg.com/@stoplight/elements/styles.min.css">
			</head>
			<body>

				<elements-api
				apiDescriptionUrl="`+apiBasePath+`/openapi.yaml"
				router="hash"
				layout="sidebar"
				/>

			</body>
			</html>`)
	}
}

func ScalarDocs(apiBasePath string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.HTML(http.StatusOK, `<!doctype html>
			<html>
			<head>
				<title>API Reference</title>
				<meta charset="utf-8" />
				<meta
				name="viewport"
				content="width=device-width, initial-scale=1" />
			</head>
			<body>
				<script
				id="api-reference"
				data-url="`+apiBasePath+`/openapi.json"></script>
				<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
			</body>
			</html>`)
	}
}

func SwaggerUI(apiBasePath string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.HTML(http.StatusOK, `<!DOCTYPE html>
			<html lang="en">
			<head>
			<meta charset="utf-8" />
			<meta name="viewport" content="width=device-width, initial-scale=1" />
			<meta name="description" content="SwaggerUI" />
			<title>SwaggerUI</title>
			<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
			</head>
			<body>
			<div id="swagger-ui"></div>
			<script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js" crossorigin></script>
			<script>
			window.onload = () => {
				window.ui = SwaggerUIBundle({
				url: '`+apiBasePath+`/openapi.json',
				dom_id: '#swagger-ui',
				});
			};
			</script>
			</body>
			</html>`)
	}
}
