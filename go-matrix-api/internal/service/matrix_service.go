package service

import (
	"errors"
	"math"
	"strings"

	"github.com/JhardelLP31/Challengue-Interseguro/go-matrix-api/internal/model"
)

const zeroEpsilon = 1e-12

// Errores de validación de entrada para operaciones matriciales.
var (
	ErrEmptyMatrix     = errors.New("la matriz no puede estar vacía")
	ErrEmptyColumn     = errors.New("la matriz debe tener al menos una columna")
	ErrNonRectangular  = errors.New("todas las filas de la matriz deben tener la misma longitud")
	ErrNonFiniteNumber = errors.New("la matriz contiene valores numéricos inválidos (NaN o Inf)")
	ErrInvalidRotation = errors.New("dirección de rotación inválida; use 'clockwise' o 'counterclockwise'")
)

// MatrixService encapsula las operaciones matemáticas y validaciones sobre matrices.
type MatrixService struct{}

// NewMatrixService crea una nueva instancia del servicio matricial.
func NewMatrixService() *MatrixService {
	return &MatrixService{}
}

// ValidateMatrix verifica que la entrada sea una matriz rectangular no vacía de números finitos
// y retorna sus dimensiones (filas, columnas).
func (s *MatrixService) ValidateMatrix(matrix model.Matrix) (int, int, error) {
	rows := len(matrix)
	if rows == 0 {
		return 0, 0, ErrEmptyMatrix
	}

	cols := len(matrix[0])
	if cols == 0 {
		return 0, 0, ErrEmptyColumn
	}

	for _, row := range matrix {
		if len(row) != cols {
			return 0, 0, ErrNonRectangular
		}
		for _, val := range row {
			if math.IsNaN(val) || math.IsInf(val, 0) {
				return 0, 0, ErrNonFiniteNumber
			}
		}
	}

	return rows, cols, nil
}

// ComputeQRFactorization calcula la factorización QR completa de una matriz rectangular A (m x n)
// utilizando reflexiones de Householder: A = Q * R, donde Q es ortogonal (m x m) y R es
// triangular/trapezoidal superior (m x n).
func (s *MatrixService) ComputeQRFactorization(matrix model.Matrix) (model.QRResult, error) {
	rows, cols, err := s.ValidateMatrix(matrix)
	if err != nil {
		return model.QRResult{}, err
	}

	r := cloneMatrix(matrix)
	q := identityMatrix(rows)

	maxStep := rows - 1
	if cols < maxStep {
		maxStep = cols
	}

	for k := 0; k < maxStep; k++ {
		subLen := rows - k
		v := make([]float64, subLen)

		normSq := 0.0
		tailNormSq := 0.0
		for i := 0; i < subLen; i++ {
			val := r[k+i][k]
			v[i] = val
			normSq += val * val
			if i > 0 {
				tailNormSq += val * val
			}
		}

		// Si toda la subcolumna es cero o ya está en forma triangular con pivote no negativo,
		// la reflexión es la matriz identidad y no requiere transformación.
		if normSq < 1e-24 || (tailNormSq < 1e-24 && v[0] >= 0) {
			continue
		}

		normX := math.Sqrt(normSq)
		// Elección numéricamente estable del signo de alpha para evitar cancelación catastrófica.
		alpha := -math.Copysign(normX, v[0])
		v[0] -= alpha

		vDotV := 0.0
		for i := 0; i < subLen; i++ {
			vDotV += v[i] * v[i]
		}
		if vDotV < 1e-24 {
			continue
		}

		// Aplicar H_k por la izquierda a R: R[k:, k:] = (I - 2*v*v^T / vDotV) * R[k:, k:]
		for j := k; j < cols; j++ {
			dot := 0.0
			for i := 0; i < subLen; i++ {
				dot += v[i] * r[k+i][j]
			}
			scale := 2.0 * dot / vDotV
			for i := 0; i < subLen; i++ {
				r[k+i][j] -= scale * v[i]
			}
		}

		// Acumular H_k en Q por la derecha: Q[:, k:] = Q[:, k:] * H_k (H_k es simétrica)
		for row := 0; row < rows; row++ {
			dot := 0.0
			for i := 0; i < subLen; i++ {
				dot += q[row][k+i] * v[i]
			}
			scale := 2.0 * dot / vDotV
			for i := 0; i < subLen; i++ {
				q[row][k+i] -= scale * v[i]
			}
		}
	}

	cleanSmallValues(q)
	cleanSmallValues(r)

	return model.QRResult{
		Q: q,
		R: r,
	}, nil
}

// RotateMatrix rota una matriz rectangular 90 grados en sentido horario ("clockwise")
// o antihorario ("counterclockwise").
func (s *MatrixService) RotateMatrix(matrix model.Matrix, direction string) (model.Matrix, string, error) {
	rows, cols, err := s.ValidateMatrix(matrix)
	if err != nil {
		return nil, "", err
	}

	normalizedDir := strings.ToLower(strings.TrimSpace(direction))
	if normalizedDir == "" {
		normalizedDir = "clockwise"
	}

	if normalizedDir != "clockwise" && normalizedDir != "counterclockwise" {
		return nil, "", ErrInvalidRotation
	}

	rotated := make(model.Matrix, cols)
	for i := 0; i < cols; i++ {
		rotated[i] = make([]float64, rows)
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			val := normalizeZero(matrix[i][j])
			if normalizedDir == "clockwise" {
				rotated[j][rows-1-i] = val
			} else {
				rotated[cols-1-j][i] = val
			}
		}
	}

	return rotated, normalizedDir, nil
}

func cloneMatrix(matrix model.Matrix) model.Matrix {
	cloned := make(model.Matrix, len(matrix))
	for i, row := range matrix {
		cloned[i] = make([]float64, len(row))
		copy(cloned[i], row)
	}
	return cloned
}

func identityMatrix(size int) model.Matrix {
	id := make(model.Matrix, size)
	for i := 0; i < size; i++ {
		id[i] = make([]float64, size)
		id[i][i] = 1.0
	}
	return id
}

func cleanSmallValues(matrix model.Matrix) {
	for i := range matrix {
		for j := range matrix[i] {
			if math.Abs(matrix[i][j]) < zeroEpsilon {
				matrix[i][j] = 0.0
			} else {
				matrix[i][j] = normalizeZero(matrix[i][j])
			}
		}
	}
}

func normalizeZero(v float64) float64 {
	if v == 0 {
		return 0.0
	}
	return v
}
