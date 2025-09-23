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

func TestCheckObsEndpoint(t *testing.T) {
	cfg := &Config{
		Region: "region-0",
		Cloud:  "myhuaweicloud.com",
	}

	// without customizing OBS endpoint in Config
	expected := "https://obs.region-1.myhuaweicloud.com/"
	th.AssertEquals(t, expected, getObsEndpoint(cfg, "region-1"))

	// with customizing OBS endpoint in Config
	cfg.Endpoints = map[string]string{
		"obs": "https://oss.region-0.myhuaweicloud.com/",
	}

	// the region is equal to the region in customizing endpoint
	expected = "https://oss.region-0.myhuaweicloud.com/"
	th.AssertEquals(t, expected, getObsEndpoint(cfg, "region-0"))

	// the region is not equal to the region in customizing endpoint
	expected = "https://oss.region-1.myhuaweicloud.com/"
	th.AssertEquals(t, expected, getObsEndpoint(cfg, "region-1"))
}

// This method is used to mix the generation interference of the test client.
func TestCheckNewServiceClientWithDerivedAuth(t *testing.T) {
	var (
		cfg = &Config{
			Endpoints: map[string]string{
				"iotda": "xxxxxxxxx.st1.iotda-app.cn-north-4.myhuaweicloud.com",
			},
			Region:    "region-0",
			AccessKey: "access key",
			SecretKey: "security key",
			RPLock:    new(sync.Mutex),
			RegionProjectIDMap: map[string]string{
				"region-0": "project ID",
			},
			HwClient: new(golangsdk.ProviderClient),
		}
		region = "region-0"
	)

	client1, err := cfg.NewServiceClient("evs", region)
	if err != nil {
		t.Errorf("error creating EVS client: %s", err)
	}

	client2, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, true)
	if err != nil {
		t.Errorf("error creating IoTDA client: %s", err)
	}

	client3, err := cfg.NewServiceClient("evs", region)
	if err != nil {
		t.Errorf("error creating EVS client: %s", err)
	}

	client4, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, true)
	if err != nil {
		t.Errorf("error creating IoTDA client: %s", err)
	}

	th.AssertEquals(t, false, client1.AKSKAuthOptions.IsDerived)
	th.AssertEquals(t, true, client2.AKSKAuthOptions.IsDerived)
	th.AssertEquals(t, false, client3.AKSKAuthOptions.IsDerived)
	th.AssertEquals(t, true, client4.AKSKAuthOptions.IsDerived)
}
