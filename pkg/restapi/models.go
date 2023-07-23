package restapi

// Payload is the payload interface.
type Payload interface {
	Validate() error
}

// ErrorResponse is the error response.
type ErrorResponse struct {
	// StatusCode is the status code.
	StatusCode int `json:"status_code"`
	// ErrorMsg is the error message.
	ErrorMsg string `json:"error_msg"`
}
