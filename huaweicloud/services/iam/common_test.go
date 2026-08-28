package iam

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/chnsz/golangsdk"
)

// TestPollRequest_successOnFirstRequest verifies the happy path where the request function succeeds
// on its first invocation after the initial delay.
//
// PollRequest always waits iamPollInitialDelay before calling requestFunc for the first time.
// This test ensures that:
//   - A successful response is returned without error when requestFunc succeeds immediately.
//   - The returned value matches exactly what requestFunc produced.
//   - The total elapsed time is at least iamPollInitialDelay, confirming the initial delay is enforced
//     even when no retry is required.
func TestPollRequest_successOnFirstRequest(t *testing.T) {
	start := time.Now()
	expected := map[string]interface{}{"id": "test-id"}

	resp, err := PollRequest(context.Background(), func() (interface{}, error) {
		return expected, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !reflect.DeepEqual(resp, expected) {
		t.Fatalf("expected %#v, got %#v", expected, resp)
	}

	// Confirm the initial delay occurred before the first request was made.
	elapsed := time.Since(start)
	if elapsed < iamPollInitialDelay {
		t.Fatalf("expected initial delay of at least %s, elapsed %s", iamPollInitialDelay, elapsed)
	}
}

// TestPollRequest_successAfter404Retries verifies that PollRequest keeps polling when requestFunc
// returns golangsdk.ErrDefault404, and eventually returns the successful response once the resource
// becomes available.
//
// The mock requestFunc returns 404 for the first two attempts and succeeds on the third.
// This test ensures that:
//   - PollRequest retries on 404 instead of treating it as a terminal error.
//   - requestFunc is invoked exactly three times before returning.
//   - The final successful response is returned to the caller.
//   - The total elapsed time includes the initial delay plus one interval wait after each 404,
//     i.e. iamPollInitialDelay + 2*iamPollInterval.
func TestPollRequest_successAfter404Retries(t *testing.T) {
	start := time.Now()
	expected := map[string]interface{}{"id": "test-id"}
	attempts := 0

	resp, err := PollRequest(context.Background(), func() (interface{}, error) {
		attempts++
		if attempts < 3 {
			return nil, golangsdk.ErrDefault404{}
		}
		return expected, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !reflect.DeepEqual(resp, expected) {
		t.Fatalf("expected %#v, got %#v", expected, resp)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}

	// Two 404 responses require two interval waits between the three request attempts.
	minElapsed := iamPollInitialDelay + 2*iamPollInterval
	if elapsed := time.Since(start); elapsed < minElapsed {
		t.Fatalf("expected elapsed time of at least %s, got %s", minElapsed, elapsed)
	}
}

// TestPollRequest_non404Error verifies that PollRequest stops immediately when requestFunc returns
// an error other than golangsdk.ErrDefault404.
//
// Non-404 errors (e.g. 403, 500, network failures) are considered non-retryable and should be
// returned directly to the caller without entering the interval retry loop.
// This test ensures that the original error is preserved and returned unchanged.
func TestPollRequest_non404Error(t *testing.T) {
	expectedErr := errors.New("internal error")

	_, err := PollRequest(context.Background(), func() (interface{}, error) {
		return nil, expectedErr
	})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

// TestPollRequest_contextCanceledDuringInitialDelay verifies that PollRequest respects context
// cancellation during the initial delay phase, before requestFunc is ever invoked.
//
// The context is canceled before PollRequest is called. This test ensures that:
//   - PollRequest returns immediately with context.Canceled instead of waiting the full initial delay.
//   - requestFunc is never called when the context is already canceled at entry.
func TestPollRequest_contextCanceledDuringInitialDelay(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := PollRequest(ctx, func() (interface{}, error) {
		t.Fatal("request function should not be called when context is already canceled")
		return nil, nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

// TestPollRequest_contextCanceledDuringInterval verifies that PollRequest respects context
// cancellation while waiting between retry attempts after a 404 response.
//
// The mock requestFunc returns 404 on the first call and triggers context cancellation at that
// point. PollRequest should detect the canceled context during the subsequent interval wait and
// return context.Canceled instead of continuing to poll.
// A goroutine is used because PollRequest blocks until it returns or the test times out.
func TestPollRequest_contextCanceledDuringInterval(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	attempts := 0

	done := make(chan error, 1)
	go func() {
		_, err := PollRequest(ctx, func() (interface{}, error) {
			attempts++
			if attempts == 1 {
				cancel()
			}
			return nil, golangsdk.ErrDefault404{}
		})
		done <- err
	}()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for PollRequest to return after context cancellation")
	}
}

// TestPollRequest_timeout verifies that PollRequest eventually returns a timeout error when
// requestFunc keeps returning golangsdk.ErrDefault404 until the polling deadline is exceeded.
//
// The polling deadline is iamPollTimeout from the moment PollRequest starts (after the initial
// delay is set up). This test ensures that:
//   - A non-nil error is returned when the resource never becomes available within the timeout window.
//   - requestFunc is called multiple times before giving up, confirming retries did occur.
//   - The total elapsed time is at least iamPollInitialDelay + iamPollTimeout.
//
// This test is skipped when running with -short because it requires approximately 21 seconds to complete.
func TestPollRequest_timeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping timeout test in short mode")
	}

	start := time.Now()
	attempts := 0

	_, err := PollRequest(context.Background(), func() (interface{}, error) {
		attempts++
		return nil, golangsdk.ErrDefault404{}
	})
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if attempts < 2 {
		t.Fatalf("expected multiple polling attempts before timeout, got %d", attempts)
	}

	minElapsed := iamPollInitialDelay + iamPollTimeout
	if elapsed := time.Since(start); elapsed < minElapsed {
		t.Fatalf("expected elapsed time of at least %s, got %s", minElapsed, elapsed)
	}
}

// TestPollRequest_longRunningRequestExceedsDeadline verifies that PollRequest does not interrupt
// an in-flight requestFunc, even when that call blocks longer than iamPollTimeout.
//
// This test simulates the SDK 429 backoff scenario by sleeping 25 seconds inside requestFunc
// (longer than the 20-second polling deadline) and then returning 404.
//
// Expected timeline:
//   - t=0s:  PollRequest starts
//   - t=1s:  initial delay ends, requestFunc is invoked
//   - t=1s~26s: requestFunc blocks for 25 seconds (PollRequest cannot check the deadline yet)
//   - t=26s: requestFunc returns 404, PollRequest then detects that the deadline has already
//     passed and returns a timeout error
//
// Therefore the total elapsed time should fall in the ~26s range
// (iamPollInitialDelay + 25s), NOT the ~21s range (iamPollInitialDelay + iamPollTimeout).
// This confirms that requestFunc execution time is not preempted by the PollRequest deadline.
//
// This test is skipped when running with -short because it requires approximately 26 seconds.
func TestPollRequest_longRunningRequestExceedsDeadline(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running requestFunc test in short mode")
	}

	const requestSleep = 25 * time.Second
	start := time.Now()
	attempts := 0

	_, err := PollRequest(context.Background(), func() (interface{}, error) {
		attempts++
		// Simulate a long-blocking request (e.g. SDK 429 backoff).
		// lintignore:R018
		time.Sleep(requestSleep)
		return nil, golangsdk.ErrDefault404{}
	})

	elapsed := time.Since(start)
	t.Logf("PollRequest returned after %s (attempts=%d, err=%v)", elapsed, attempts, err)

	if err == nil {
		t.Fatal("expected timeout error after the long-running requestFunc returned 404, got nil")
	}
	if attempts != 1 {
		t.Fatalf("expected requestFunc to be called exactly once, got %d", attempts)
	}

	// ~21s range would mean PollRequest somehow timed out around its own deadline without
	// waiting for requestFunc. That should not happen.
	around21s := iamPollInitialDelay + iamPollTimeout // 21s
	// ~26s range means PollRequest waited for the full requestFunc sleep before returning.
	around26s := iamPollInitialDelay + requestSleep // 26s

	if elapsed < around26s {
		t.Fatalf("elapsed %s is closer to the 21s range (%s) than expected; "+
			"PollRequest should wait for requestFunc to finish (~%s)",
			elapsed, around21s, around26s)
	}

	// Allow a small scheduling overhead upper bound.
	if elapsed > around26s+2*time.Second {
		t.Fatalf("elapsed %s is unexpectedly larger than ~%s", elapsed, around26s)
	}
}
