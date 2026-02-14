package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/oswaldom-code/api-template-gin/src/adapters/http/rest/dto"
)

func (h *Handler) Ping(c *gin.Context) {
	dto.OK(c, gin.H{"ping": "pong"})
}
