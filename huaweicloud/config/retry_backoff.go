package config

import (
	"context"
	"log"
	"time"

	"github.com/chnsz/golangsdk"
)

const (
	// defaultRetryInitialBackoff is the initial sleep used by the default 429 backoff
	// strategy. With retries=0 this preserves the historical "1 minute" first wait.
	defaultRetryInitialBackoff = time.Minute
	// defaultRetryMaxBackoff caps exponential backoff at 30 minutes.
	defaultRetryMaxBackoff = 30 * time.Minute
)

// serviceRetryInitialBackoff maps ServiceCatalog.Product to an initial backoff
// duration string (time.ParseDuration format, e.g. "1s", "5m").
// Only products that need to diverge from the default should be listed here.
// This is a code-side allowlist and is not exposed as a provider schema.
var serviceRetryInitialBackoff = map[string]string{
	"IAM": "1s",
}

// lookupRetryInitialBackoff returns the initial backoff for a service product.
// ok is false when the product is not configured or the duration string is invalid
// (invalid values fall back to the default after logging a warning when used via
// applyServiceRetryBackoff).
func lookupRetryInitialBackoff(product string) (time.Duration, bool) {
	raw, ok := serviceRetryInitialBackoff[product]
	if !ok {
		return 0, false
	}

	d, err := time.ParseDuration(raw)
	if err != nil || d <= 0 {
		log.Printf("[WARN] invalid retry initial backoff %q for product %q: %v, fallback to default %s",
			raw, product, err, defaultRetryInitialBackoff)
		return defaultRetryInitialBackoff, true
	}
	return d, true
}

// calculateRetryBackoffSleep returns the sleep duration for a 429 retry attempt.
// Formula: initial * 2^retries, capped at defaultRetryMaxBackoff.
func calculateRetryBackoffSleep(initial time.Duration, retries uint) time.Duration {
	if initial <= 0 {
		initial = defaultRetryInitialBackoff
	}

	sleep := initial << retries // initial * 2^retries
	if sleep > defaultRetryMaxBackoff || sleep <= 0 {
		// Cap overflow (shift too large) and the configured max.
		return defaultRetryMaxBackoff
	}
	return sleep
}

// newRetryBackoffFunc builds a 429 RetryBackoffFunc with the given initial interval.
// Sleep duration is initial * 2^retries, capped at defaultRetryMaxBackoff.
func newRetryBackoffFunc(initial time.Duration) golangsdk.RetryFunc {
	if initial <= 0 {
		initial = defaultRetryInitialBackoff
	}

	return func(ctx context.Context, _ *golangsdk.ErrUnexpectedResponseCode, e error, retries uint) error {
		sleep := calculateRetryBackoffSleep(initial, retries)

		log.Printf("[WARN] Received StatusTooManyRequests response code, try to sleep %s (initial=%s, retries=%d)",
			sleep, initial, retries)

		if ctx != nil {
			select {
			case <-time.After(sleep):
			case <-ctx.Done():
				return e
			}
		} else {
			// lintignore:R018
			time.Sleep(sleep)
		}

		return nil
	}
}

// applyServiceRetryBackoff overrides the cloned ProviderClient RetryBackoffFunc
// when the service product has a custom initial backoff configured.
func applyServiceRetryBackoff(client *golangsdk.ProviderClient, product string) {
	if client == nil || client.RetryBackoffFunc == nil {
		return
	}

	initial, ok := lookupRetryInitialBackoff(product)
	if !ok {
		return
	}

	client.RetryBackoffFunc = newRetryBackoffFunc(initial)
}
