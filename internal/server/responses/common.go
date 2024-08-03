package responses

type StatusResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []any  `json:"data,omitempty"`
}
