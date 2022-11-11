package config

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/chnsz/golangsdk"
	th "github.com/chnsz/golangsdk/testhelper"
)

func testRequestRetry(t *testing.T, count int) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	retryCount := count

	var info = struct {
		retries int
		mut     *sync.RWMutex
	}{
		0,
		new(sync.RWMutex),
	}

	th.Mux.HandleFunc("/route/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		_, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Error hadling test request")
		}
		if info.retries < retryCount {
			info.mut.RLock()
			info.retries++
			info.mut.RUnlock()
			//lintignore:R009
			panic(err) // simulate EOF
		}
		w.WriteHeader(500)
		_, _ = fmt.Fprintf(w, `%v`, info.retries)
	})

	cfg := &Config{MaxRetries: retryCount}
	_, err := genClient(cfg, golangsdk.AuthOptions{
		IdentityEndpoint: fmt.Sprintf("%s/route", th.Endpoint()),
	})
	_, ok := err.(golangsdk.ErrDefault500)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, retryCount, info.retries)
}

func TestRequestRetry(t *testing.T) {
	t.Run("TestRequestMultipleRetries", func(t *testing.T) { testRequestRetry(t, 5) })
	t.Run("TestRequestSingleRetry", func(t *testing.T) { testRequestRetry(t, 1) })
	t.Run("TestRequestZeroRetry", func(t *testing.T) { testRequestRetry(t, 0) })
}
