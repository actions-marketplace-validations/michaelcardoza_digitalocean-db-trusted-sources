package main

import "slices"

const ruleTypeIP = "ip_addr"

type rule struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type firewallRules struct {
	Rules []rule `json:"rules"`
}

func (r rule) matchesIP(ip string) bool {
	return r.Type == ruleTypeIP && r.Value == ip
}

func addIP(rules []rule, ip string) ([]rule, bool) {
	exists := slices.ContainsFunc(rules, func(r rule) bool {
		return r.matchesIP(ip)
	})
	if exists {
		return rules, false
	}

	updated := slices.Clone(rules)
	updated = append(updated, rule{Type: ruleTypeIP, Value: ip})
	return updated, true
}

func removeIP(rules []rule, ip string) ([]rule, bool) {
	idx := slices.IndexFunc(rules, func(r rule) bool {
		return r.matchesIP(ip)
	})
	if idx == -1 {
		return rules, false
	}

	updated := slices.Clone(rules)
	updated = slices.Delete(updated, idx, idx+1)
	return updated, true
}
