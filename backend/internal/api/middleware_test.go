package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCORSMiddleware_SetsHeaders(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := CORSMiddleware(next, "http://localhost:5173")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, "http://localhost:5173", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), "POST")
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCORSMiddleware_PreflightShortCircuits(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) { called = true })
	handler := CORSMiddleware(next, "http://localhost:5173")

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/calculate", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.False(t, called, "preflight must not reach the wrapped handler")
}

func TestHandleCalculate_MalformedJSON(t *testing.T) {
	rec := postCalculate(t, `{not json`)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	resp := decodeResponse(t, rec)
	assert.False(t, resp.Success)
	assert.Nil(t, resp.Result)
	assert.NotNil(t, resp.Error)
}
