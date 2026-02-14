package infrastructure

import "github.com/gin-gonic/gin"

// ServerInterface represents all server handlers.
type ServerInterface interface {
	Ping(c *gin.Context)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL     string
	Middlewares []gin.HandlerFunc
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router *gin.Engine, si ServerInterface) *gin.Engine {
	return RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with public and protected route groups.
func RegisterHandlersWithOptions(router *gin.Engine, si ServerInterface, options GinServerOptions) *gin.Engine {
	public := router.Group(options.BaseURL)
	{
		public.GET("/ping", func(c *gin.Context) {
			si.Ping(c)
		})
	}

	protected := router.Group(options.BaseURL)
	for _, m := range options.Middlewares {
		protected.Use(m)
	}
	{
		// Add protected routes here as the API grows
		// Example: protected.GET("/users", si.ListUsers)
	}

	return router
}
