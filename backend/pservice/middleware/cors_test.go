package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCORSConfigFromString(t *testing.T) {
	t.Parallel()

	t.Run("empty means allow none", func(t *testing.T) {
		t.Parallel()
		cfg, err := ParseCORSConfigFromString("")
		require.NoError(t, err)
		assert.True(t, cfg.AllowNone)
		assert.False(t, cfg.IsOriginAllowed("http://localhost:5173"))
	})

	t.Run("star means allow all", func(t *testing.T) {
		t.Parallel()
		cfg, err := ParseCORSConfigFromString("*")
		require.NoError(t, err)
		assert.True(t, cfg.AllowAll)
		assert.True(t, cfg.IsOriginAllowed("https://example.com"))
	})

	t.Run("specific origins are trimmed and validated", func(t *testing.T) {
		t.Parallel()
		cfg, err := ParseCORSConfigFromString(" http://localhost:5173,https://example.com ")
		require.NoError(t, err)
		assert.True(t, cfg.IsOriginAllowed("http://localhost:5173"))
		assert.True(t, cfg.IsOriginAllowed("https://example.com"))
		assert.False(t, cfg.IsOriginAllowed("https://other.example.com"))
	})

	t.Run("path is rejected", func(t *testing.T) {
		t.Parallel()
		_, err := ParseCORSConfigFromString("https://example.com/app")
		require.Error(t, err)
	})
}

func setupTestRouter(config *CORSConfig) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(config.Middleware())
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	r.POST("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	return r
}

func TestCORSMiddleware(t *testing.T) {
	t.Run("no origin passes without CORS headers", func(t *testing.T) {
		cfg, err := ParseCORSConfigFromString("http://localhost:5173")
		require.NoError(t, err)
		r := setupTestRouter(cfg)
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("allow all simple request", func(t *testing.T) {
		cfg, err := ParseCORSConfigFromString("*")
		require.NoError(t, err)
		r := setupTestRouter(cfg)
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set("Origin", "https://example.com")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "POST, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	})

	t.Run("allow all preflight", func(t *testing.T) {
		cfg, err := ParseCORSConfigFromString("*")
		require.NoError(t, err)
		r := setupTestRouter(cfg)
		req := httptest.NewRequest(http.MethodOptions, "/test", nil)
		req.Header.Set("Origin", "https://example.com")
		req.Header.Set("Access-Control-Request-Method", "POST")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("allow none blocks preflight only", func(t *testing.T) {
		cfg, err := ParseCORSConfigFromString("")
		require.NoError(t, err)
		r := setupTestRouter(cfg)
		simpleReq := httptest.NewRequest(http.MethodPost, "/test", nil)
		simpleReq.Header.Set("Origin", "https://example.com")
		simpleW := httptest.NewRecorder()
		r.ServeHTTP(simpleW, simpleReq)

		preflightReq := httptest.NewRequest(http.MethodOptions, "/test", nil)
		preflightReq.Header.Set("Origin", "https://example.com")
		preflightReq.Header.Set("Access-Control-Request-Method", "POST")
		preflightW := httptest.NewRecorder()
		r.ServeHTTP(preflightW, preflightReq)

		assert.Equal(t, http.StatusOK, simpleW.Code)
		assert.Empty(t, simpleW.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, http.StatusForbidden, preflightW.Code)
	})

	t.Run("specific origin reflects allowed origin and blocks unknown preflight", func(t *testing.T) {
		cfg, err := ParseCORSConfigFromString("http://localhost:5173")
		require.NoError(t, err)
		r := setupTestRouter(cfg)
		allowedReq := httptest.NewRequest(http.MethodPost, "/test", nil)
		allowedReq.Header.Set("Origin", "http://localhost:5173")
		allowedW := httptest.NewRecorder()
		r.ServeHTTP(allowedW, allowedReq)

		blockedReq := httptest.NewRequest(http.MethodOptions, "/test", nil)
		blockedReq.Header.Set("Origin", "https://example.com")
		blockedReq.Header.Set("Access-Control-Request-Method", "POST")
		blockedW := httptest.NewRecorder()
		r.ServeHTTP(blockedW, blockedReq)

		assert.Equal(t, http.StatusOK, allowedW.Code)
		assert.Equal(t, "http://localhost:5173", allowedW.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "Origin", allowedW.Header().Get("Vary"))
		assert.Equal(t, http.StatusForbidden, blockedW.Code)
	})
}
