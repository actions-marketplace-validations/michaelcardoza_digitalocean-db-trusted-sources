package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type actionKind string

const (
	actionKindAdd    actionKind = "add"
	actionKindRemove actionKind = "remove"
)

func (a actionKind) valid() bool {
	return a == actionKindAdd || a == actionKindRemove
}

type config struct {
	action actionKind
	dbID   string
	token  string
}

func loadConfig() (config, error) {
	cfg := config{
		action: actionKind(strings.TrimSpace(os.Getenv("INPUT_ACTION"))),
		dbID:   strings.TrimSpace(os.Getenv("INPUT_DATABASE_ID")),
		token:  strings.TrimSpace(os.Getenv("INPUT_DIGITALOCEAN_TOKEN")),
	}

	var errs []error
	if !cfg.action.valid() {
		errs = append(errs, fmt.Errorf("action must be %q or %q, got %q", actionKindAdd, actionKindRemove, cfg.action))
	}
	if cfg.dbID == "" {
		errs = append(errs, errors.New("database_id is required"))
	}
	if cfg.token == "" {
		errs = append(errs, errors.New("digitalocean_token is required"))
	}
	if len(errs) > 0 {
		return cfg, errors.Join(errs...)
	}
	return cfg, nil
}
