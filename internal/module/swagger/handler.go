package swagger

import (
	"core-ledger/model/dto"
	"embed"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

//go:embed swagger.json
var swaggerJSON embed.FS

// SwaggerHandler handles Swagger documentation endpoints
type SwaggerHandler struct{}

// NewSwaggerHandler creates a new SwaggerHandler
func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

// GetSwaggerJSON serves the swagger.json file
func (h *SwaggerHandler) GetSwaggerJSON(c *gin.Context) {
	// Try to read from embedded file first
	data, err := swaggerJSON.ReadFile("swagger.json")
	if err != nil {
		// Fallback: try to read from root directory
		rootDir, _ := os.Getwd()
		swaggerPath := filepath.Join(rootDir, "swagger.json")
		data, err = os.ReadFile(swaggerPath)
		if err != nil {
			// Try relative path
			swaggerPath = "swagger.json"
			data, err = os.ReadFile(swaggerPath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, dto.PreResponse{
					Error: &dto.ResponseError{
						Message: "Failed to read swagger.json: " + err.Error(),
					},
				})
				return
			}
		}
	}

	c.Header("Content-Type", "application/json")
	c.Data(http.StatusOK, "application/json", data)
}

// GetSwaggerUI serves the Swagger UI HTML page
func (h *SwaggerHandler) GetSwaggerUI(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Core Ledger API Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui.css" />
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }
        *, *:before, *:after {
            box-sizing: inherit;
        }
        body {
            margin:0;
            background: #fafafa;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: "/api/v1/swagger/swagger.json",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                validatorUrl: null,
                tryItOutEnabled: true
            });
        };
    </script>
</body>
</html>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// GetSwaggerUIAlternative serves an alternative Swagger UI using ReDoc
func (h *SwaggerHandler) GetSwaggerUIAlternative(c *gin.Context) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Core Ledger API Documentation - ReDoc</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
    <style>
        body {
            margin: 0;
            padding: 0;
        }
    </style>
</head>
<body>
    <redoc spec-url="/api/v1/swagger/swagger.json"></redoc>
    <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
</body>
</html>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

