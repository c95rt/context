package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ResponseWriter wrap common http responses
type ResponseWriter struct {
	writer http.ResponseWriter
}

type generalResponse struct {
	Errors  []*errorResponse `json:"errors"`
	Success bool             `json:"success"`
	Data    interface{}      `json:"data"`
}

type errorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Scope   string      `json:"scope"`
	Type    int         `json:"type"`
	Data    interface{} `json:"data"`
}

// ErrOption custom type to handle functional options
type ErrOption func(*errorResponse)

// WithErrorType modifies error type from any error response struct
func WithErrorType(errType int) ErrOption {
	return func(err *errorResponse) {
		err.Type = errType
	}
}

// WithErrorScope modifies error scope from any error response struct
func WithErrorScope(scope string) ErrOption {
	return func(err *errorResponse) {
		err.Scope = scope
	}
}

func (r *ResponseWriter) writeJSONResponse(code int, errors []*errorResponse, data interface{}) {
	r.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := &generalResponse{Errors: errors, Success: errors == nil, Data: data}
	b, err := json.Marshal(response)
	if err != nil {
		r.writer.WriteHeader(http.StatusInternalServerError)
		r.writer.Write([]byte(fmt.Sprintf("unexpected error: %v", err)))
	}
	r.writer.WriteHeader(code)
	if code, err := r.writer.Write(b); err != nil {
		fmt.Sprintf("could not response - code: %d", code)
	}
}

// JSON converts the given interface in a JSON response
// adding the expected headers.
func (r *ResponseWriter) JSON(code int, data interface{}) {
	r.writeJSONResponse(code, nil, data)
}

func (r *ResponseWriter) writePlainJSONResponse(statusCode int, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		r.writer.WriteHeader(http.StatusInternalServerError)
		r.writer.Write([]byte(fmt.Sprintf("unexpected error: %v", err)))
	}

	r.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	r.writer.WriteHeader(statusCode)

	if code, err := r.writer.Write(b); err != nil {
		fmt.Sprintf("could not response - code: %d", code)
	}
}

func (r *ResponseWriter) WriteJSON(statusCode int, data interface{}) {
	r.writePlainJSONResponse(statusCode, data)
}

// Stringf similar to `fmt.Printf` creates a plain response
// based on the given format.
func (r *ResponseWriter) Stringf(code int, format string, args ...interface{}) {
	r.writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	r.writer.WriteHeader(code)
	if code, err := r.writer.Write([]byte(fmt.Sprintf(format, args...))); err != nil {
		fmt.Sprintf("could not response - code: %d", code)
	}
}

// Errorf similar to `fmt.Printf` creates a JSON response
// based on the given format.
func (r *ResponseWriter) Errorf(code int, format string, args ...interface{}) {
	errors := []*errorResponse{
		{Code: code, Message: fmt.Sprintf(format, args...)},
	}
	r.writeJSONResponse(code, errors, nil)
}

// ErrorJ creates a JSON response with error
func (r *ResponseWriter) ErrorJ(code int, format string, data interface{}) {
	errors := []*errorResponse{
		{Code: code, Message: format, Data: data},
	}
	r.writeJSONResponse(code, errors, nil)
}

// String writes a plain response
func (r *ResponseWriter) String(code int, msg string) {
	r.writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	r.writer.WriteHeader(code)
	if code, err := r.writer.Write([]byte(msg)); err != nil {
		fmt.Sprintf("could not response - code: %d", code)
	}
}

// Errorf writes an standard error using the `errorHandler` struct.
func (r *ResponseWriter) Error(code int, msg string, opts ...ErrOption) {
	err := &errorResponse{Code: code, Message: msg}
	for _, With := range opts {
		With(err)
	}
	r.writeJSONResponse(code, []*errorResponse{err}, nil)
}
