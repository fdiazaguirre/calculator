// Package calculator holds the pure arithmetic logic for the service.
// It has no knowledge of HTTP, JSON, or transport concerns.
package calculator

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

// Supported operation identifiers.
const (
	OpAdd        = "add"
	OpSubtract   = "subtract"
	OpMultiply   = "multiply"
	OpDivide     = "divide"
	OpPower      = "power"
	OpSqrt       = "sqrt"
	OpPercentage = "percentage"
)

// SupportedOperations is the human-readable, comma-separated list used in error
// messages and shared with validation code.
const SupportedOperations = "add, subtract, multiply, divide, power, sqrt, percentage"

// MaxMagnitude bounds operands and results (FR-009). Anything with a larger
// absolute value is rejected.
const MaxMagnitude = 1e15

// significantDigits is the precision used to strip binary floating-point noise
// from results before they leave the service (FR-008).
const significantDigits = 12

// Error messages are drawn verbatim from the API contract error catalogue so
// backend and frontend can assert on the same strings.
var (
	ErrDivideByZero     = errors.New("Division by zero is not allowed")
	ErrResultOutOfRange = errors.New("Result out of range")
	ErrSqrtNegative     = errors.New("Square root of a negative number is not allowed")
	ErrNotReal          = errors.New("Result is not a real number")
)

// arities maps each known operation to the number of operands it consumes.
var arities = map[string]int{
	OpAdd:        2,
	OpSubtract:   2,
	OpMultiply:   2,
	OpDivide:     2,
	OpPower:      2,
	OpSqrt:       1,
	OpPercentage: 2,
}

// Arity reports how many operands an operation takes and whether the operation
// is known. Callers use this to validate requests before calling Compute.
func Arity(op string) (arity int, known bool) {
	arity, known = arities[op]
	return arity, known
}

// Compute performs the requested operation on the operands and returns the
// result rounded to significantDigits significant figures. Operands are assumed
// to have already passed range and arity validation at the transport boundary;
// Compute is responsible for math-domain errors only.
func Compute(op string, a, b float64) (float64, error) {
	switch op {
	case OpAdd:
		return finalize(a + b)
	case OpSubtract:
		return finalize(a - b)
	case OpMultiply:
		return finalize(a * b)
	case OpDivide:
		if b == 0 {
			return 0, ErrDivideByZero
		}
		return finalize(a / b)
	case OpPower:
		return power(a, b)
	case OpSqrt:
		if a < 0 {
			return 0, ErrSqrtNegative
		}
		return finalize(math.Sqrt(a))
	case OpPercentage:
		return finalize(a * b / 100)
	default:
		return 0, fmt.Errorf("Unsupported operation '%s'. Supported: %s", op, SupportedOperations)
	}
}

// power computes a^b, mapping the two domain edge cases to their contract
// errors: 0 raised to a negative exponent is division by zero, and a negative
// base with a fractional exponent has no real result.
func power(a, b float64) (float64, error) {
	if a == 0 && b < 0 {
		return 0, ErrDivideByZero
	}
	result := math.Pow(a, b)
	if math.IsNaN(result) {
		return 0, ErrNotReal
	}
	return finalize(result)
}

// finalize enforces the result range (FR-009) and strips float noise (FR-008).
func finalize(v float64) (float64, error) {
	if math.IsInf(v, 0) || math.IsNaN(v) || math.Abs(v) > MaxMagnitude {
		return 0, ErrResultOutOfRange
	}
	return roundResult(v), nil
}

// roundResult removes IEEE-754 representation noise by formatting the value to a
// fixed number of significant digits and parsing it back. For example this maps
// 0.30000000000000004 to 0.3.
func roundResult(v float64) float64 {
	s := strconv.FormatFloat(v, 'g', significantDigits, 64)
	rounded, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return v
	}
	return rounded
}
