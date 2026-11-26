package model

// Matrix representa una matriz numérica bidimensional.
type Matrix [][]float64

// MatrixRequest representa el payload de entrada para operaciones sobre matrices.
type MatrixRequest struct {
	Matrix    Matrix `json:"matrix"`
	Direction string `json:"direction,omitempty"`
}

// MatrixInputMetadata describe las dimensiones y datos de la matriz recibida.
type MatrixInputMetadata struct {
	Rows   int    `json:"rows"`
	Cols   int    `json:"cols"`
	Matrix Matrix `json:"matrix"`
}

// QRResult contiene las matrices Q (ortogonal) y R (triangular superior) de la factorización QR.
type QRResult struct {
	Q Matrix `json:"q"`
	R Matrix `json:"r"`
}

// StatsRequest es el payload enviado por HTTP hacia la API en Node.js.
type StatsRequest struct {
	Matrices map[string]Matrix `json:"matrices"`
}

// StatsResponse refleja el resultado devuelto por la API de estadísticas en Node.js.
type StatsResponse struct {
	Max              float64         `json:"max"`
	Min              float64         `json:"min"`
	Average          float64         `json:"average"`
	Sum              float64         `json:"sum"`
	TotalElements    int             `json:"totalElements"`
	IsAnyDiagonal    bool            `json:"isAnyDiagonal"`
	DiagonalByMatrix map[string]bool `json:"diagonalByMatrix"`
}

// QRResponse es la respuesta final del endpoint de factorización QR.
type QRResponse struct {
	Operation  string              `json:"operation"`
	Input      MatrixInputMetadata `json:"input"`
	QR         QRResult            `json:"qr"`
	Statistics *StatsResponse      `json:"statistics"`
}

// RotateResponse es la respuesta final del endpoint de rotación de matriz.
type RotateResponse struct {
	Operation  string              `json:"operation"`
	Direction  string              `json:"direction"`
	Input      MatrixInputMetadata `json:"input"`
	Rotated    Matrix              `json:"rotated"`
	Statistics *StatsResponse      `json:"statistics"`
}

// TokenRequest representa las credenciales enviadas para solicitar un token JWT.
type TokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TokenResponse representa la respuesta con el JWT generado.
type TokenResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"tokenType"`
	ExpiresIn int64  `json:"expiresIn"`
}

// ErrorResponse estandariza el formato de errores HTTP de la API.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
