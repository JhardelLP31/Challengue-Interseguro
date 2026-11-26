package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/JhardelLP31/Challengue-Interseguro/go-matrix-api/internal/model"
)

// StatsClient define el contrato para comunicarse por HTTP con la API de estadísticas en Node.js.
type StatsClient interface {
	CalculateStatistics(ctx context.Context, authHeader string, matrices map[string]model.Matrix) (*model.StatsResponse, error)
}

type httpStatsClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewHTTPStatsClient crea un cliente HTTP con timeout configurable hacia node-stats-api.
func NewHTTPStatsClient(baseURL string, timeout time.Duration) StatsClient {
	return &httpStatsClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// CalculateStatistics envía las matrices procesadas a POST /api/v1/stats propagando el token JWT.
func (c *httpStatsClient) CalculateStatistics(
	ctx context.Context,
	authHeader string,
	matrices map[string]model.Matrix,
) (*model.StatsResponse, error) {
	payload := model.StatsRequest{Matrices: matrices}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializando payload para API de estadísticas: %w", err)
	}

	endpoint := fmt.Sprintf("%s/api/v1/stats", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("error creando petición HTTP hacia API de estadísticas: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error conectando con la API de estadísticas (%s): %w", endpoint, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta de la API de estadísticas: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("la API de estadísticas respondió HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var stats model.StatsResponse
	if err := json.Unmarshal(respBody, &stats); err != nil {
		return nil, fmt.Errorf("error deserializando respuesta de estadísticas: %w", err)
	}

	return &stats, nil
}
