package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"go-backend/errors"
)

func ValidationMiddleware(maxSize int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > maxSize {
				SendResponse(w, r, Error(errors.NewRequestTooLargeError(maxSize, r.ContentLength), r))
				return
			}
			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
				if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
					SendResponse(w, r, Error(errors.NewValidationError("Content-Type must be application/json", nil), r))
					return
				}
				body, err := io.ReadAll(io.LimitReader(r.Body, maxSize))
				if err != nil {
					SendResponse(w, r, Error(errors.NewInvalidJSONError(err), r))
					return
				}
				if len(body) == 0 && (r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch) {
					SendResponse(w, r, Error(errors.NewValidationError("Request body cannot be empty", nil), r))
					return
				}
				if !json.Valid(body) {
					SendResponse(w, r, Error(errors.NewInvalidJSONError(nil), r))
					return
				}
				r.Body = io.NopCloser(strings.NewReader(string(body)))
			}
			next.ServeHTTP(w, r)
		})
	}
}
