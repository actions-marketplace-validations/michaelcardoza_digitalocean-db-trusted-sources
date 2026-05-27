package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

// version is set at build time via -ldflags "-X main.version=..."
var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println(version)
		return
	}

	cfg, err := loadConfig()
	if err != nil {
		fail(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := run(ctx, cfg); err != nil {
		fail(err)
	}
}

func run(ctx context.Context, cfg config) error {
	ip, err := lookupRunnerIP(ctx)
	if err != nil {
		return fmt.Errorf("lookup runner IP: %w", err)
	}
	fmt.Println("runner IP:", ip)

	doClient := newDOClient(cfg.token)
	rules, err := doClient.GetFirewallRules(ctx, cfg.dbID)
	if err != nil {
		return fmt.Errorf("get firewall rules: %w", err)
	}

	var changed bool
	switch cfg.action {
	case actionKindAdd:
		rules, changed = addIP(rules, ip)
	case actionKindRemove:
		rules, changed = removeIP(rules, ip)
	}
	if !changed {
		fmt.Println("no changes needed: rule already up to date")
		return nil
	}

	if err := doClient.PutFirewallRules(ctx, cfg.dbID, rules); err != nil {
		return fmt.Errorf("update firewall rules: %w", err)
	}
	return nil
}

func fail(err error) {
	fmt.Printf("::error::%s\n", err)
	os.Exit(1)
}
