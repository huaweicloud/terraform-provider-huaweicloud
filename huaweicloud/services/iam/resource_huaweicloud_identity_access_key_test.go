package iam

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

const testV3AccessKeyMetadata = `{"credential":{"access":"test-access-key",` +
	`"status":"active","create_time":"2026-01-01T00:00:00Z"}}`

type accessKeyTestAPI struct {
	creates      int32
	reads        int32
	createStatus int
	readStatuses []int
	onRead       func(*http.Request)
}

type accessKeyTestRoundTripper func(*http.Request) (*http.Response, error)

func (f accessKeyTestRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func newAccessKeyTestConfig(t *testing.T, api *accessKeyTestAPI) *config.Config {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method + " " + r.URL.Path {
		case "GET /v3.0/OS-USER/users/test-user":
			_, _ = fmt.Fprint(w, `{"user":{"id":"test-user","name":"test-user"}}`)
		case "POST /v3.0/OS-CREDENTIAL/credentials":
			atomic.AddInt32(&api.creates, 1)
			if api.createStatus != 0 {
				w.WriteHeader(api.createStatus)
				return
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = fmt.Fprint(w, `{"credential":{"access":"test-access-key","secret":"test-secret","user_id":"test-user"}}`)
		case "GET /v3.0/OS-CREDENTIAL/credentials/test-access-key":
			attempt := int(atomic.AddInt32(&api.reads, 1))
			if api.onRead != nil {
				api.onRead(r)
			}
			status := http.StatusOK
			if len(api.readStatuses) > 0 {
				index := attempt - 1
				if index >= len(api.readStatuses) {
					index = len(api.readStatuses) - 1
				}
				status = api.readStatuses[index]
			}
			w.WriteHeader(status)
			if status != http.StatusOK {
				_, _ = fmt.Fprint(w, `{"error":{"message":"test error"}}`)
				return
			}
			_, _ = fmt.Fprint(w, testV3AccessKeyMetadata)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	t.Cleanup(server.Close)
	client := server.Client()
	client.Transport = &config.LogRoundTripper{Rt: client.Transport, MaxRetries: 5}
	return &config.Config{
		DomainClient: &golangsdk.ProviderClient{HTTPClient: *client},
		Endpoints:    map[string]string{"iam": server.URL + "/"},
	}
}

func applyTestV3AccessKey(ctx context.Context, t *testing.T, cfg *config.Config,
	values map[string]interface{}) (*terraform.InstanceState, diag.Diagnostics) {
	t.Helper()
	if values == nil {
		values = make(map[string]interface{})
	}
	values["user_id"] = "test-user"
	if _, ok := values["secret_file"]; !ok {
		values["secret_file"] = filepath.Join(t.TempDir(), "credentials.csv")
	}
	resource := ResourceV3AccessKey()
	diff, err := resource.Diff(ctx, nil, terraform.NewResourceConfigRaw(values), cfg)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	return resource.Apply(ctx, nil, diff, cfg)
}

func TestResourceV3AccessKeyCreate_retriesNotFound(t *testing.T) {
	api := &accessKeyTestAPI{readStatuses: []int{404, 404, 200, 404}}
	cfg := newAccessKeyTestConfig(t, api)
	state, diags := applyTestV3AccessKey(context.Background(), t, cfg, nil)
	assert.Empty(t, diags)
	if !assert.NotNil(t, state, "successful POST must not lose resource state") {
		t.FailNow()
	}
	assert.Equal(t, "test-access-key", state.ID)
	assert.Equal(t, "active", state.Attributes["status"])
	assert.Equal(t, "2026-01-01T00:00:00Z", state.Attributes["create_time"])
	assert.EqualValues(t, 1, atomic.LoadInt32(&api.creates))
	// A fourth, unprotected GET would clear the ID again.
	assert.EqualValues(t, 3, atomic.LoadInt32(&api.reads))
}

func TestResourceV3AccessKeyCreate_reauthentication(t *testing.T) {
	for _, tc := range []struct {
		name          string
		readStatuses  []int
		reauthErr     error
		createTimeout string
		wantReads     int32
		wantError     string
	}{
		{
			name: "not_found_then_available", readStatuses: []int{401, 404, 200, 404},
			createTimeout: "5s", wantReads: 3,
		},
		{
			name: "forbidden", readStatuses: []int{401, 403, 200}, createTimeout: "5s", wantReads: 2,
			wantError: "Successfully re-authenticated, but got error executing request: Action forbidden",
		},
		{
			name: "server_error", readStatuses: []int{401, 500, 200}, createTimeout: "5s", wantReads: 2,
			wantError: "Successfully re-authenticated, but got error executing request: Internal Server Error",
		},
		{
			name: "reauthentication_failed", readStatuses: []int{401, 200}, reauthErr: errors.New("test reauthentication failed"),
			createTimeout: "5s", wantReads: 1, wantError: "Unable to re-authenticate",
		},
		{
			name: "persistent_not_found", readStatuses: []int{401, 404}, createTimeout: "2s", wantReads: 3,
			wantError: "deadline exceeded|timeout while waiting",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			api := &accessKeyTestAPI{readStatuses: tc.readStatuses}
			cfg := newAccessKeyTestConfig(t, api)
			var reauths int32
			cfg.DomainClient.ReauthFunc = func() error {
				atomic.AddInt32(&reauths, 1)
				return tc.reauthErr
			}
			start := time.Now()
			state, diags := applyTestV3AccessKey(context.Background(), t, cfg, map[string]interface{}{
				"timeouts": map[string]interface{}{"create": tc.createTimeout},
			})
			assert.Less(t, time.Since(start), 3*time.Second, "must not retry permanent errors or exceed the verification budget")
			if !assert.NotNil(t, state, "SDK must retain the created key even when verification fails") {
				t.FailNow()
			}
			assert.Equal(t, "test-access-key", state.ID)
			assert.EqualValues(t, 1, atomic.LoadInt32(&api.creates))
			assert.Equal(t, tc.wantReads, atomic.LoadInt32(&api.reads))
			assert.EqualValues(t, 1, atomic.LoadInt32(&reauths))
			if tc.wantError == "" {
				assert.Empty(t, diags)
				assert.Equal(t, "active", state.Attributes["status"])
				assert.Equal(t, "2026-01-01T00:00:00Z", state.Attributes["create_time"])
			} else {
				assert.True(t, diags.HasError())
				message := fmt.Sprint(diags)
				assert.Regexp(t, tc.wantError, message)
				assert.Contains(t, message, "The key was created and its ID has been retained")
				assert.NotContains(t, message, "test-access-key")
				assert.NotContains(t, message, "test-secret")
			}
		})
	}
}

func TestResourceV3AccessKeyCreate_immediateSuccess(t *testing.T) {
	api := &accessKeyTestAPI{}
	cfg := newAccessKeyTestConfig(t, api)
	secretFile := filepath.Join(t.TempDir(), "credentials.csv")
	state, diags := applyTestV3AccessKey(context.Background(), t, cfg, map[string]interface{}{"secret_file": secretFile})
	assert.Empty(t, diags)
	if !assert.NotNil(t, state) {
		t.FailNow()
	}
	assert.Equal(t, "test-access-key", state.ID)
	assert.Equal(t, "test-user", state.Attributes["user_name"])
	assert.Equal(t, "active", state.Attributes["status"])
	assert.Empty(t, state.Attributes["secret"])
	assert.EqualValues(t, 1, atomic.LoadInt32(&api.creates))
	assert.EqualValues(t, 1, atomic.LoadInt32(&api.reads))
	contents, err := os.ReadFile(secretFile)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.Contains(t, string(contents), "test-user,test-access-key,test-secret")
	info, err := os.Stat(secretFile)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
}

func TestResourceV3AccessKeyCreate_verificationErrorsPreserveState(t *testing.T) {
	for _, tc := range []struct {
		name    string
		status  int
		timeout string
		message string
	}{
		{name: "persistent_not_found", status: 404, timeout: "100ms", message: "deadline exceeded|timeout while waiting"},
		{name: "forbidden", status: 403, timeout: "5s", message: "Action forbidden"},
		{name: "server_error", status: 500, timeout: "5s", message: "Internal Server Error"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			api := &accessKeyTestAPI{readStatuses: []int{tc.status}}
			cfg := newAccessKeyTestConfig(t, api)
			start := time.Now()
			state, diags := applyTestV3AccessKey(context.Background(), t, cfg, map[string]interface{}{
				"timeouts": map[string]interface{}{"create": tc.timeout},
			})
			assert.True(t, diags.HasError(), "expected verification error, got %v", diags)
			if !assert.NotNil(t, state, "SDK must return the created state alongside diagnostics") {
				t.FailNow()
			}
			assert.Equal(t, "test-access-key", state.ID)
			assert.Equal(t, "test-user", state.Attributes["user_name"])
			assert.EqualValues(t, 1, atomic.LoadInt32(&api.creates))
			assert.EqualValues(t, 1, atomic.LoadInt32(&api.reads))
			assert.Less(t, time.Since(start), 3*time.Second, "must honor timeout and fail fast on non-404 errors")
			message := fmt.Sprint(diags)
			assert.Regexp(t, tc.message, message)
			assert.Contains(t, message, "The key was created and its ID has been retained")
			assert.NotContains(t, message, "test-access-key")
			assert.NotContains(t, message, "test-secret")
		})
	}
}

func TestResourceV3AccessKeyCreate_postFailure(t *testing.T) {
	api := &accessKeyTestAPI{createStatus: http.StatusForbidden}
	cfg := newAccessKeyTestConfig(t, api)
	state, diags := applyTestV3AccessKey(context.Background(), t, cfg, nil)
	assert.True(t, diags.HasError())
	assert.Nil(t, state)
	assert.EqualValues(t, 1, atomic.LoadInt32(&api.creates))
	assert.Zero(t, atomic.LoadInt32(&api.reads))
}

func TestResourceV3AccessKeyCreate_cancelWhileWaiting(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	api := &accessKeyTestAPI{
		readStatuses: []int{http.StatusNotFound},
		onRead:       func(_ *http.Request) { cancel() },
	}
	cfg := newAccessKeyTestConfig(t, api)
	start := time.Now()
	state, diags := applyTestV3AccessKey(ctx, t, cfg, nil)
	assert.True(t, diags.HasError())
	assert.Contains(t, fmt.Sprint(diags), "context canceled")
	if !assert.NotNil(t, state) {
		t.FailNow()
	}
	assert.Equal(t, "test-access-key", state.ID)
	assert.EqualValues(t, 1, atomic.LoadInt32(&api.creates))
	assert.EqualValues(t, 1, atomic.LoadInt32(&api.reads))
	assert.Less(t, time.Since(start), time.Second)
}

func TestResourceV3AccessKeyCreate_timeoutCancelsRequest(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	requestCanceled := make(chan struct{})
	transportReturned := make(chan struct{})
	releaseHandler := make(chan struct{})
	// Release the handler before server.Close runs, even if a cancellation assertion fails.
	defer close(releaseHandler)
	api := &accessKeyTestAPI{
		onRead: func(r *http.Request) {
			select {
			case <-r.Context().Done():
				close(requestCanceled)
			case <-releaseHandler:
				// Cleanup must not be mistaken for request cancellation.
			}
		},
	}
	cfg := newAccessKeyTestConfig(t, api)
	transport := cfg.DomainClient.HTTPClient.Transport
	cfg.DomainClient.HTTPClient.Transport = accessKeyTestRoundTripper(func(r *http.Request) (*http.Response, error) {
		if r.Method == http.MethodGet && r.URL.Path == "/v3.0/OS-CREDENTIAL/credentials/test-access-key" {
			defer close(transportReturned)
		}
		return transport.RoundTrip(r)
	})
	state, diags := applyTestV3AccessKey(ctx, t, cfg, nil)
	assert.True(t, diags.HasError())
	assert.NoError(t, ctx.Err(), "verification must return before the outer create deadline")
	if !assert.NotNil(t, state) {
		t.FailNow()
	}
	assert.Equal(t, "test-access-key", state.ID)
	select {
	case <-requestCanceled:
	case <-time.After(time.Second):
		t.Fatal("create timeout did not cancel the in-flight IAM request")
	}
	select {
	case <-transportReturned:
	case <-time.After(time.Second):
		t.Fatal("HTTP retry transport continued running after create timed out")
	}
	assert.Nil(t, cfg.DomainClient.Context, "operation context must not leak into the shared client")
}

func TestResourceV3AccessKeyCreate_timeoutBudgets(t *testing.T) {
	for _, tc := range []struct {
		name          string
		createTimeout string
		parentTimeout time.Duration
		createBudget  time.Duration
		readBudget    time.Duration
	}{
		{name: "default", createBudget: 20 * time.Minute, readBudget: time.Minute},
		{name: "custom", createTimeout: "5m", createBudget: 5 * time.Minute, readBudget: time.Minute},
		{name: "short", createTimeout: "5s", createBudget: 5 * time.Second, readBudget: 4500 * time.Millisecond},
		{name: "parent_deadline", parentTimeout: 5 * time.Second, createBudget: 5 * time.Second, readBudget: 4500 * time.Millisecond},
	} {
		t.Run(tc.name, func(t *testing.T) {
			api := &accessKeyTestAPI{}
			cfg := newAccessKeyTestConfig(t, api)
			budgets := make(chan time.Duration, 3)
			transport := cfg.DomainClient.HTTPClient.Transport
			cfg.DomainClient.HTTPClient.Transport = accessKeyTestRoundTripper(func(r *http.Request) (*http.Response, error) {
				deadline, ok := r.Context().Deadline()
				assert.True(t, ok, "every IAM request must have the appropriate operation deadline")
				select {
				case budgets <- time.Until(deadline):
				default:
					t.Error("unexpected extra IAM request")
				}
				return transport.RoundTrip(r)
			})
			values := make(map[string]interface{})
			if tc.createTimeout != "" {
				values["timeouts"] = map[string]interface{}{"create": tc.createTimeout}
			}
			ctx := context.Background()
			if tc.parentTimeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, tc.parentTimeout)
				defer cancel()
			}
			state, diags := applyTestV3AccessKey(ctx, t, cfg, values)
			assert.Empty(t, diags)
			assert.NotNil(t, state)
			if !assert.Len(t, budgets, 3) {
				t.FailNow()
			}
			for _, expected := range []time.Duration{tc.createBudget, tc.createBudget, tc.readBudget} {
				actual := <-budgets
				assert.LessOrEqual(t, actual, expected)
				assert.Greater(t, actual, expected-250*time.Millisecond)
			}
		})
	}
}

