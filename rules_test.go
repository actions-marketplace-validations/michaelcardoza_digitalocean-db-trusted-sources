package main

import (
	"slices"
	"testing"
)

func TestRuleMatchesIP(t *testing.T) {
	cases := []struct {
		name string
		r    rule
		ip   string
		want bool
	}{
		{"match", rule{Type: "ip_addr", Value: "1.1.1.1"}, "1.1.1.1", true},
		{"wrong value", rule{Type: "ip_addr", Value: "1.1.1.1"}, "2.2.2.2", false},
		{"wrong type", rule{Type: "tag", Value: "1.1.1.1"}, "1.1.1.1", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.r.matchesIP(c.ip); got != c.want {
				t.Errorf("matchesIP(%q) = %v, want %v", c.ip, got, c.want)
			}
		})
	}
}

func TestAddIP(t *testing.T) {
	cases := []struct {
		name        string
		rules       []rule
		ip          string
		wantRules   []rule
		wantChanged bool
	}{
		{
			name:        "add to empty",
			rules:       nil,
			ip:          "1.1.1.1",
			wantRules:   []rule{{Type: "ip_addr", Value: "1.1.1.1"}},
			wantChanged: true,
		},
		{
			name:        "add new alongside existing",
			rules:       []rule{{Type: "ip_addr", Value: "2.2.2.2"}},
			ip:          "1.1.1.1",
			wantRules:   []rule{{Type: "ip_addr", Value: "2.2.2.2"}, {Type: "ip_addr", Value: "1.1.1.1"}},
			wantChanged: true,
		},
		{
			name:        "already present, no-op",
			rules:       []rule{{Type: "ip_addr", Value: "1.1.1.1"}},
			ip:          "1.1.1.1",
			wantRules:   []rule{{Type: "ip_addr", Value: "1.1.1.1"}},
			wantChanged: false,
		},
		{
			name:        "matching value but different type, still adds",
			rules:       []rule{{Type: "tag", Value: "1.1.1.1"}},
			ip:          "1.1.1.1",
			wantRules:   []rule{{Type: "tag", Value: "1.1.1.1"}, {Type: "ip_addr", Value: "1.1.1.1"}},
			wantChanged: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, changed := addIP(c.rules, c.ip)
			if !slices.Equal(got, c.wantRules) {
				t.Errorf("rules = %+v, want %+v", got, c.wantRules)
			}
			if changed != c.wantChanged {
				t.Errorf("changed = %v, want %v", changed, c.wantChanged)
			}
		})
	}
}

func TestRemoveIP(t *testing.T) {
	cases := []struct {
		name        string
		rules       []rule
		ip          string
		wantRules   []rule
		wantChanged bool
	}{
		{
			name:        "remove existing",
			rules:       []rule{{Type: "ip_addr", Value: "1.1.1.1"}},
			ip:          "1.1.1.1",
			wantRules:   []rule{},
			wantChanged: true,
		},
		{
			name:        "remove from middle",
			rules:       []rule{{Type: "ip_addr", Value: "a"}, {Type: "ip_addr", Value: "b"}, {Type: "ip_addr", Value: "c"}},
			ip:          "b",
			wantRules:   []rule{{Type: "ip_addr", Value: "a"}, {Type: "ip_addr", Value: "c"}},
			wantChanged: true,
		},
		{
			name:        "not present, no-op",
			rules:       []rule{{Type: "ip_addr", Value: "1.1.1.1"}},
			ip:          "9.9.9.9",
			wantRules:   []rule{{Type: "ip_addr", Value: "1.1.1.1"}},
			wantChanged: false,
		},
		{
			name:        "empty input, no-op",
			rules:       nil,
			ip:          "1.1.1.1",
			wantRules:   nil,
			wantChanged: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, changed := removeIP(c.rules, c.ip)
			if !slices.Equal(got, c.wantRules) {
				t.Errorf("rules = %+v, want %+v", got, c.wantRules)
			}
			if changed != c.wantChanged {
				t.Errorf("changed = %v, want %v", changed, c.wantChanged)
			}
		})
	}
}

func TestAddIPDoesNotMutateInput(t *testing.T) {
	original := []rule{{Type: "ip_addr", Value: "1.1.1.1"}}
	snapshot := slices.Clone(original)

	_, _ = addIP(original, "2.2.2.2")

	if !slices.Equal(original, snapshot) {
		t.Errorf("input mutated: got %+v, want %+v", original, snapshot)
	}
}

func TestRemoveIPDoesNotMutateInput(t *testing.T) {
	original := []rule{{Type: "ip_addr", Value: "1.1.1.1"}, {Type: "ip_addr", Value: "2.2.2.2"}}
	snapshot := slices.Clone(original)

	_, _ = removeIP(original, "1.1.1.1")

	if !slices.Equal(original, snapshot) {
		t.Errorf("input mutated: got %+v, want %+v", original, snapshot)
	}
}
