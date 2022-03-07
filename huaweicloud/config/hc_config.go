package config

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	hc_config "github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/httphandler"
	tms "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/tms/v1"
	vpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

/*
This file is used to impl the configuration of huaweicloud-sdk-go-v3 package and
genetate service clients.
*/
func buildAuthCredentials(c *Config, region string) (*basic.Credentials, error) {
	if c.AccessKey == "" || c.SecretKey == "" {
		return nil, fmt.Errorf("access_key or secret_key is missing in the provider")
	}

	credentials := basic.Credentials{
		AK:            c.AccessKey,
		SK:            c.SecretKey,
		SecurityToken: c.SecurityToken,
		IamEndpoint:   c.IdentityEndpoint,
	}

	c.RPLock.Lock()
	defer c.RPLock.Unlock()
	projectID, ok := c.RegionProjectIDMap[region]
	if !ok {
		// Not find in the map, then try to query and store.
		client := c.HwClient
		err := c.loadUserProjects(client, region)
		if err != nil {
			return nil, err
		}
		projectID, _ = c.RegionProjectIDMap[region]
	}

	credentials.ProjectId = projectID
	return &credentials, nil
}

func buildGlobalAuthCredentials(c *Config, region string) (*global.Credentials, error) {
	if c.AccessKey == "" || c.SecretKey == "" {
		return nil, fmt.Errorf("access_key or secret_key is missing in the provider")
	}

	credentials := global.Credentials{
		AK:            c.AccessKey,
		SK:            c.SecretKey,
		DomainId:      c.DomainID,
		SecurityToken: c.SecurityToken,
		IamEndpoint:   c.IdentityEndpoint,
	}

	return &credentials, nil
}

func buildHTTPConfig(c *Config) *hc_config.HttpConfig {
	httpConfig := hc_config.DefaultHttpConfig()

	if c.MaxRetries > 0 {
		httpConfig = httpConfig.WithRetries(c.MaxRetries)
	}

	if c.Insecure {
		httpConfig = httpConfig.WithIgnoreSSLVerification(true)
	}

	httpHandler := httphandler.NewHttpHandler().
		AddRequestHandler(logRequestHandler).
		AddResponseHandler(logResponseHandler)
	httpConfig = httpConfig.WithHttpHandler(httpHandler)

	if proxyURL := getProxyFromEnv(); proxyURL != "" {
		if parsed, err := url.Parse(proxyURL); err == nil {
			logp.Printf("[DEBUG] using https proxy: %s://%s", parsed.Scheme, parsed.Host)

			httpProxy := hc_config.Proxy{
				Schema:   parsed.Scheme,
				Host:     parsed.Host,
				Username: parsed.User.Username(),
			}
			if pwd, ok := parsed.User.Password(); ok {
				httpProxy.Password = pwd
			}

			httpConfig = httpConfig.WithProxy(&httpProxy)
		} else {
			logp.Printf("[WARN] parsing https proxy failed: %s", err)
		}
	}

	return httpConfig
}

func getServiceEndpoint(c *Config, srv, region string) string {
	// try to get the endpoint from customizing map
	if endpoint, ok := c.Endpoints[srv]; ok {
		return endpoint
	}

	// get the endpoint from build-in catalog
	catalog, ok := allServiceCatalog[srv]
	if !ok {
		return ""
	}

	var ep string
	if catalog.Scope == "global" && !c.RegionClient {
		ep = fmt.Sprintf("https://%s.%s/", catalog.Name, c.Cloud)
	} else {
		ep = fmt.Sprintf("https://%s.%s.%s/", catalog.Name, region, c.Cloud)
	}
	return ep
}

