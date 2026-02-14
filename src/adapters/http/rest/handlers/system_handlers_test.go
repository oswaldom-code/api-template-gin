package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oswaldom-code/api-template-gin/src/adapters/http/rest/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPing_Returns200WithEnvelopeFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	handler := NewRestHandler()
	router.GET("/ping", handler.Ping)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	dataMap := resp.Data.(map[string]interface{})
	assert.Equal(t, "pong", dataMap["ping"])

	assert.NotEmpty(t, resp.Meta.Timestamp)
}

func TestPing_DoesNotContainOldFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	handler := NewRestHandler()
	router.GET("/ping", handler.Ping)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var raw map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &raw)
	require.NoError(t, err)
}
