package common

import (
	"crypto/tls"
	"net/http"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// NewCustomClient creates a custom client assembled from user-provided endpoints.
// URLs will be assembled according to the endpoints array, separated each element by slashes.
// for example, array ["https://www.example.com", "v2", "test", ...] will form the address
// "https://www.example.com/v2/test/.../".
// NOTE: You can decide whether to skip the SSL certificate check with the insecure parameter.
func NewCustomClient(insecure bool, endpoints ...string) *golangsdk.ServiceClient {
	p := new(golangsdk.ProviderClient)
	p.HTTPClient = http.Client{
		Transport: &config.LogRoundTripper{
			Rt: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				TLSClientConfig: &tls.Config{
					// The version fields of the TLS ver.1.2 and ver.1.3 are filled with 0x0303 (ver.1.2)
					MinVersion:         tls.VersionTLS12,
					InsecureSkipVerify: insecure, // Pay attention to the security risks after skip verify.
				},
			},
		},
		Timeout: 30 * time.Minute,
	}

	return &golangsdk.ServiceClient{
		ProviderClient: p,
		ResourceBase:   strings.Join(endpoints, "/") + "/",
	}
}
