package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

const ipifyTimeout = 10 * time.Second

var ipifyURL = "https://api.ipify.org"

func lookupRunnerIP(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ipifyURL, nil)
	if err != nil {
		return "", fmt.Errorf("build ipify request: %w", err)
	}

	httpClient := &http.Client{Timeout: ipifyTimeout}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ipify: status %d", resp.StatusCode)
	}

	rawIP, err := io.ReadAll(io.LimitReader(resp.Body, 64))
	if err != nil {
		return "", err
	}

	ip := strings.TrimSpace(string(rawIP))
	if net.ParseIP(ip) == nil {
		return "", fmt.Errorf("ipify: invalid IP %q", ip)
	}

	return ip, nil
}