func TestResourceV3AccessKeyCreate_csvFallbackOnVerificationError(t *testing.T) {
	api := &accessKeyTestAPI{readStatuses: []int{http.StatusForbidden}}
	cfg := newAccessKeyTestConfig(t, api)
	state, diags := applyTestV3AccessKey(context.Background(), t, cfg, map[string]interface{}{
		"secret_file": filepath.Join(t.TempDir(), "missing", "credentials.csv"),
	})
	assert.True(t, diags.HasError())
	assert.Len(t, diags, 2)
	assert.Equal(t, diag.Warning, diags[0].Severity)
	if !assert.NotNil(t, state) {
		t.FailNow()
	}
	assert.Equal(t, "test-access-key", state.ID)
	assert.Equal(t, "test-secret", state.Attributes["secret"])
}

func TestResourceV3AccessKeyCreate_pgpWithRetry(t *testing.T) {
	entity, err := openpgp.NewEntity("test", "", "test@example.com", nil)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	var publicKey bytes.Buffer
	if !assert.NoError(t, entity.Serialize(&publicKey)) {
		t.FailNow()
	}
	api := &accessKeyTestAPI{readStatuses: []int{404, 200}}
	cfg := newAccessKeyTestConfig(t, api)
	state, diags := applyTestV3AccessKey(context.Background(), t, cfg, map[string]interface{}{
		"pgp_key": base64.StdEncoding.EncodeToString(publicKey.Bytes()),
	})
	assert.Empty(t, diags)
	if !assert.NotNil(t, state) {
		t.FailNow()
	}
	assert.NotEmpty(t, state.Attributes["key_fingerprint"])
	assert.Empty(t, state.Attributes["secret"])
	encrypted, err := base64.StdEncoding.DecodeString(state.Attributes["encrypted_secret"])
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	message, err := openpgp.ReadMessage(bytes.NewReader(encrypted), openpgp.EntityList{entity}, nil, nil)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	secret, err := io.ReadAll(message.UnverifiedBody)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.Equal(t, "test-secret", string(secret))
	assert.EqualValues(t, 1, atomic.LoadInt32(&api.creates))
	assert.EqualValues(t, 2, atomic.LoadInt32(&api.reads))
}

func TestResourceV3AccessKeyRead_deleted(t *testing.T) {
	api := &accessKeyTestAPI{readStatuses: []int{http.StatusNotFound}}
	cfg := newAccessKeyTestConfig(t, api)
	resource := ResourceV3AccessKey()
	data := resource.Data(&terraform.InstanceState{ID: "test-access-key"})
	diags := resource.ReadContext(context.Background(), data, cfg)
	assert.False(t, diags.HasError())
	assert.Empty(t, data.Id())
	assert.Len(t, diags, 1)
	assert.Equal(t, diag.Warning, diags[0].Severity)
	assert.Contains(t, diags[0].Detail, "removed in Terraform state")
	assert.EqualValues(t, 1, atomic.LoadInt32(&api.reads))
	assert.Zero(t, atomic.LoadInt32(&api.creates))
}
