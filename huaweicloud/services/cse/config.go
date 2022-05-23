package cse

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cse/dedicated/v4/auth"
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

// getAuthorizationToken is a method to request the CSE API and get the authorization token.
// The format of is "Bearer {token}".
func getAuthorizationToken(c *golangsdk.ServiceClient, username, password string) (string, error) {
	tokenOpts := auth.CreateOpts{
		Name:     username,
		Password: password,
	}
	resp, err := auth.Create(c, tokenOpts)
	if err != nil {
		return "", fmt.Errorf("unable to create the authorization token: %v", err)
	}

	return fmt.Sprintf("Bearer %s", resp.Token), nil
}

// GetAuthorizationToken is a method for creating an authorization information for a CSE microservice to connect to the
// specified dedicated engine.
func GetAuthorizationToken(connAddr, username, password string) (string, error) {
	if username == "" {
		return "", nil
	}
	client := NewCustomClient(connAddr, "v4")
	return getAuthorizationToken(client, username, password)
}
