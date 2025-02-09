package infrastructure

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	metrics "github.com/penglongli/gin-metrics/ginmetrics"

	"github.com/oswaldom-code/api-template-gin/pkg/config"
	"github.com/oswaldom-code/api-template-gin/pkg/log"
	"github.com/oswaldom-code/api-template-gin/src/adapters/http/rest/handlers"
	"github.com/oswaldom-code/api-template-gin/src/adapters/http/rest/infrastructure/middlewares"
	"github.com/oswaldom-code/api-template-gin/src/adapters/http/rest/infrastructure/routes"
)

func setMetrics(router *gin.Engine) {
	monitor := metrics.GetMonitor()
	if monitor == nil {
		log.Error("[ERROR] No se pudo obtener el monitor de métricas")
		return
	}

	monitor.SetMetricPath("/metrics")
	slowTime := int32(5)
	monitor.SetSlowTime(slowTime) // default is 5 seconds TODO: make it configurable
	durationThresholds := []float64{0.5, 1, 3, 5, 10}
	monitor.SetDuration(durationThresholds)

	monitor.Use(router)

	log.Info("Métricas configuradas correctamente", log.Fields{
		"metric_path":          "/metrics",
		"slow_time":            slowTime,
		"duration_percentiles": durationThresholds,
	})
}

// NewGinServer creates a new Gin server.
func NewGinServer(handler routes.ServerInterface) *gin.Engine {
	// get configuration
	serverConfig := config.GetServerConfig()
	// validate parameters configuration
	if ok := serverConfig.Validate(); ok != nil {
		panic("[ERROR] server configuration is not valid")
	}

	log.Info("Running in mode: ", log.Fields{"mode": serverConfig.Mode})
	gin.SetMode(serverConfig.Mode)

	// if debug mode is release write the logs to a file
	if serverConfig.Mode == "release" {
		// Disable Console Color, you don't need console color when writing the logs to file.
		gin.DisableConsoleColor()

		if err := os.MkdirAll("log", os.ModePerm); err != nil {
			log.Error("[ERROR] failed to create log directory", log.Fields{"error": err})
			return nil
		}

		f, ok := os.Create("log/error.log")
		if ok != nil {
			log.Error("[ERROR] error creating log file", log.Fields{"error": ok})
			return nil
		}
		defer f.Close()

		// Usa un writer para los logs de Gin y tu paquete de logs
		multiWriter := io.MultiWriter(f, os.Stdout) // También puedes redirigir a stdout o a tu logger personalizado
		gin.DefaultWriter = multiWriter
	}

	// create routes
	router := gin.Default()
	if router == nil {
		log.Error("[ERROR] failed to create gin router")
		return nil
	}

	setMetrics(router)

	ginServerOptions := routes.GinServerOptions{BaseURL: "/"}
	ginServerOptions.Middlewares = append(ginServerOptions.Middlewares, middlewares.BasicAuthorizationMiddleware)

	// register Handler, router and middleware to gin
	routes.RegisterHandlersWithOptions(router, handler, ginServerOptions)

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
