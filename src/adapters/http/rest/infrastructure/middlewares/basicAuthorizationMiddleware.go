package middlewares

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oswaldom-code/api-template-gin/pkg/config"
)

var BasicAuthorizationMiddleware MiddlewareFunc = func(c *gin.Context) {
	// get token from header
	token := c.GetHeader("Authorization")
	// validate token
	if token != "Basic "+base64.StdEncoding.EncodeToString([]byte(config.GetAuthenticationKey().Secret)) {
		// response unauthorized status code
		c.Redirect(http.StatusFound, "/authorization")
	}
}
