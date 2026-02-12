package network

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
)

func TestCheckValidSubnetMiddleware_EmptySubnet_AllowsRequest(t *testing.T) {
	var called int32
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&called, 1)
		w.WriteHeader(http.StatusOK)
	})

	mw := CheckValidSubnetMiddleware("")
	h := mw(next)

	req := httptest.NewRequest(http.MethodPost, "/updates/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d, body=%q", http.StatusOK, rr.Code, rr.Body.String())
	}
	if atomic.LoadInt32(&called) != 1 {
		t.Fatalf("expected next to be called once, got %d", called)
	}
}

func TestCheckValidSubnetMiddleware_MissingHeader_Returns400_AndDoesNotCallNext(t *testing.T) {
	var called int32
	logger.Init()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&called, 1)
		w.WriteHeader(http.StatusOK)
	})

	mw := CheckValidSubnetMiddleware("127.0.0.1/32")
	h := mw(next)

	req := httptest.NewRequest(http.MethodPost, "/updates/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d, body=%q", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
	if atomic.LoadInt32(&called) != 0 {
		t.Fatalf("expected next NOT to be called, got %d", called)
	}
}

func TestCheckValidSubnetMiddleware_InvalidHeaderIP_Returns400_AndDoesNotCallNext(t *testing.T) {
	var called int32
	logger.Init()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&called, 1)
		w.WriteHeader(http.StatusOK)
	})

	mw := CheckValidSubnetMiddleware("127.0.0.1/32")
	h := mw(next)

	req := httptest.NewRequest(http.MethodPost, "/updates/", nil)
	req.Header.Set("X-Real-IP", "not-an-ip")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d, body=%q", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
	if atomic.LoadInt32(&called) != 0 {
		t.Fatalf("expected next NOT to be called, got %d", called)
	}
}

func TestCheckValidSubnetMiddleware_IPNotInSubnet_Returns403_AndDoesNotCallNext(t *testing.T) {
	var called int32
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&called, 1)
		w.WriteHeader(http.StatusOK)
	})

	mw := CheckValidSubnetMiddleware("127.0.0.1/32")
	h := mw(next)

	req := httptest.NewRequest(http.MethodPost, "/updates/", nil)
	req.Header.Set("X-Real-IP", "127.0.0.2")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected %d, got %d, body=%q", http.StatusForbidden, rr.Code, rr.Body.String())
	}
	if atomic.LoadInt32(&called) != 0 {
		t.Fatalf("expected next NOT to be called, got %d", called)
	}
}

func TestCheckValidSubnetMiddleware_IPInSubnet_AllowsRequest(t *testing.T) {
	var called int32
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&called, 1)
		w.WriteHeader(http.StatusOK)
	})

	mw := CheckValidSubnetMiddleware("127.0.0.0/8")
	h := mw(next)

	req := httptest.NewRequest(http.MethodPost, "/updates/", nil)
	req.Header.Set("X-Real-IP", "127.0.0.1")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d, body=%q", http.StatusOK, rr.Code, rr.Body.String())
	}
	if atomic.LoadInt32(&called) != 1 {
		t.Fatalf("expected next to be called once, got %d", called)
	}
}
