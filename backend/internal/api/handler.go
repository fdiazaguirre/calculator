package api

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"

	"calculator/backend/internal/calculator"
)

// maxRequestBytes caps the calculation request body. A valid request is a few
// dozen bytes; this rejects oversized payloads before decoding them.
const maxRequestBytes = 4096

// CalculationRequest is the JSON body accepted by POST /api/v1/calculate.
// B is a pointer so a missing operand is distinguishable from a zero operand.
type CalculationRequest struct {
	Operation string   `json:"operation"`
	A         float64  `json:"a"`
	B         *float64 `json:"b"`
}

// CalculationResponse is the envelope returned for every calculation outcome.
// Exactly one of Result / Error is non-nil (FR-007).
type CalculationResponse struct {
	Success bool     `json:"success"`
	Result  *float64 `json:"result"`
	Error   *string  `json:"error"`
}

func successResponse(result float64) CalculationResponse {
	return CalculationResponse{Success: true, Result: &result, Error: nil}
}

func errorResponse(message string) CalculationResponse {
	return CalculationResponse{Success: false, Result: nil, Error: &message}
}

// HandleCalculate decodes a calculation request, validates it, runs it, and
// writes the response envelope. Structural problems (bad JSON, unknown
// operation, wrong arity, out-of-range operands) are HTTP 400. Math-domain
// failures (division by zero, out-of-range result) are HTTP 200 with
// success=false, because the request itself was well-formed.
func HandleCalculate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBytes)

	var req CalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse("Request body must be valid JSON"))
		return
	}

	if msg, ok := validate(req); !ok {
		writeJSON(w, http.StatusBadRequest, errorResponse(msg))
		return
	}

	b := 0.0
	if req.B != nil {
		b = *req.B
	}

	result, err := calculator.Compute(req.Operation, req.A, b)
	if err != nil {
		writeJSON(w, http.StatusOK, errorResponse(err.Error()))
		return
	}

	writeJSON(w, http.StatusOK, successResponse(result))
}

// validate checks the structural correctness of a request. It returns a
// user-facing message and false on the first problem found.
func validate(req CalculationRequest) (string, bool) {
	arity, known := calculator.Arity(req.Operation)
	if !known {
		return fmt.Sprintf("Unsupported operation '%s'. Supported: %s", req.Operation, calculator.SupportedOperations), false
	}

	if msg, ok := checkRange("a", req.A); !ok {
		return msg, false
	}
	if req.B != nil {
		if msg, ok := checkRange("b", *req.B); !ok {
			return msg, false
		}
	}

	switch arity {
	case 2:
		if req.B == nil {
			return fmt.Sprintf("Operation '%s' requires operand 'b'", req.Operation), false
		}
	case 1:
		if req.B != nil {
			return fmt.Sprintf("Operation '%s' takes one operand; 'b' must be omitted or null", req.Operation), false
		}
	}

	return "", true
}

func checkRange(name string, value float64) (string, bool) {
	// JSON has no NaN/Inf literals, so those guards are unreachable via the HTTP
	// path; they defend against any future non-JSON caller of validate.
	if math.IsNaN(value) || math.IsInf(value, 0) || math.Abs(value) > calculator.MaxMagnitude {
		return fmt.Sprintf("Operand '%s' is out of the supported range (±1e15)", name), false
	}
	return "", true
}

// HandleHealth reports service liveness for probes.
func HandleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	// The status header is already written, so an encode failure cannot change
	// the HTTP status — logging is the only recovery.
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("writeJSON encode error: %v", err)
	}
}
