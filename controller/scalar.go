package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const scalarTemplate = `
<!doctype html>
<html>
  <head>
    <title>Scalar API Reference</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>

  <body>
    <div id="app"></div>

    <!-- Load the Script -->
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>

    <!-- Initialize the Scalar API Reference -->
    <script>
      Scalar.createApiReference('#app', {
        // The URL of the OpenAPI/Swagger document
        url: '/openapi.json',
        // Avoid CORS issues
        // proxyUrl: 'https://proxy.scalar.com',
      })
    </script>
  </body>
</html>
`

func Scalar(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(scalarTemplate))
}
