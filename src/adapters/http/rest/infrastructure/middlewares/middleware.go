package middlewares

import "github.com/gin-gonic/gin"

type MiddlewareFunc func(c *gin.Context)
