package payloads

import (
	"github.com/go-chi/render"
	"net/http"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	Message string `json:"message"`          // user-level message
	ErrorCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		Message:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		Message:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

func ErrServer(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		Message:     "Server Error",
		ErrorText:      err.Error(),
	}
}

var ErrUnauthorized = &ErrResponse{HTTPStatusCode: 404, Message: "Unauthorized"}
var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, Message: "Resource not found."}