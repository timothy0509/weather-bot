package hko

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBodySizeLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(make([]byte, maxBodySize+1))
	}))
	defer server.Close()

	c := NewClient(5 * time.Second)
	var v map[string]interface{}
	err := c.Get(server.URL, &v)
	if err == nil {
		t.Error("expected error for oversized response")
	}
}

func TestRetryOn503(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n <= 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("unavailable"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	c := NewClient(5 * time.Second)
	var v map[string]string
	err := c.Get(server.URL, &v)
	if err != nil {
		t.Fatalf("expected success after retries, got: %v", err)
	}
	if v["status"] != "ok" {
		t.Errorf("expected status=ok, got %v", v)
	}
	if got := atomic.LoadInt32(&attempts); got != 3 {
		t.Errorf("expected 3 attempts, got %d", got)
	}
}

func TestRetryExhausted(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("unavailable"))
	}))
	defer server.Close()

	c := NewClient(5 * time.Second)
	var v map[string]string
	err := c.Get(server.URL, &v)
	if err == nil {
		t.Error("expected error after exhausting retries")
	}
}

func TestNoRetryOn4xx(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
	}))
	defer server.Close()

	c := NewClient(5 * time.Second)
	var v map[string]string
	err := c.Get(server.URL, &v)
	if err == nil {
		t.Error("expected error for 400")
	}
	if got := atomic.LoadInt32(&attempts); got != 1 {
		t.Errorf("expected 1 attempt (no retry on 4xx), got %d", got)
	}
}

func TestCacheHit(t *testing.T) {
	var requests int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"value": "fresh"})
	}))
	defer server.Close()

	c := NewClient(5 * time.Second)

	var v1 map[string]string
	if err := c.GetWithTTL(server.URL, &v1, 5*time.Minute); err != nil {
		t.Fatalf("first request: %v", err)
	}

	var v2 map[string]string
	if err := c.GetWithTTL(server.URL, &v2, 5*time.Minute); err != nil {
		t.Fatalf("second request: %v", err)
	}

	if v1["value"] != "fresh" || v2["value"] != "fresh" {
		t.Errorf("expected cached value, got v1=%v v2=%v", v1, v2)
	}
	if got := atomic.LoadInt32(&requests); got != 1 {
		t.Errorf("expected 1 HTTP request (cache hit on 2nd), got %d", got)
	}
}

func TestCacheExpiry(t *testing.T) {
	var requests int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&requests, 1)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"request": int(n)})
	}))
	defer server.Close()

	c := NewClient(5 * time.Second)

	var v1 map[string]int
	if err := c.GetWithTTL(server.URL, &v1, 1*time.Millisecond); err != nil {
		t.Fatalf("first request: %v", err)
	}

	time.Sleep(5 * time.Millisecond)

	var v2 map[string]int
	if err := c.GetWithTTL(server.URL, &v2, 1*time.Millisecond); err != nil {
		t.Fatalf("second request: %v", err)
	}

	if v1["request"] == v2["request"] {
		t.Error("expected different values after cache expiry")
	}
	if got := atomic.LoadInt32(&requests); got != 2 {
		t.Errorf("expected 2 HTTP requests after cache expiry, got %d", got)
	}
}

func TestSingleflightDedup(t *testing.T) {
	var requests int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		time.Sleep(50 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"value": "deduped"})
	}))
	defer server.Close()

	c := NewClient(5 * time.Second)

	var wg sync.WaitGroup
	results := make([]map[string]string, 5)
	errors := make([]error, 5)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			errors[idx] = c.Get(server.URL, &results[idx])
		}(i)
	}
	wg.Wait()

	for i, err := range errors {
		if err != nil {
			t.Errorf("goroutine %d failed: %v", i, err)
		}
	}

	if got := atomic.LoadInt32(&requests); got != 1 {
		t.Errorf("expected 1 HTTP request (singleflight dedup), got %d", got)
	}
}

func TestGetWithTTLZeroTTLNoCache(t *testing.T) {
	var requests int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"n": int(requests)})
	}))
	defer server.Close()

	c := NewClient(5 * time.Second)

	var v1, v2 map[string]int
	c.GetWithTTL(server.URL, &v1, 0)
	c.GetWithTTL(server.URL, &v2, 0)

	if got := atomic.LoadInt32(&requests); got != 2 {
		t.Errorf("expected 2 requests with TTL=0 (no caching), got %d", got)
	}
}
