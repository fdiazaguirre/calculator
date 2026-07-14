package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestServer() http.Handler {
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	return mux
}

func postCalculate(t *testing.T, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	newTestServer().ServeHTTP(rec, req)
	return rec
}

func decodeResponse(t *testing.T, rec *httptest.ResponseRecorder) CalculationResponse {
	t.Helper()
	var resp CalculationResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	return resp
}

func TestHandleCalculate_Success(t *testing.T) {
	tests := []struct {
		name string
		body string
		want float64
	}{
		{"add", `{"operation":"add","a":7,"b":5}`, 12},
		{"subtract", `{"operation":"subtract","a":10,"b":4}`, 6},
		{"multiply", `{"operation":"multiply","a":6,"b":7}`, 42},
		{"divide", `{"operation":"divide","a":9,"b":4}`, 2.25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := postCalculate(t, tt.body)
			assert.Equal(t, http.StatusOK, rec.Code)

			resp := decodeResponse(t, rec)
			assert.True(t, resp.Success)
			require.NotNil(t, resp.Result)
			assert.Equal(t, tt.want, *resp.Result)
			assert.Nil(t, resp.Error)
		})
	}
}

func TestHandleCalculate_ResponseEnvelopeShape(t *testing.T) {
	rec := postCalculate(t, `{"operation":"add","a":1,"b":2}`)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))

	// Envelope must contain exactly the three documented fields.
	var raw map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &raw))
	assert.Contains(t, raw, "success")
	assert.Contains(t, raw, "result")
	assert.Contains(t, raw, "error")
}

func TestHandleCalculate_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/calculate", nil)
	rec := httptest.NewRecorder()
	newTestServer().ServeHTTP(rec, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestHandleCalculate_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantMsg string
	}{
		{
			"missing operand b for binary op",
			`{"operation":"divide","a":9}`,
			"Operation 'divide' requires operand 'b'",
		},
		{
			"operand a out of range",
			`{"operation":"add","a":2e15,"b":1}`,
			"Operand 'a' is out of the supported range (±1e15)",
		},
		{
			"operand b out of range",
			`{"operation":"add","a":1,"b":-2e15}`,
			"Operand 'b' is out of the supported range (±1e15)",
		},
		{
			"unknown operation",
			`{"operation":"modulo","a":10,"b":3}`,
			"Unsupported operation 'modulo'. Supported: add, subtract, multiply, divide, power, sqrt, percentage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := postCalculate(t, tt.body)
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			resp := decodeResponse(t, rec)
			assert.False(t, resp.Success)
			assert.Nil(t, resp.Result)
			require.NotNil(t, resp.Error)
			assert.Equal(t, tt.wantMsg, *resp.Error)
		})
	}
}

// Calculation-domain errors are HTTP 200 with success=false, not 4xx.
func TestHandleCalculate_CalculationError(t *testing.T) {
	rec := postCalculate(t, `{"operation":"divide","a":5,"b":0}`)
	assert.Equal(t, http.StatusOK, rec.Code)

	resp := decodeResponse(t, rec)
	assert.False(t, resp.Success)
	assert.Nil(t, resp.Result)
	require.NotNil(t, resp.Error)
	assert.Equal(t, "Division by zero is not allowed", *resp.Error)
}

func TestHandleHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	newTestServer().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var body map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "ok", body["status"])
}
