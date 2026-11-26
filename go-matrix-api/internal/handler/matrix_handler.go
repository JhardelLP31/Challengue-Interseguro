package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/JhardelLP31/Challengue-Interseguro/go-matrix-api/internal/client"
	"github.com/JhardelLP31/Challengue-Interseguro/go-matrix-api/internal/model"
	"github.com/JhardelLP31/Challengue-Interseguro/go-matrix-api/internal/service"
)

// MatrixHandler expone los endpoints REST para factorización QR y rotación de matrices.
type MatrixHandler struct {
	matrixService *service.MatrixService
	statsClient   client.StatsClient
}

// NewMatrixHandler crea una instancia del controlador matricial con sus dependencias inyectadas.
func NewMatrixHandler(matrixService *service.MatrixService, statsClient client.StatsClient) *MatrixHandler {
	return &MatrixHandler{
		matrixService: matrixService,
		statsClient:   statsClient,
	}
}

// FactorizeQR recibe una matriz rectangular, calcula su factorización QR (Q y R),
// envía ambas matrices a la API de Node.js por HTTP y devuelve el resultado consolidado.
func (h *MatrixHandler) FactorizeQR(c *fiber.Ctx) error {
	var req model.MatrixRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "InvalidJson",
			Message: "El cuerpo de la petición debe ser un JSON válido con el campo 'matrix'.",
		})
	}

	rows, cols, err := h.matrixService.ValidateMatrix(req.Matrix)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "ValidationError",
			Message: err.Error(),
		})
	}

	qrResult, err := h.matrixService.ComputeQRFactorization(req.Matrix)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "CalculationError",
			Message: err.Error(),
		})
	}

	authHeader, _ := c.Locals("authHeader").(string)
	stats, err := h.statsClient.CalculateStatistics(c.UserContext(), authHeader, map[string]model.Matrix{
		"Q": qrResult.Q,
		"R": qrResult.R,
	})
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(model.ErrorResponse{
			Error:   "UpstreamServiceError",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.QRResponse{
		Operation: "qr_factorization",
		Input: model.MatrixInputMetadata{
			Rows:   rows,
			Cols:   cols,
			Matrix: req.Matrix,
		},
		QR:         qrResult,
		Statistics: stats,
	})
}

// RotateMatrix rota una matriz rectangular 90 grados (horario o antihorario),
// envía la matriz rotada a la API de Node.js por HTTP y devuelve matriz rotada + estadísticas.
func (h *MatrixHandler) RotateMatrix(c *fiber.Ctx) error {
	var req model.MatrixRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "InvalidJson",
			Message: "El cuerpo de la petición debe ser un JSON válido con el campo 'matrix'.",
		})
	}

	rows, cols, err := h.matrixService.ValidateMatrix(req.Matrix)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "ValidationError",
			Message: err.Error(),
		})
	}

	rotated, direction, err := h.matrixService.RotateMatrix(req.Matrix, req.Direction)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "ValidationError",
			Message: err.Error(),
		})
	}

	authHeader, _ := c.Locals("authHeader").(string)
	stats, err := h.statsClient.CalculateStatistics(c.UserContext(), authHeader, map[string]model.Matrix{
		"rotated": rotated,
	})
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(model.ErrorResponse{
			Error:   "UpstreamServiceError",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.RotateResponse{
		Operation: "matrix_rotation",
		Direction: direction,
		Input: model.MatrixInputMetadata{
			Rows:   rows,
			Cols:   cols,
			Matrix: req.Matrix,
		},
		Rotated:    rotated,
		Statistics: stats,
	})
}
