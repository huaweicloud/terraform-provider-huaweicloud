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
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/signer/algorithm"
	hcconfig "github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/httphandler"
	hcregion "github.com/huaweicloud/huaweicloud-sdk-go-v3/core/region"
	iamv3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
)

/*
This file is used to impl the configuration of huaweicloud-sdk-go-v3 package and
genetate service clients.
*/
func buildAuthCredentials(c *Config, region string, isDerived bool) (*basic.Credentials, error) {
	if c.AccessKey == "" || c.SecretKey == "" {
		return nil, fmt.Errorf("access_key or secret_key is missing in the provider")
	}

	credentials := basic.Credentials{
		AK:            c.AccessKey,
		SK:            c.SecretKey,
		SecurityToken: c.SecurityToken,
		IamEndpoint:   c.IdentityEndpoint,
	}

	if isDerived {
		credentials.DerivedPredicate = basic.DefaultDerivedPredicate
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
		projectID = c.RegionProjectIDMap[region]
	}

	credentials.ProjectId = projectID
	return &credentials, nil
}

func buildGlobalAuthCredentials(c *Config) (*global.Credentials, error) {
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
	if c.SigningAlgorithm != "" {
		httpConfig.WithSigningAlgorithm(algorithm.SigningAlgorithm(c.SigningAlgorithm))
	}

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

	if proxyURL, err := parseProxyFromEnv(); err == nil {
		if proxyURL != nil {
			log.Printf("[DEBUG] using https proxy: %s://%s", proxyURL.Scheme, proxyURL.Host)

			httpProxy := hcconfig.Proxy{
				Schema:   proxyURL.Scheme,
				Host:     proxyURL.Host,
				Username: proxyURL.User.Username(),
			}
			if pwd, ok := proxyURL.User.Password(); ok {
				httpProxy.Password = pwd
			}

			httpConfig = httpConfig.WithProxy(&httpProxy)
		}
	} else {
		log.Printf("[WARN] parsing https proxy failed: %s", err)
	}

	return httpConfig
}

// HcIamV3Client is the IAM service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcIamV3Client(region string) (*iamv3.IamClient, error) {
	hcClient, err := NewHcClient(c, region, "iam", true)
	if err != nil {
		return nil, err
	}
	return iamv3.NewIamClient(hcClient), nil
}

// NewHcClient is the common client using huaweicloud-sdk-go-v3 package
func NewHcClient(c *Config, region, product string, isGlobal bool) (*core.HcHttpClient, error) {
	return implNewHcClient(c, region, product, isGlobal, false)
}

func implNewHcClient(c *Config, region, product string, isGlobal, isDerived bool) (*core.HcHttpClient, error) {
	endpoint := GetServiceEndpoint(c, product, region)
	if endpoint == "" {
		return nil, fmt.Errorf("failed to get the endpoint of %q service in region %s", product, region)
	}

	builder := core.NewHcHttpClientBuilder().
		WithRegion(hcregion.NewRegion(region, endpoint)).
		WithHttpConfig(buildHTTPConfig(c))

	if isGlobal {
		credentials, err := buildGlobalAuthCredentials(c)
		if err != nil {
			return nil, err
		}
		builder.WithCredentialsType("global.Credentials").WithCredential(credentials)
	} else {
		credentials, err := buildAuthCredentials(c, region, isDerived)
		if err != nil {
			return nil, err
		}

		builder.WithCredential(credentials)
		if isDerived {
			// the derivedAuthServiceName is fixed to "iotdm", now only IoTDA service need derived sign
			builder.WithDerivedAuthServiceName("iotdm")
		}
	}

	headers := make(map[string]string)
	customUserAgent := os.Getenv("HW_TF_CUSTOM_UA")
	if customUserAgent != "" {
		headers["User-Agent"] = fmt.Sprintf("%s;%s", providerUserAgent, customUserAgent)
	} else {
		headers["User-Agent"] = providerUserAgent
	}

	return builder.Build().PreInvoke(headers), nil
}

func parseProxyFromEnv() (*url.URL, error) {
	var proxy string

	envNames := []string{"HTTPS_PROXY", "https_proxy"}
	for _, n := range envNames {
		if val := os.Getenv(n); val != "" {
			proxy = val
			break
		}
	}

	if proxy == "" {
		return nil, nil
	}

	proxyURL, err := url.Parse(proxy)
	if err != nil ||
		(proxyURL.Scheme != "http" &&
			proxyURL.Scheme != "https" &&
			proxyURL.Scheme != "socks5") {
		// proxy was bogus. Try prepending "http://" to it and
		// see if that parses correctly. If not, we fall
		// through and complain about the original one.
		if proxyURL, err := url.Parse("http://" + proxy); err == nil {
			return proxyURL, nil
		}
	}
	if err != nil {
		return nil, fmt.Errorf("invalid https proxy address %q: %v", proxy, err)
	}
	return proxyURL, nil
}

func logRequestHandler(request http.Request) {
	requestAt := fmt.Sprintf("%d-0", time.Now().UnixMilli())
	log.Printf("[DEBUG] [%s] API Request URL: %s %s\nAPI Request Headers:\n%s",
		requestAt, request.Method, request.URL, FormatHeaders(request.Header, "\n"))

	if request.Body != nil {
		if err := logRequest(request.Body, request.Header.Get("Content-Type"), requestAt); err != nil {
			log.Printf("[WARN] [%s] failed to log API Request Body: %s", requestAt, err)
		}
	}
}

func logResponseHandler(response http.Response) {
	responseAt := fmt.Sprintf("%d-0", time.Now().UnixMilli())
	log.Printf("[DEBUG] [%s] API Response Code: %d\nAPI Response Headers:\n%s",
		responseAt, response.StatusCode, FormatHeaders(response.Header, "\n"))

	if response.Body != nil {
		if err := logResponse(response.Body, response.Header.Get("Content-Type"), responseAt); err != nil {
			log.Printf("[WARN] [%s] failed to log API Response Body: %s", responseAt, err)
		}
	}
}

// logRequest will log the HTTP Request details, then close the original.
// If the body is JSON, it will attempt to be pretty-formatted.
func logRequest(original io.ReadCloser, contentType, requestAt string) error {
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
		debugInfo := formatJSON(body[index:], requestAt, true)
		log.Printf("[DEBUG] [%s] API Request Body: %s", requestAt, debugInfo)
	} else {
		log.Printf("[DEBUG] [%s] Not logging because the request body isn't JSON", requestAt)
	}

	return nil
}

// logResponse will log the HTTP Response details, then close the original.
// If the body is JSON, it will attempt to be pretty-formatted.
func logResponse(original io.ReadCloser, contentType, responseAt string) error {
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
		debugInfo := formatJSON(body[index:], responseAt, true)
		log.Printf("[DEBUG] [%s] API Response Body: %s", responseAt, debugInfo)
	} else {
		log.Printf("[DEBUG] [%s] Not logging because the response body isn't JSON", responseAt)
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
