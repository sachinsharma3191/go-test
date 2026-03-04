package middleware

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// errReader returns an error on Read to trigger io.ReadAll error path.
type errReader struct{}

func (errReader) Read(_ []byte) (int, error) {
	return 0, errors.New("read failed")
}

func makeValidationHandler(maxSize int64, nextCalled *bool) http.Handler {
	vmw := ValidationMiddleware(maxSize)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*nextCalled = true
		w.WriteHeader(http.StatusNoContent)
	})
	return vmw(next)
}

func TestValidationMiddleware_AllowsValidJSON(t *testing.T) {
	called := false
	h := makeValidationHandler(1024, &called)

	body := `{"name":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if !called {
		t.Fatalf("expected next handler to be called")
	}
	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204 from next handler, got %d", rr.Code)
	}
}

func TestValidationMiddleware_RejectsTooLargeContentLength(t *testing.T) {
	called := false
	h := makeValidationHandler(10, &called)

	body := `{"name":"this is longer than 10 bytes"}`
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.ContentLength = int64(len(body))
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if called {
		t.Fatalf("did not expect next handler to be called")
	}
	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413, got %d", rr.Code)
	}
}

func TestValidationMiddleware_ReadBodyError(t *testing.T) {
	called := false
	h := makeValidationHandler(1024, &called)

	req := httptest.NewRequest(http.MethodPost, "/test", io.NopCloser(errReader{}))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if called {
		t.Fatalf("did not expect next handler to be called when body read fails")
	}
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 on read error, got %d", rr.Code)
	}
}

func TestValidationMiddleware_RejectsNonJSONContentType(t *testing.T) {
	called := false
	h := makeValidationHandler(1024, &called)

	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("body"))
	req.Header.Set("Content-Type", "text/plain")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if called {
		t.Fatalf("did not expect next handler to be called")
	}
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestValidationMiddleware_RejectsEmptyBody(t *testing.T) {
	called := false
	h := makeValidationHandler(1024, &called)

	req := httptest.NewRequest(http.MethodPost, "/test", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if called {
		t.Fatalf("did not expect next handler to be called")
	}
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestValidationMiddleware_RejectsInvalidJSON(t *testing.T) {
	called := false
	h := makeValidationHandler(1024, &called)

	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if called {
		t.Fatalf("did not expect next handler to be called")
	}
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

// Ensure that non-body methods (GET) bypass JSON validation but still go through next.
func TestValidationMiddleware_AllowsGETWithoutBody(t *testing.T) {
	called := false
	h := makeValidationHandler(10, &called)

	req := httptest.NewRequest(http.MethodGet, "/test", io.NopCloser(strings.NewReader("whatever")))
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if !called {
		t.Fatalf("expected next handler to be called for GET")
	}
}

func TestValidationMiddleware_RejectsEmptyBody_PUT(t *testing.T) {
	called := false
	h := makeValidationHandler(1024, &called)
	req := httptest.NewRequest(http.MethodPut, "/test", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if called {
		t.Fatal("should not call next for empty PUT body")
	}
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestValidationMiddleware_RejectsEmptyBody_PATCH(t *testing.T) {
	called := false
	h := makeValidationHandler(1024, &called)
	req := httptest.NewRequest(http.MethodPatch, "/test", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if called {
		t.Fatal("should not call next for empty PATCH body")
	}
}
