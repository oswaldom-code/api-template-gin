package dto

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestSuccessResponse_JSONSerialization(t *testing.T) {
	resp := SuccessResponse{
		Data: map[string]string{"key": "value"},
		Meta: Meta{Timestamp: "2026-01-01T00:00:00Z"},
	}

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.NotNil(t, result["data"])
	assert.NotNil(t, result["meta"])

	dataMap := result["data"].(map[string]interface{})
	assert.Equal(t, "value", dataMap["key"])

	metaMap := result["meta"].(map[string]interface{})
	assert.Equal(t, "2026-01-01T00:00:00Z", metaMap["timestamp"])
}

func TestErrorResponse_JSONSerialization(t *testing.T) {
	resp := ErrorResponse{
		Error: ErrorDetail{
			Code:    ErrUnauthorized,
			Message: "Invalid token",
		},
	}

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	errObj := result["error"].(map[string]interface{})
	assert.Equal(t, "UNAUTHORIZED", errObj["code"])
	assert.Equal(t, "Invalid token", errObj["message"])
}

func TestOK_Helper(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	OK(c, gin.H{"ping": "pong"})

	assert.Equal(t, http.StatusOK, w.Code)

	var resp SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	dataMap := resp.Data.(map[string]interface{})
	assert.Equal(t, "pong", dataMap["ping"])
	assert.NotEmpty(t, resp.Meta.Timestamp)
}

func TestCreated_Helper(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Created(c, gin.H{"id": "123"})

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	dataMap := resp.Data.(map[string]interface{})
	assert.Equal(t, "123", dataMap["id"])
}

func TestBadRequest_Helper(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	BadRequest(c, "invalid input")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, ErrBadRequest, resp.Error.Code)
	assert.Equal(t, "invalid input", resp.Error.Message)
}

func TestUnauthorized_Helper(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Unauthorized(c, "missing token")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.True(t, c.IsAborted(), "Unauthorized must abort the context")

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, ErrUnauthorized, resp.Error.Code)
	assert.Equal(t, "missing token", resp.Error.Message)
}

func TestNotFound_Helper(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	NotFound(c, "resource not found")

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, ErrNotFound, resp.Error.Code)
}

func TestInternalError_Helper(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	InternalError(c, "unexpected failure")

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, ErrInternalServer, resp.Error.Code)
	assert.Equal(t, "unexpected failure", resp.Error.Message)
}

func TestSuccess_GenericStatusCode(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Success(c, http.StatusAccepted, gin.H{"status": "processing"})

	assert.Equal(t, http.StatusAccepted, w.Code)
}

func TestError_GenericStatusCode(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Error(c, http.StatusConflict, ErrConflict, "resource already exists")

	assert.Equal(t, http.StatusConflict, w.Code)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, ErrConflict, resp.Error.Code)
}

func TestAbortWithError_AbortsContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	AbortWithError(c, http.StatusForbidden, ErrForbidden, "access denied")

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusForbidden, w.Code)
}
