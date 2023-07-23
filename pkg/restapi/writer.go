package restapi

import (
	"bytes"
	"io"
	"net/http"
)

// ResponseWriter is a wrapper around http.ResponseWriter that we
// can use to capture the response body and status code and headers
type ResponseWriter struct {
	orig       http.ResponseWriter
	statusCode int
	buffer     *bytes.Buffer
	writer     io.Writer
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	buffer := new(bytes.Buffer)

	return &ResponseWriter{
		orig:   w,
		buffer: buffer,
		writer: io.MultiWriter(w, buffer),
	}
}

// WriteHeader writes the header to the response and stores the status code
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.orig.WriteHeader(statusCode)
}

// Write writes the body to the ResponseWriter and the buffer
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	return rw.writer.Write(b)
}

// StatusCode returns the status code
func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}

// Body returns the response body that was written to the buffer
func (rw *ResponseWriter) Body() []byte {
	return rw.buffer.Bytes()
}

func (rw *ResponseWriter) Header() http.Header {
	return rw.orig.Header()
}
