package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/JhardelLP31/Challengue-Interseguro/go-matrix-api/internal/config"
	"github.com/JhardelLP31/Challengue-Interseguro/go-matrix-api/internal/model"
)

// AuthHandler gestiona la emisión de tokens JWT para consumir ambas APIs.
type AuthHandler struct {
	cfg config.Config
}

// NewAuthHandler crea una nueva instancia de AuthHandler.
func NewAuthHandler(cfg config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

// IssueToken valida las credenciales del cliente y devuelve un JWT firmado con HS256.
func (h *AuthHandler) IssueToken(c *fiber.Ctx) error {
	var req model.TokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "InvalidJson",
			Message: "El cuerpo de la petición no es un JSON válido.",
		})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "ValidationError",
			Message: "Los campos username y password son requeridos.",
		})
	}

	if req.Username != h.cfg.AuthUsername || req.Password != h.cfg.AuthPassword {
		return c.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Credenciales inválidas.",
		})
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  req.Username,
		"role": "api-client",
		"iat":  now.Unix(),
		"exp":  now.Add(h.cfg.JWTExpiresIn).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error:   "InternalServerError",
			Message: "No se pudo generar el token JWT.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.TokenResponse{
		Token:     signedToken,
		TokenType: "Bearer",
		ExpiresIn: int64(h.cfg.JWTExpiresIn.Seconds()),
	})
}