// NewVpcClient is the VPC service client using huaweicloud-sdk-go-v3 package
func NewVpcClient(c *Config, region string) (*vpc.VpcClient, error) {
	credentials, err := buildAuthCredentials(c, region)
	if err != nil {
		return nil, err
	}

	vpcEndpoint := getServiceEndpoint(c, "vpc", region)
	if vpcEndpoint == "" {
		return nil, fmt.Errorf("failed to get the endpoint of VPC service")
	}

	return vpc.NewVpcClient(
		vpc.VpcClientBuilder().
			WithEndpoint(vpcEndpoint).
			WithCredential(*credentials).
			WithHttpConfig(buildHTTPConfig(c)).
			Build()), nil
}

// NewTmsClient is the TMS service client using huaweicloud-sdk-go-v3 package
func NewTmsClient(c *Config, region string) (*tms.TmsClient, error) {
	credentials, err := buildGlobalAuthCredentials(c, region)
	if err != nil {
		return nil, err
	}

	tmsEndpoint := getServiceEndpoint(c, "tms", region)
	if tmsEndpoint == "" {
		return nil, fmt.Errorf("failed to get the endpoint of TMS service")
	}

	return tms.NewTmsClient(
		tms.TmsClientBuilder().
			WithEndpoint(tmsEndpoint).
			WithCredential(*credentials).
			WithHttpConfig(buildHTTPConfig(c)).
			Build()), nil
}

func getProxyFromEnv() string {
	var url string

	envNames := []string{"HTTPS_PROXY", "https_proxy"}
	for _, n := range envNames {
		if val := os.Getenv(n); val != "" {
			url = val
			break
		}
	}

	return url
}

func logRequestHandler(request http.Request) {
	if !logging.IsDebugOrHigher() {
		return
	}

	log.Printf("[DEBUG] API Request URL: %s %s", request.Method, request.URL)
	log.Printf("[DEBUG] API Request Headers:\n%s", FormatHeaders(request.Header, "\n"))
	if request.Body != nil {
		if err := logRequest(request.Body, request.Header.Get("Content-Type")); err != nil {
			log.Printf("[WARN] failed to get request body: %s", err)
		}
	}
}

func logResponseHandler(response http.Response) {
	if !logging.IsDebugOrHigher() {
		return
	}

	log.Printf("[DEBUG] API Response Code: %d", response.StatusCode)
	log.Printf("[DEBUG] API Response Headers:\n%s", FormatHeaders(response.Header, "\n"))

	if err := logResponse(response.Body, response.Header.Get("Content-Type")); err != nil {
		log.Printf("[WARN] failed to get response body: %s", err)
	}
}

func logRequest(original io.ReadCloser, contentType string) error {
	defer original.Close()

	var bs bytes.Buffer
	_, err := io.Copy(&bs, original)
	if err != nil {
		return err
	}

	body := bs.Bytes()
	index := findJSONIndex(body)
	if index == -1 {
		return nil
	}

	// Handle request contentType
	if strings.HasPrefix(contentType, "application/json") {
		debugInfo := formatJSON(body[index:], true)
		log.Printf("[DEBUG] API Request Body: %s", debugInfo)
	} else {
		log.Printf("[DEBUG] Not logging because the request body isn't JSON")
	}

	return nil
}

// logResponse will log the HTTP Response details.
// If the body is JSON, it will attempt to be pretty-formatted.
func logResponse(original io.ReadCloser, contentType string) error {
	defer original.Close()

	var bs bytes.Buffer
	_, err := io.Copy(&bs, original)
	if err != nil {
		return err
	}

	body := bs.Bytes()
	index := findJSONIndex(body)
	if index == -1 {
		return nil
	}

	if strings.HasPrefix(contentType, "application/json") {
		debugInfo := formatJSON(body[index:], true)
		log.Printf("[DEBUG] API Response Body: %s", debugInfo)
	} else {
		log.Printf("[DEBUG] Not logging because the response body isn't JSON")
	}

	return nil
}

func findJSONIndex(raw []byte) int {
	var index = -1
	for i, v := range raw {
		if v == '{' {
			index = i
			break
		}
	}

	return index
}
