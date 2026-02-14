package infrastructure

import (
	"encoding/base64"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	metrics "github.com/penglongli/gin-metrics/ginmetrics"

	"github.com/oswaldom-code/api-template-gin/pkg/config"
	"github.com/oswaldom-code/api-template-gin/src/adapters/http/rest/dto"
	"github.com/oswaldom-code/api-template-gin/src/adapters/http/rest/handlers"
)

var basicAuthorizationMiddleware gin.HandlerFunc = func(c *gin.Context) {
	// get token from header
	token := c.GetHeader("Authorization")
	expected := "Basic " + base64.StdEncoding.EncodeToString([]byte(config.GetAuthenticationKey().Secret))
	// validate token
	if token != expected {
		dto.Unauthorized(c, "Invalid or missing auth token")
		return
	}
	c.Next()
}

func setMetrics(router *gin.Engine) {
	// get global Monitor object
	monitor := metrics.GetMonitor()
	// +optional set metric path, default /debug/metrics
	monitor.SetMetricPath("/metrics")
	// +optional set slow time, default 5s
	monitor.SetSlowTime(10)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	monitor.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	monitor.Use(router)
}

// NewGinServer creates a new Gin server.
func NewGinServer(handler ServerInterface) *gin.Engine {
	// get configuration
	serverConfig := config.GetServerConfig()
	// validate parameters configuration
	if err := serverConfig.Validate(); err != nil {
		panic("[ERROR] server configuration is not valid")
	}

	// set gin mode (debug or release)
	gin.SetMode(serverConfig.Mode)
	// if debug mode is release write the logs to a file
	if serverConfig.Mode == "release" {
		// Disable Console Color, you don't need console color when writing the logs to file.
		gin.DisableConsoleColor()
		// Logging to a file.
		f, err := os.Create("log/error.log")
		if err != nil {
			panic("[ERROR] error creating log file")
		}
		gin.DefaultWriter = io.MultiWriter(f)
	}

	// create routes
	router := gin.Default()
	// set metrics
	setMetrics(router)
	// register handlers with route groups (public + protected)
	ginServerOptions := GinServerOptions{
		BaseURL:     "/",
		Middlewares: []gin.HandlerFunc{basicAuthorizationMiddleware},
	}
	RegisterHandlersWithOptions(router, handler, ginServerOptions)
	return router
}

func loadHandlers() *handlers.Handler {
	return handlers.NewRestHandler()
}

func NewServer() *gin.Engine {
	return NewGinServer(
		loadHandlers(),
	)
}
