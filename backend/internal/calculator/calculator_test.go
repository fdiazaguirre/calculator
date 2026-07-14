package calculator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompute_BasicArithmetic(t *testing.T) {
	tests := []struct {
		name string
		op   string
		a    float64
		b    float64
		want float64
	}{
		{"addition", OpAdd, 7, 5, 12},
		{"subtraction", OpSubtract, 10, 4, 6},
		{"multiplication", OpMultiply, 6, 7, 42},
		{"division", OpDivide, 9, 4, 2.25},
		{"addition with negatives", OpAdd, -3, -9, -12},
		{"subtraction to negative", OpSubtract, 4, 10, -6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compute(tt.op, tt.a, tt.b)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// FR-008: exact decimal display up to 12 significant digits, no float noise.
func TestCompute_FloatingPointPrecision(t *testing.T) {
	got, err := Compute(OpAdd, 0.1, 0.2)
	require.NoError(t, err)
	assert.Equal(t, 0.3, got, "0.1 + 0.2 must display as 0.3, not 0.30000000000000004")
}

func TestArity(t *testing.T) {
	tests := []struct {
		op        string
		wantArity int
		wantKnown bool
	}{
		{OpAdd, 2, true},
		{OpDivide, 2, true},
		{OpPower, 2, true},
		{OpSqrt, 1, true},
		{OpPercentage, 2, true},
		{"modulo", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.op, func(t *testing.T) {
			arity, known := Arity(tt.op)
			assert.Equal(t, tt.wantKnown, known)
			assert.Equal(t, tt.wantArity, arity)
		})
	}
}

func TestCompute_AdvancedOperations(t *testing.T) {
	tests := []struct {
		name string
		op   string
		a    float64
		b    float64
		want float64
	}{
		{"power", OpPower, 2, 10, 1024},
		{"power zero exponent", OpPower, 0, 0, 1},
		{"power negative base integer exponent", OpPower, -8, 3, -512},
		{"square root", OpSqrt, 144, 0, 12},
		{"square root of zero", OpSqrt, 0, 0, 0},
		{"percentage", OpPercentage, 15, 200, 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compute(tt.op, tt.a, tt.b)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCompute_AdvancedErrors(t *testing.T) {
	tests := []struct {
		name    string
		op      string
		a       float64
		b       float64
		wantMsg string
	}{
		{"zero to negative power", OpPower, 0, -2, "Division by zero is not allowed"},
		{"negative base fractional exponent", OpPower, -8, 0.5, "Result is not a real number"},
		{"square root of negative", OpSqrt, -4, 0, "Square root of a negative number is not allowed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Compute(tt.op, tt.a, tt.b)
			require.Error(t, err)
			assert.EqualError(t, err, tt.wantMsg)
		})
	}
}

func TestCompute_Errors(t *testing.T) {
	tests := []struct {
		name    string
		op      string
		a       float64
		b       float64
		wantMsg string
	}{
		{"division by zero", OpDivide, 5, 0, "Division by zero is not allowed"},
		{"result out of range via overflow", OpMultiply, 1e300, 1e300, "Result out of range"},
		{"result out of range above 1e15", OpAdd, 9e14, 9e14, "Result out of range"},
		{"unknown operation", "modulo", 10, 3, "Unsupported operation 'modulo'. Supported: add, subtract, multiply, divide, power, sqrt, percentage"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Compute(tt.op, tt.a, tt.b)
			require.Error(t, err)
			assert.EqualError(t, err, tt.wantMsg)
		})
	}
}
