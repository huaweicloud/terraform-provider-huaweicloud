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
	aomv2 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/aom/v2"
	ccev3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cce/v3"
	cdnv1 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v1"
	cdnv2 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2"
	cptsv1 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cpts/v1"
	cssv1 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/css/v1"
	cssv2 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/css/v2"
	ctsv3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cts/v3"
	hssv5 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/hss/v5"
	iamv3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	iotdav5 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5"
	dmsv2 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/kafka/v2"
	kpsv3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/kps/v3"
	livev1 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/live/v1"
	mpcv1 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/mpc/v1"
	omsv2 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/oms/v2"
	rdsv3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3"
	tmsv1 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/tms/v1"
	vodv1 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vod/v1"
	vpcv3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
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
			logp.Printf("[DEBUG] using https proxy: %s://%s", proxyURL.Scheme, proxyURL.Host)

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
		logp.Printf("[WARN] parsing https proxy failed: %s", err)
	}

	return httpConfig
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

// HcCptsV1Client is the CPTS service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcCptsV1Client(region string) (*cptsv1.CptsClient, error) {
	hcClient, err := NewHcClient(c, region, "cpts", false)
	if err != nil {
		return nil, err
	}
	return cptsv1.NewCptsClient(hcClient), nil
}

// HcVodV1Client is the AOM service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcVodV1Client(region string) (*vodv1.VodClient, error) {
	hcClient, err := NewHcClient(c, region, "vod", false)
	if err != nil {
		return nil, err
	}
	return vodv1.NewVodClient(hcClient), nil
}

// HcAomV2Client is the AOM service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcAomV2Client(region string) (*aomv2.AomClient, error) {
	hcClient, err := NewHcClient(c, region, "aom", false)
	if err != nil {
		return nil, err
	}
	return aomv2.NewAomClient(hcClient), nil
}

// HcLiveV1Client is the live service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcLiveV1Client(region string) (*livev1.LiveClient, error) {
	hcClient, err := NewHcClient(c, region, "live", false)
	if err != nil {
		return nil, err
	}
	return livev1.NewLiveClient(hcClient), nil
}

// HcMpcV1Client is the MPC service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcMpcV1Client(region string) (*mpcv1.MpcClient, error) {
	hcClient, err := NewHcClient(c, region, "mpc", false)
	if err != nil {
		return nil, err
	}
	return mpcv1.NewMpcClient(hcClient), nil
}

// HcIoTdaV5Client is the IoTDA service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcIoTdaV5Client(region string, isDerived bool) (*iotdav5.IoTDAClient, error) {
	hcClient, err := implNewHcClient(c, region, "iotda", false, isDerived)
	if err != nil {
		return nil, err
	}
	return iotdav5.NewIoTDAClient(hcClient), nil
}

// HcMpcV1Client is the MPC service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcOmsV2Client(region string) (*omsv2.OmsClient, error) {
	hcClient, err := NewHcClient(c, region, "oms", false)
	if err != nil {
		return nil, err
	}
	return omsv2.NewOmsClient(hcClient), nil
}

// HcCssV1Client is the css service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcCssV1Client(region string) (*cssv1.CssClient, error) {
	hcClient, err := NewHcClient(c, region, "css", false)
	if err != nil {
		return nil, err
	}
	return cssv1.NewCssClient(hcClient), nil
}

// HcCssV2Client is the css service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcCssV2Client(region string) (*cssv2.CssClient, error) {
	hcClient, err := NewHcClient(c, region, "css", false)
	if err != nil {
		return nil, err
	}
	return cssv2.NewCssClient(hcClient), nil
}

// HcCdnV1Client is the CDN service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcCdnV1Client(region string) (*cdnv1.CdnClient, error) {
	hcClient, err := NewHcClient(c, region, "cdn", false)
	if err != nil {
		return nil, err
	}
	return cdnv1.NewCdnClient(hcClient), nil
}

// HcCdnV2Client is the CDN service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcCdnV2Client(region string) (*cdnv2.CdnClient, error) {
	hcClient, err := NewHcClient(c, region, "cdn", false)
	if err != nil {
		return nil, err
	}
	return cdnv2.NewCdnClient(hcClient), nil
}

// HcDmsV2Client is the DMS service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcDmsV2Client(region string) (*dmsv2.KafkaClient, error) {
	hcClient, err := NewHcClient(c, region, "dmsv2", false)
	if err != nil {
		return nil, err
	}
	return dmsv2.NewKafkaClient(hcClient), nil
}

// HcHssV5Client is the HSS service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcHssV5Client(region string) (*hssv5.HssClient, error) {
	hcClient, err := NewHcClient(c, region, "hss", false)
	if err != nil {
		return nil, err
	}
	return hssv5.NewHssClient(hcClient), nil
}

// HcCceV3Client is the CCE service client using huaweicloud-sdk-go-v3 package
func (c *Config) HcCceV3Client(region string) (*ccev3.CceClient, error) {
	hcClient, err := NewHcClient(c, region, "cce", false)
	if err != nil {
		return nil, err
	}
	return ccev3.NewCceClient(hcClient), nil
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
		credentials, err := buildGlobalAuthCredentials(c, region)
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
