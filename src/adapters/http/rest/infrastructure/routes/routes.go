package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/oswaldom-code/api-template-gin/src/adapters/http/rest/infrastructure/middlewares"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	Health(c *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []middlewares.MiddlewareFunc
}

// Health operation middleware
func (siw *ServerInterfaceWrapper) Health(c *gin.Context) {
	siw.Handler.Health(c)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL     string
	Middlewares []middlewares.MiddlewareFunc
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router *gin.Engine, si ServerInterface) *gin.Engine {
	return RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router *gin.Engine, si ServerInterface, options GinServerOptions) *gin.Engine {
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
	}

	router.GET(options.BaseURL+"/health", wrapper.Health)

	return router
}
