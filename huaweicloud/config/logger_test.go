package config

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type logTestRoundTripper func(*http.Request) (*http.Response, error)

func (f logTestRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestLogRoundTripper_canceledBeforeRequest(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", nil)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	attempts := 0
	transport := &LogRoundTripper{
		MaxRetries: 1,
		Rt: logTestRoundTripper(func(_ *http.Request) (*http.Response, error) {
			attempts++
			return nil, ctx.Err()
		}),
	}
	start := time.Now()
	response, err := transport.RoundTrip(request)
	assert.Nil(t, response)
	assert.ErrorIs(t, err, context.Canceled)
	assert.Zero(t, attempts, "an already canceled request must not reach the underlying transport")
	assert.Less(t, time.Since(start), time.Second)
}

func TestLogRoundTripper_canceledDuringBackoff(t *testing.T) {
	for _, withResponse := range []bool{false, true} {
		name := "connection_error"
		if withResponse {
			name = "connection_reset_with_response"
		}
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", nil)
			if !assert.NoError(t, err) {
				t.FailNow()
			}
			attempts := 0
			transport := &LogRoundTripper{
				MaxRetries: 1,
				Rt: logTestRoundTripper(func(_ *http.Request) (*http.Response, error) {
					attempts++
					if attempts == 1 {
						time.AfterFunc(10*time.Millisecond, cancel)
					}
					if withResponse {
						return &http.Response{StatusCode: http.StatusBadGateway}, errors.New("connection reset by peer")
					}
					return nil, io.EOF
				}),
			}
			start := time.Now()
			response, err := transport.RoundTrip(request)
			assert.Nil(t, response)
			assert.ErrorIs(t, err, context.Canceled)
			assert.Equal(t, 1, attempts, "cancellation must interrupt backoff before another request")
			assert.Less(t, time.Since(start), time.Second)
		})
	}
}

func TestLogRoundTripper_deadlineDuringBackoff(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", nil)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	attempts := 0
	transport := &LogRoundTripper{
		MaxRetries: 1,
		Rt: logTestRoundTripper(func(_ *http.Request) (*http.Response, error) {
			attempts++
			return nil, io.EOF
		}),
	}
	start := time.Now()
	response, err := transport.RoundTrip(request)
	assert.Nil(t, response)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Equal(t, 1, attempts)
	assert.Less(t, time.Since(start), time.Second)
}

func TestLogRoundTripper_connectionResetRecovers(t *testing.T) {
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com", nil)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	attempts := 0
	transport := &LogRoundTripper{
		MaxRetries: 1,
		Rt: logTestRoundTripper(func(_ *http.Request) (*http.Response, error) {
			attempts++
			if attempts == 1 {
				return &http.Response{StatusCode: http.StatusBadGateway}, errors.New("connection reset by peer")
			}
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("ok"))}, nil
		}),
	}
	response, err := transport.RoundTrip(request)
	if !assert.NoError(t, err) || !assert.NotNil(t, response) {
		t.FailNow()
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.Equal(t, "ok", string(body))
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 2, attempts)
}
