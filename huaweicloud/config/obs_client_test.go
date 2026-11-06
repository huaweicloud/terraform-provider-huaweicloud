package config

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/chnsz/golangsdk"
)

// TestObjectStorageClient_retrySignsAgain checks that a request which is retried after a connection error is
// signed again, with a new Date header. The first request stalls for longer than the resolution of the Date header
// (one second) and then fails, so a retry that resends the first request would carry the same Date header.
func TestObjectStorageClient_retrySignsAgain(t *testing.T) {
	var (
		mu    sync.Mutex
		dates []string
		auths []string
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		dates = append(dates, r.Header.Get("Date"))
		auths = append(auths, r.Header.Get("Authorization"))
		attempt := len(dates)
		mu.Unlock()

		if attempt == 1 {
			time.Sleep(1100 * time.Millisecond)
			conn, _, err := w.(http.Hijacker).Hijack()
			if err != nil {
				t.Errorf("failed to hijack the connection: %s", err)
				return
			}
			_ = conn.Close()
			return
		}

		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><ListAllMyBucketsResult></ListAllMyBucketsResult>`))
	}))
	defer server.Close()

	// "localhost" has no dots, so getObsEndpoint keeps the endpoint as it is.
	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("failed to parse the server URL: %s", err)
	}
	endpoint := "http://localhost:" + serverURL.Port()

	clientFuncs := map[string]func(*Config, string) error{
		"ObjectStorageClient": func(c *Config, region string) error {
			client, err := c.ObjectStorageClient(region)
			if err != nil {
				return err
			}
			_, err = client.ListBuckets(nil)
			return err
		},
		"ObjectStorageClientWithSignature": func(c *Config, region string) error {
			client, err := c.ObjectStorageClientWithSignature(region)
			if err != nil {
				return err
			}
			_, err = client.ListBuckets(nil)
			return err
		},
	}

	for name, call := range clientFuncs {
		t.Run(name, func(t *testing.T) {
			mu.Lock()
			dates = nil
			auths = nil
			mu.Unlock()

			domainClient := &golangsdk.ProviderClient{}
			domainClient.HTTPClient = http.Client{
				Transport: &LogRoundTripper{
					Rt:         &http.Transport{},
					MaxRetries: 1,
				},
			}
			c := &Config{
				AccessKey:    "access-key",
				SecretKey:    "secret-key",
				MaxRetries:   1,
				Endpoints:    map[string]string{"obs": endpoint},
				DomainClient: domainClient,
			}

			if err := call(c, "region-1"); err != nil {
				t.Fatalf("the request failed: %s", err)
			}

			mu.Lock()
			defer mu.Unlock()
			if len(dates) != 2 {
				t.Fatalf("expected 2 attempts, got %d: %v", len(dates), dates)
			}
			if dates[0] == "" || dates[0] == dates[1] {
				t.Errorf("the retry was not signed again, the Date headers of both attempts are %q and %q",
					dates[0], dates[1])
			}
			if auths[0] == "" || auths[1] == "" || auths[0] == auths[1] {
				t.Errorf("expected both attempts to carry their own signature")
			}

			lrt, ok := c.DomainClient.HTTPClient.Transport.(*LogRoundTripper)
			if !ok || lrt.MaxRetries != 1 {
				t.Errorf("the HTTP client shared with the other services was modified")
			}
		})
	}
}
