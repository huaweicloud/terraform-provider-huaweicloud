package common

import (
	"crypto/tls"
	"net/http"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// NewCustomClient creates a custom client assembled from user-provided endpoints.
// URLs will be assembled according to the endpoints array, separated each element by slashes.
// for example, array ["https://www.example.com", "v2", "test", ...] will form the address
// "https://www.example.com/v2/test/.../".
// NOTE: This client will skip the SSL certificate check.
func NewCustomClient(endpoints ...string) *golangsdk.ServiceClient {
	p := new(golangsdk.ProviderClient)
	p.HTTPClient = http.Client{
		Transport: &config.LogRoundTripper{
			Rt: &http.Transport{
				Proxy:           http.ProxyFromEnvironment,
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			OsDebug: logging.IsDebugOrHigher(),
		},
		Timeout: 30 * time.Minute,
	}

	return &golangsdk.ServiceClient{
		ProviderClient: p,
		ResourceBase:   strings.Join(endpoints, "/") + "/",
	}
}
