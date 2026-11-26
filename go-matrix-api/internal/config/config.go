package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config centraliza la configuración de la API en Go cargada desde variables de entorno.
type Config struct {
	Port           string
	StatsAPIURL    string
	JWTSecret      string
	JWTExpiresIn   time.Duration
	AuthUsername   string
	AuthPassword   string
	HTTPTimeout    time.Duration
}

// Load construye la configuración leyendo variables de entorno con valores por defecto para desarrollo.
func Load() Config {
	expiresHours := parsePositiveInt(os.Getenv("JWT_EXPIRES_IN_HOURS"), 24)
	timeoutSeconds := parsePositiveInt(os.Getenv("HTTP_TIMEOUT_SECONDS"), 10)

	return Config{
		Port:         getEnvOrDefault("PORT", "8080"),
		StatsAPIURL:  strings.TrimRight(getEnvOrDefault("STATS_API_URL", "http://localhost:3000"), "/"),
		JWTSecret:    getEnvOrDefault("JWT_SECRET", "interseguro-coding-challenge-super-secret-key-2024"),
		JWTExpiresIn: time.Duration(expiresHours) * time.Hour,
		AuthUsername: getEnvOrDefault("AUTH_USERNAME", "interseguro"),
		AuthPassword: getEnvOrDefault("AUTH_PASSWORD", "challenge2024"),
		HTTPTimeout:  time.Duration(timeoutSeconds) * time.Second,
	}
}

func getEnvOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func parsePositiveInt(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
