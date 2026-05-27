package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLookupRunnerIP(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = io.WriteString(w, "203.0.113.42\n")
		}))
		t.Cleanup(ts.Close)
		ipifyURL = ts.URL

		got, err := lookupRunnerIP(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "203.0.113.42" {
			t.Errorf("ip = %q, want %q", got, "203.0.113.42")
		}
	})

	t.Run("invalid IP body", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = io.WriteString(w, "not-an-ip")
		}))
		t.Cleanup(ts.Close)
		ipifyURL = ts.URL

		_, err := lookupRunnerIP(context.Background())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "invalid IP") {
			t.Errorf("error missing 'invalid IP': %s", err)
		}
	})

	t.Run("non-200 status", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		t.Cleanup(ts.Close)
		ipifyURL = ts.URL

		_, err := lookupRunnerIP(context.Background())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "503") {
			t.Errorf("error missing status: %s", err)
		}
	})

	t.Run("IPv6", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = io.WriteString(w, "2001:db8::1")
		}))
		t.Cleanup(ts.Close)
		ipifyURL = ts.URL

		got, err := lookupRunnerIP(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "2001:db8::1" {
			t.Errorf("ip = %q", got)
		}
	})
}
