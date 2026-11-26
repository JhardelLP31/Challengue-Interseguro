package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/JhardelLP31/Challengue-Interseguro/go-matrix-api/internal/model"
)

// JWTProtected valida que la petición incluya un token JWT HS256 válido en el header Authorization.
func JWTProtected(jwtSecret string) fiber.Handler {
	secretBytes := []byte(jwtSecret)

	return func(c *fiber.Ctx) error {
		authHeader := strings.TrimSpace(c.Get("Authorization"))
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Falta el header Authorization con formato Bearer <token>.",
			})
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse{
				Error:   "Unauthorized",
				Message: "El token JWT no puede estar vacío.",
			})
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado: %v", t.Header["alg"])
			}
			return secretBytes, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Token JWT inválido o expirado.",
			})
		}

		c.Locals("authHeader", authHeader)
		return c.Next()
	}
}
