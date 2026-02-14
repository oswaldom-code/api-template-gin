package infrastructure

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oswaldom-code/api-template-gin/src/adapters/http/rest/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicAuthMiddleware_Unauthorized_NoHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(basicAuthorizationMiddleware)
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, dto.ErrUnauthorized, resp.Error.Code)
	assert.NotEmpty(t, resp.Error.Message)
}

func TestBasicAuthMiddleware_Unauthorized_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(basicAuthorizationMiddleware)
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Basic invalid_token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, dto.ErrUnauthorized, resp.Error.Code)
}

func TestBasicAuthMiddleware_Unauthorized_ResponseDoesNotContainOldFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(basicAuthorizationMiddleware)
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var raw map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &raw)
	require.NoError(t, err)
}

func TestBasicAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	os.Setenv("auth.secret", "test_secret")
	defer os.Unsetenv("auth.secret")

	router := gin.New()
	router.Use(basicAuthorizationMiddleware)
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	token := "Basic " + base64.StdEncoding.EncodeToString([]byte("test_secret"))
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "true")
}
