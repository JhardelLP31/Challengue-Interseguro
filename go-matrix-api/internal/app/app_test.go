package app

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/JhardelLP31/Challengue-Interseguro/go-matrix-api/internal/client"
	"github.com/JhardelLP31/Challengue-Interseguro/go-matrix-api/internal/config"
	"github.com/JhardelLP31/Challengue-Interseguro/go-matrix-api/internal/model"
)

func testConfig(statsURL string) config.Config {
	return config.Config{
		Port:         "8080",
		StatsAPIURL:  statsURL,
		JWTSecret:    "test-secret-key",
		JWTExpiresIn: time.Hour,
		AuthUsername: "interseguro",
		AuthPassword: "challenge2024",
		HTTPTimeout:  5 * time.Second,
	}
}

func obtainTestJWT(t *testing.T, fiberApp interface {
	Test(*http.Request, ...int) (*http.Response, error)
}) string {
	t.Helper()

	body := `{"username":"interseguro","password":"challenge2024"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/token", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := fiberApp.Test(req)
	if err != nil {
		t.Fatalf("error en POST /api/v1/auth/token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("esperado HTTP 200 al pedir token, obtenido %d", resp.StatusCode)
	}

	var tokenResp model.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		t.Fatalf("error decodificando token: %v", err)
	}
	return tokenResp.Token
}

func TestHealthEndpoint(t *testing.T) {
	cfg := testConfig("http://localhost:3000")
	fiberApp := New(cfg, client.NewHTTPStatsClient(cfg.StatsAPIURL, cfg.HTTPTimeout))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp, err := fiberApp.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando /health: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("esperado 200 OK, obtenido %d", resp.StatusCode)
	}
}

func TestQREndpoint_RequiresJWT(t *testing.T) {
	cfg := testConfig("http://localhost:3000")
	fiberApp := New(cfg, client.NewHTTPStatsClient(cfg.StatsAPIURL, cfg.HTTPTimeout))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/matrices/qr", bytes.NewBufferString(`{"matrix":[[1,2],[3,4]]}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := fiberApp.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando petición: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("esperado 401 Unauthorized sin JWT, obtenido %d", resp.StatusCode)
	}
}

func TestQREndpoint_CallsStatsAPIWithJWT(t *testing.T) {
	var receivedAuthHeader string
	var receivedMatrices map[string]model.Matrix

	mockStatsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuthHeader = r.Header.Get("Authorization")
		rawBody, _ := io.ReadAll(r.Body)

		var payload model.StatsRequest
		if err := json.Unmarshal(rawBody, &payload); err != nil {
			t.Fatalf("payload inválido enviado a stats API: %v", err)
		}
		receivedMatrices = payload.Matrices

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(model.StatsResponse{
			Max:           5,
			Min:           0,
			Average:       1.125,
			Sum:           9,
			TotalElements: 8,
			IsAnyDiagonal: true,
			DiagonalByMatrix: map[string]bool{
				"Q": true,
				"R": true,
			},
		})
	}))
	defer mockStatsServer.Close()

	cfg := testConfig(mockStatsServer.URL)
	fiberApp := New(cfg, client.NewHTTPStatsClient(cfg.StatsAPIURL, cfg.HTTPTimeout))
	token := obtainTestJWT(t, fiberApp)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/matrices/qr", bytes.NewBufferString(`{"matrix":[[3,0],[0,5]]}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := fiberApp.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando POST /api/v1/matrices/qr: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("esperado 200 OK, obtenido %d: %s", resp.StatusCode, string(body))
	}

	if receivedAuthHeader != "Bearer "+token {
		t.Fatalf("esperado propagación del header Authorization Bearer, obtenido %q", receivedAuthHeader)
	}
	if _, ok := receivedMatrices["Q"]; !ok {
		t.Fatal("esperado que se envíe matriz Q a la API de estadísticas")
	}
	if _, ok := receivedMatrices["R"]; !ok {
		t.Fatal("esperado que se envíe matriz R a la API de estadísticas")
	}

	var qrResp model.QRResponse
	if err := json.NewDecoder(resp.Body).Decode(&qrResp); err != nil {
		t.Fatalf("error decodificando QRResponse: %v", err)
	}
	if qrResp.Operation != "qr_factorization" || qrResp.Statistics == nil || !qrResp.Statistics.IsAnyDiagonal {
		t.Fatalf("respuesta inesperada: %+v", qrResp)
	}
}

func TestRotateEndpoint_Success(t *testing.T) {
	mockStatsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(model.StatsResponse{
			Max:           6,
			Min:           1,
			Average:       3.5,
			Sum:           21,
			TotalElements: 6,
			IsAnyDiagonal: false,
			DiagonalByMatrix: map[string]bool{
				"rotated": false,
			},
		})
	}))
	defer mockStatsServer.Close()

	cfg := testConfig(mockStatsServer.URL)
	fiberApp := New(cfg, client.NewHTTPStatsClient(cfg.StatsAPIURL, cfg.HTTPTimeout))
	token := obtainTestJWT(t, fiberApp)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/matrices/rotate",
		bytes.NewBufferString(`{"matrix":[[1,2,3],[4,5,6]],"direction":"clockwise"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := fiberApp.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando POST /api/v1/matrices/rotate: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("esperado 200 OK, obtenido %d", resp.StatusCode)
	}

	var rotateResp model.RotateResponse
	if err := json.NewDecoder(resp.Body).Decode(&rotateResp); err != nil {
		t.Fatalf("error decodificando RotateResponse: %v", err)
	}
	if len(rotateResp.Rotated) != 3 || len(rotateResp.Rotated[0]) != 2 || rotateResp.Rotated[0][0] != 4 {
		t.Fatalf("matriz rotada inesperada: %+v", rotateResp.Rotated)
	}
}
