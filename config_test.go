package main

import (
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("valid add", func(t *testing.T) {
		t.Setenv("INPUT_ACTION", "add")
		t.Setenv("INPUT_DATABASE_ID", "db-uuid")
		t.Setenv("INPUT_DIGITALOCEAN_TOKEN", "tok")

		cfg, err := loadConfig()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.action != actionKindAdd {
			t.Errorf("action = %q, want %q", cfg.action, actionKindAdd)
		}
		if cfg.dbID != "db-uuid" {
			t.Errorf("dbID = %q, want %q", cfg.dbID, "db-uuid")
		}
		if cfg.token != "tok" {
			t.Errorf("token = %q, want %q", cfg.token, "tok")
		}
	})

	t.Run("trims whitespace", func(t *testing.T) {
		t.Setenv("INPUT_ACTION", "  remove  ")
		t.Setenv("INPUT_DATABASE_ID", " db ")
		t.Setenv("INPUT_DIGITALOCEAN_TOKEN", " tok ")

		cfg, err := loadConfig()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.action != actionKindRemove {
			t.Errorf("action = %q, want %q", cfg.action, actionKindRemove)
		}
		if cfg.dbID != "db" || cfg.token != "tok" {
			t.Errorf("not trimmed: %+v", cfg)
		}
	})

	t.Run("invalid action", func(t *testing.T) {
		t.Setenv("INPUT_ACTION", "delete")
		t.Setenv("INPUT_DATABASE_ID", "db")
		t.Setenv("INPUT_DIGITALOCEAN_TOKEN", "tok")

		_, err := loadConfig()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "action must be") {
			t.Errorf("error %q missing action message", err)
		}
	})

	t.Run("missing all", func(t *testing.T) {
		t.Setenv("INPUT_ACTION", "")
		t.Setenv("INPUT_DATABASE_ID", "")
		t.Setenv("INPUT_DIGITALOCEAN_TOKEN", "")

		_, err := loadConfig()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		msg := err.Error()
		for _, want := range []string{"action must be", "database_id is required", "digitalocean_token is required"} {
			if !strings.Contains(msg, want) {
				t.Errorf("error missing %q\n got: %s", want, msg)
			}
		}
	})

	t.Run("missing dbID only", func(t *testing.T) {
		t.Setenv("INPUT_ACTION", "add")
		t.Setenv("INPUT_DATABASE_ID", "")
		t.Setenv("INPUT_DIGITALOCEAN_TOKEN", "tok")

		_, err := loadConfig()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "database_id is required") {
			t.Errorf("error %q missing dbID message", err)
		}
	})
}
