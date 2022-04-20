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
	hcconfig "github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/httphandler"
	cptsv1 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cpts/v1"
	ctsv3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cts/v3"
	gaussdbv3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/gaussdb/v3"
	iamv3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	kpsv3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/kps/v3"
	rdsv3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3"
	tmsv1 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/tms/v1"
	vpcv3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
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
		// S1005
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

func buildHTTPConfig(c *Config) *hcconfig.HttpConfig {
	httpConfig := hcconfig.DefaultHttpConfig()

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

			httpProxy := hcconfig.Proxy{
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

// try to get the endpoint from customizing map
func getServiceEndpoint(c *Config, srv, region string) string {
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

// HcVpcV3Client is the VPC service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcVpcV3Client(region string) (*vpcv3.VpcClient, error) {
	hcClient, err := NewHcClient(c, region, "vpc", false)
	if err != nil {
		return nil, err
	}

	return vpcv3.NewVpcClient(hcClient), nil
}

// HcTmsV1Client is the TMS service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcTmsV1Client(region string) (*tmsv1.TmsClient, error) {
	hcClient, err := NewHcClient(c, region, "tms", true)
	if err != nil {
		return nil, err
	}
	return tmsv1.NewTmsClient(hcClient), nil
}

// HcKmsV3Client is the KMS service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcKmsV3Client(region string) (*kpsv3.KpsClient, error) {
	hcClient, err := NewHcClient(c, region, "kms", false)
	if err != nil {
		return nil, err
	}
	return kpsv3.NewKpsClient(hcClient), nil
}

// HcIamV3Client is the IAM service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcIamV3Client(region string) (*iamv3.IamClient, error) {
	hcClient, err := NewHcClient(c, region, "iam", true)
	if err != nil {
		return nil, err
	}
	return iamv3.NewIamClient(hcClient), nil
}

// HcGaussdbV3Client is the Gaussdb service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcGaussdbV3Client(region string) (*gaussdbv3.GaussDBClient, error) {
	hcClient, err := NewHcClient(c, region, "gaussdb", false)
	if err != nil {
		return nil, err
	}
	return gaussdbv3.NewGaussDBClient(hcClient), nil
}

// HcCtsV3Client is the CTS service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcCtsV3Client(region string) (*ctsv3.CtsClient, error) {
	hcClient, err := NewHcClient(c, region, "cts", false)
	if err != nil {
		return nil, err
	}
	return ctsv3.NewCtsClient(hcClient), nil
}

// HcRdsV3Client is the RDS service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcRdsV3Client(region string) (*rdsv3.RdsClient, error) {
	hcClient, err := NewHcClient(c, region, "rds", false)
	if err != nil {
		return nil, err
	}
	return rdsv3.NewRdsClient(hcClient), nil
}

// HcCptsV3Client is the CPTS service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcCptsV3Client(region string) (*cptsv1.CptsClient, error) {
	hcClient, err := NewHcClient(c, region, "cpts", false)
	if err != nil {
		return nil, err
	}
	return cptsv1.NewCptsClient(hcClient), nil
}

// NewHcClient is the common client using huaweicloud-sdk-go-v3 package
func NewHcClient(c *Config, region, product string, globalFlag bool) (*core.HcHttpClient, error) {
	endpoint := getServiceEndpoint(c, product, region)
	if endpoint == "" {
		return nil, fmt.Errorf("failed to get the endpoint of %q service in region %s", product, region)
	}

	builder := core.NewHcHttpClientBuilder().WithEndpoint(endpoint).WithHttpConfig(buildHTTPConfig(c))

	if globalFlag {
		credentials, err := buildGlobalAuthCredentials(c, region)
		if err != nil {
			return nil, err
		}
		builder.WithCredentialsType("global.Credentials").WithCredential(*credentials)
	} else {
		credentials, err := buildAuthCredentials(c, region)
		if err != nil {
			return nil, err
		}
		builder.WithCredential(*credentials)
	}

	return builder.Build(), nil
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
