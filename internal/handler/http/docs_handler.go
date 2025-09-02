package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterDocsRoutes(r *gin.Engine) {
	// Serve the raw OpenAPI YAML
	r.StaticFile("/api/docs/openapi.yaml", "./docs/openapi.yaml")

	// Simple Swagger UI static page that loads the YAML from /api/docs/openapi.yaml
	r.GET("/api/docs", func(c *gin.Context) {
		html := `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>NewsBrief API Docs</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui.css" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui-bundle.js" crossorigin></script>
    <script>
      window.onload = () => {
        window.ui = SwaggerUIBundle({
          url: '/api/docs/openapi.yaml',
          dom_id: '#swagger-ui',
          presets: [SwaggerUIBundle.presets.apis],
          layout: 'BaseLayout'
        });
      };
    </script>
  </body>
</html>`
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})
}


