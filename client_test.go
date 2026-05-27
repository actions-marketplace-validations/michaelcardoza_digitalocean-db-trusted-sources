package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) (*doClient, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	return &doClient{
		httpClient: ts.Client(),
		baseURL:    ts.URL,
		token:      "test-token",
	}, ts
}

func TestGetFirewallRules(t *testing.T) {
	want := []rule{
		{Type: "ip_addr", Value: "1.1.1.1"},
		{Type: "tag", Value: "trusted"},
	}

	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/v2/databases/db-uuid/firewall" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("auth header = %q", got)
		}
		_ = json.NewEncoder(w).Encode(firewallRules{Rules: want})
	})

	got, err := c.GetFirewallRules(context.Background(), "db-uuid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("rule[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestPutFirewallRules(t *testing.T) {
	sent := []rule{{Type: "ip_addr", Value: "9.9.9.9"}}

	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/v2/databases/db-uuid/firewall" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("content-type = %q", got)
		}
		var body firewallRules
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if len(body.Rules) != 1 || body.Rules[0] != sent[0] {
			t.Errorf("body.Rules = %+v, want %+v", body.Rules, sent)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	if err := c.PutFirewallRules(context.Background(), "db-uuid", sent); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDoErrorIncludesStatusAndBody(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = io.WriteString(w, `{"message":"bad rule"}`)
	})

	_, err := c.GetFirewallRules(context.Background(), "db-uuid")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	msg := err.Error()
	if !strings.Contains(msg, "422") {
		t.Errorf("error missing status: %s", msg)
	}
	if !strings.Contains(msg, "bad rule") {
		t.Errorf("error missing body: %s", msg)
	}
}
