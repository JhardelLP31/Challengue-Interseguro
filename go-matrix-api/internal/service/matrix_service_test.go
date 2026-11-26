package service

import (
	"errors"
	"math"
	"testing"

	"github.com/JhardelLP31/Challengue-Interseguro/go-matrix-api/internal/model"
)

func multiplyMatrices(a, b model.Matrix) model.Matrix {
	rows := len(a)
	cols := len(b[0])
	inner := len(b)
	res := make(model.Matrix, rows)
	for i := 0; i < rows; i++ {
		res[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			sum := 0.0
			for k := 0; k < inner; k++ {
				sum += a[i][k] * b[k][j]
			}
			res[i][j] = sum
		}
	}
	return res
}

func transposeMatrix(m model.Matrix) model.Matrix {
	rows := len(m)
	cols := len(m[0])
	t := make(model.Matrix, cols)
	for i := 0; i < cols; i++ {
		t[i] = make([]float64, rows)
		for j := 0; j < rows; j++ {
			t[i][j] = m[j][i]
		}
	}
	return t
}

func assertMatricesAlmostEqual(t *testing.T, expected, actual model.Matrix, tol float64) {
	t.Helper()
	if len(expected) != len(actual) || len(expected[0]) != len(actual[0]) {
		t.Fatalf("dimensiones distintas: esperado %dx%d, obtenido %dx%d",
			len(expected), len(expected[0]), len(actual), len(actual[0]))
	}
	for i := range expected {
		for j := range expected[i] {
			if math.Abs(expected[i][j]-actual[i][j]) > tol {
				t.Fatalf("diferencia en (%d,%d): esperado %v, obtenido %v", i, j, expected[i][j], actual[i][j])
			}
		}
	}
}

func assertQRProperties(t *testing.T, input model.Matrix, qr model.QRResult) {
	t.Helper()
	rows := len(input)
	cols := len(input[0])

	if len(qr.Q) != rows || len(qr.Q[0]) != rows {
		t.Fatalf("Q debe ser %dx%d, obtenido %dx%d", rows, rows, len(qr.Q), len(qr.Q[0]))
	}
	if len(qr.R) != rows || len(qr.R[0]) != cols {
		t.Fatalf("R debe ser %dx%d, obtenido %dx%d", rows, cols, len(qr.R), len(qr.R[0]))
	}

	// 1. Verificar A = Q * R
	reconstructed := multiplyMatrices(qr.Q, qr.R)
	assertMatricesAlmostEqual(t, input, reconstructed, 1e-9)

	// 2. Verificar ortogonalidad: Q^T * Q = I
	qtq := multiplyMatrices(transposeMatrix(qr.Q), qr.Q)
	assertMatricesAlmostEqual(t, identityMatrix(rows), qtq, 1e-9)

	// 3. Verificar que R sea triangular/trapezoidal superior (ceros debajo de la diagonal)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if i > j && math.Abs(qr.R[i][j]) > 1e-9 {
				t.Fatalf("R[%d][%d] debería ser 0, obtenido %v", i, j, qr.R[i][j])
			}
		}
	}
}

func TestComputeQRFactorization_TallRectangularMatrix(t *testing.T) {
	svc := NewMatrixService()
	input := model.Matrix{
		{1, 2},
		{3, 4},
		{5, 6},
	}

	qr, err := svc.ComputeQRFactorization(input)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	assertQRProperties(t, input, qr)
}

func TestComputeQRFactorization_WideRectangularMatrix(t *testing.T) {
	svc := NewMatrixService()
	input := model.Matrix{
		{1, 2, 3},
		{4, 5, 6},
	}

	qr, err := svc.ComputeQRFactorization(input)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	assertQRProperties(t, input, qr)
}

func TestComputeQRFactorization_DiagonalMatrixPreservesDiagonality(t *testing.T) {
	svc := NewMatrixService()
	input := model.Matrix{
		{3, 0},
		{0, 5},
	}

	qr, err := svc.ComputeQRFactorization(input)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	assertQRProperties(t, input, qr)
	assertMatricesAlmostEqual(t, identityMatrix(2), qr.Q, 1e-9)
	assertMatricesAlmostEqual(t, input, qr.R, 1e-9)
}

func TestRotateMatrix_ClockwiseAndCounterclockwise(t *testing.T) {
	svc := NewMatrixService()
	input := model.Matrix{
		{1, 2, 3},
		{4, 5, 6},
	}

	cw, dir, err := svc.RotateMatrix(input, "")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if dir != "clockwise" {
		t.Fatalf("dirección esperada clockwise, obtenida %s", dir)
	}
	expectedCW := model.Matrix{
		{4, 1},
		{5, 2},
		{6, 3},
	}
	assertMatricesAlmostEqual(t, expectedCW, cw, 1e-12)

	ccw, dir, err := svc.RotateMatrix(input, "counterclockwise")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if dir != "counterclockwise" {
		t.Fatalf("dirección esperada counterclockwise, obtenida %s", dir)
	}
	expectedCCW := model.Matrix{
		{3, 6},
		{2, 5},
		{1, 4},
	}
	assertMatricesAlmostEqual(t, expectedCCW, ccw, 1e-12)
}

func TestValidateMatrix_Errors(t *testing.T) {
	svc := NewMatrixService()

	if _, _, err := svc.ValidateMatrix(model.Matrix{}); !errors.Is(err, ErrEmptyMatrix) {
		t.Fatalf("esperado ErrEmptyMatrix, obtenido %v", err)
	}
	if _, _, err := svc.ValidateMatrix(model.Matrix{{}}); !errors.Is(err, ErrEmptyColumn) {
		t.Fatalf("esperado ErrEmptyColumn, obtenido %v", err)
	}
	if _, _, err := svc.ValidateMatrix(model.Matrix{{1, 2}, {3}}); !errors.Is(err, ErrNonRectangular) {
		t.Fatalf("esperado ErrNonRectangular, obtenido %v", err)
	}
	if _, _, err := svc.ValidateMatrix(model.Matrix{{math.NaN()}}); !errors.Is(err, ErrNonFiniteNumber) {
		t.Fatalf("esperado ErrNonFiniteNumber, obtenido %v", err)
	}
}
