package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const doBaseURL = "https://api.digitalocean.com"

type doClient struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

func newDOClient(token string) *doClient {
	return &doClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    doBaseURL,
		token:      token,
	}
}

func (c *doClient) GetFirewallRules(ctx context.Context, dbID string) ([]rule, error) {
	path := "/v2/databases/" + dbID + "/firewall"
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var body firewallRules
	if err := c.do(req, &body); err != nil {
		return nil, err
	}

	return body.Rules, nil
}

func (c *doClient) PutFirewallRules(ctx context.Context, dbID string, rules []rule) error {
	path := "/v2/databases/" + dbID + "/firewall"
	body := firewallRules{Rules: rules}
	req, err := c.newRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

func (c *doClient) newRequest(ctx context.Context, method, path string, payload any) (*http.Request, error) {
	var reqBody io.Reader
	if payload != nil {
		encoded, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("marshal payload: %w", err)
		}
		reqBody = bytes.NewReader(encoded)
	}

	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *doClient) do(req *http.Request, out any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		errBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("%s %s: status %d: %s", req.Method, req.URL.Path, resp.StatusCode, errBody)
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
