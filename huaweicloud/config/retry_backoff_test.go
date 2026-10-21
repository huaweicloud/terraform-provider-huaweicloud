package config

import (
	"testing"
	"time"

	"github.com/chnsz/golangsdk"
	th "github.com/chnsz/golangsdk/testhelper"
)

func TestLookupRetryInitialBackoff(t *testing.T) {
	d, ok := lookupRetryInitialBackoff("IAM")
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, time.Second, d)

	_, ok = lookupRetryInitialBackoff("ECS")
	th.AssertEquals(t, false, ok)

	_, ok = lookupRetryInitialBackoff("")
	th.AssertEquals(t, false, ok)
}

func TestCalculateRetryBackoffSleep(t *testing.T) {
	// Default initial (1m): historical behavior.
	th.AssertEquals(t, 1*time.Minute, calculateRetryBackoffSleep(defaultRetryInitialBackoff, 0))
	th.AssertEquals(t, 2*time.Minute, calculateRetryBackoffSleep(defaultRetryInitialBackoff, 1))
	th.AssertEquals(t, 4*time.Minute, calculateRetryBackoffSleep(defaultRetryInitialBackoff, 2))
	th.AssertEquals(t, 8*time.Minute, calculateRetryBackoffSleep(defaultRetryInitialBackoff, 3))
	th.AssertEquals(t, 16*time.Minute, calculateRetryBackoffSleep(defaultRetryInitialBackoff, 4))
	th.AssertEquals(t, 30*time.Minute, calculateRetryBackoffSleep(defaultRetryInitialBackoff, 5))
	th.AssertEquals(t, 30*time.Minute, calculateRetryBackoffSleep(defaultRetryInitialBackoff, 10))

	// IAM-style initial (1s).
	th.AssertEquals(t, 1*time.Second, calculateRetryBackoffSleep(time.Second, 0))
	th.AssertEquals(t, 2*time.Second, calculateRetryBackoffSleep(time.Second, 1))
	th.AssertEquals(t, 4*time.Second, calculateRetryBackoffSleep(time.Second, 2))
	th.AssertEquals(t, 8*time.Second, calculateRetryBackoffSleep(time.Second, 3))

	// Invalid initial falls back to default base.
	th.AssertEquals(t, 1*time.Minute, calculateRetryBackoffSleep(0, 0))
}

func TestApplyServiceRetryBackoff(t *testing.T) {
	client := &golangsdk.ProviderClient{
		RetryBackoffFunc:  newRetryBackoffFunc(defaultRetryInitialBackoff),
		MaxBackoffRetries: 5,
	}

	// Non-configured product is a no-op (func remains set).
	applyServiceRetryBackoff(client, "ECS")
	if client.RetryBackoffFunc == nil {
		t.Fatal("expected RetryBackoffFunc to remain set for non-configured product")
	}

	// IAM product overrides RetryBackoffFunc.
	applyServiceRetryBackoff(client, "IAM")
	if client.RetryBackoffFunc == nil {
		t.Fatal("expected IAM RetryBackoffFunc to be set")
	}

	// Nil client / nil func are no-ops.
	applyServiceRetryBackoff(nil, "IAM")
	empty := &golangsdk.ProviderClient{}
	applyServiceRetryBackoff(empty, "IAM")
	if empty.RetryBackoffFunc != nil {
		t.Fatal("expected no override when RetryBackoffFunc was nil")
	}
}

func TestLookupRetryInitialBackoffInvalidDuration(t *testing.T) {
	original := serviceRetryInitialBackoff["IAM"]
	serviceRetryInitialBackoff["IAM"] = "not-a-duration"
	defer func() {
		serviceRetryInitialBackoff["IAM"] = original
	}()

	d, ok := lookupRetryInitialBackoff("IAM")
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, defaultRetryInitialBackoff, d)
}
