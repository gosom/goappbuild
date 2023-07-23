package restapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var (
	// ErrEmptyBody is returned when the request body is empty.
	ErrEmptyBody = errors.New("empty body")
	// ErrInvalidJSON is returned when the request body is not a valid JSON.
	ErrInvalidJSON = errors.New("json syntax error")
	// ErrBadJSONData is returned when the request body contains invalid data.
	ErrBadJSONData = errors.New("invalid data. please check your request")
	// ErrNotFound is returned when the requested resource is not found.
	ErrNotFound = errors.New("not found")
	// ErrMethodNotAllowed is returned when the requested method is not allowed.
	ErrMethodNotAllowed = errors.New("method not allowed")
)

const (
	contentType     = "Content-Type"
	jsonContentType = "application/json"
)

// Controller is the base controller.
type Controller struct {
}

// Success writes a success response.
func (o Controller) Success(w http.ResponseWriter, _ *http.Request, code int, resp any) {
	w.Header().Set(contentType, jsonContentType)
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// TODO log error
		_ = err
	}
}

// Error writes an error response.
func (o Controller) Error(w http.ResponseWriter, _ *http.Request, code int, err error) {
	resp := ErrorResponse{
		StatusCode: code,
		ErrorMsg:   err.Error(),
	}

	w.Header().Set(contentType, jsonContentType)
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// TODO log error
		_ = err
	}
}

// StringParam returns a string parameter from the request url.
func (o Controller) StringURLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

// QueryParam returns a query parameter from the request url.
func (o Controller) QueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// HeaderKey returns a header value from the request.
func (o Controller) HeaderKey(r *http.Request, key string) string {
	return r.Header.Get(key)
}

// DecodeJSON decodes JSON and returns an error if reading body fails or validation fails.
func (o Controller) DecodeBody(r *http.Request, v Payload) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		var (
			jsonErr       *json.UnmarshalTypeError
			jsonSyntaxErr *json.SyntaxError
		)

		switch {
		case errors.Is(err, io.EOF):
			err = ErrEmptyBody
		case errors.As(err, &jsonErr):
			err = ErrBadJSONData
		case errors.As(err, &jsonSyntaxErr):
			err = ErrInvalidJSON
		}

		return err
	}

	return v.Validate()
}
