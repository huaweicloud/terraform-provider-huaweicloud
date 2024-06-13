package config

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/auth"
	huaweisdk "github.com/chnsz/golangsdk/openstack"

	iam_model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/pathorcontents"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	securityKeyURL     string = "http://169.254.169.254/openstack/latest/securitykey"
	keyExpiresDuration int64  = 600
	assumeRoleDuration int32  = 12 * 60 * 60
)

// CLI Shared Config
type SharedConfig struct {
	Current  string    `json:"current"`
	Profiles []Profile `json:"profiles"`
}

type Profile struct {
	Name             string  `json:"name"`
	Mode             string  `json:"mode"`
	AccessKeyId      string  `json:"accessKeyId"`
	SecretAccessKey  string  `json:"secretAccessKey"`
	SecurityToken    string  `json:"securityToken"`
	Region           string  `json:"region"`
	ProjectId        string  `json:"projectId"`
	DomainId         string  `json:"domainId"`
	AgencyDomainId   string  `json:"agencyDomainId"`
	AgencyDomainName string  `json:"agencyDomainName"`
	AgencyName       string  `json:"agencyName"`
	SsoAuth          SsoAuth `json:"ssoAuth"`
}

type SsoAuth struct {
	StsToken StsToken `json:"stsToken"`
}

type StsToken struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	SecurityToken   string `json:"securityToken"`
}

func buildClient(c *Config) error {
	if c.Token != "" {
		return buildClientByToken(c)
	} else if c.AccessKey != "" && c.SecretKey != "" {
		return buildClientByAKSK(c)
	} else if c.Password != "" && (c.Username != "" || c.UserID != "") {
		return buildClientByPassword(c)
	}

	return buildClientByMeta(c)
}

func generateTLSConfig(c *Config) (*tls.Config, error) {
	config := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: c.Insecure,
	}

	if c.CACertFile != "" {
		caCert, _, err := pathorcontents.Read(c.CACertFile)
		if err != nil {
			return nil, fmt.Errorf("Error reading CA Cert: %s", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))
		config.RootCAs = caCertPool
	}

	if c.ClientCertFile != "" && c.ClientKeyFile != "" {
		clientCert, _, err := pathorcontents.Read(c.ClientCertFile)
		if err != nil {
			return nil, fmt.Errorf("Error reading Client Cert: %s", err)
		}
		clientKey, _, err := pathorcontents.Read(c.ClientKeyFile)
		if err != nil {
			return nil, fmt.Errorf("Error reading Client Key: %s", err)
		}

		cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil {
			return nil, err
		}

		config.Certificates = []tls.Certificate{cert}
		config.BuildNameToCertificate()
	}

	return config, nil
}

func genClient(c *Config, ao golangsdk.AuthOptionsProvider) (*golangsdk.ProviderClient, error) {
	client, err := huaweisdk.NewClient(ao.GetIdentityEndpoint())
	if err != nil {
		return nil, err
	}

	// Set UserAgent
	client.UserAgent.Prepend(providerUserAgent)
	customUserAgent := os.Getenv("HW_TF_CUSTOM_UA")
	if customUserAgent != "" {
		client.UserAgent.Prepend(customUserAgent)
	}

	config, err := generateTLSConfig(c)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: config,
	}

	client.HTTPClient = http.Client{
		Transport: &LogRoundTripper{
			Rt:         transport,
			MaxRetries: c.MaxRetries,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if client.AKSKAuthOptions.AccessKey != "" {
				err := auth.Sign(req, client.AKSKAuthOptions.AccessKey, client.AKSKAuthOptions.SecretKey)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	if c.MaxRetries > 0 {
		client.MaxBackoffRetries = uint(c.MaxRetries)
		client.RetryBackoffFunc = retryBackoffFunc
	}

	// Validate authentication normally.
	err = huaweisdk.Authenticate(client, ao)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func genClients(c *Config, projectAuthOptions, domainAuthOptions golangsdk.AuthOptionsProvider) error {
	client, err := genClient(c, projectAuthOptions)
	if err != nil {
		return err
	}
	c.HwClient = client

	client, err = genClient(c, domainAuthOptions)
	if err == nil {
		c.DomainClient = client
	}
	return err
}

func buildClientByToken(c *Config) error {
	var projectAuthOptions, domainAuthOptions golangsdk.AuthOptions

	if c.AgencyDomainName != "" && c.AgencyName != "" {
		projectAuthOptions = golangsdk.AuthOptions{
			AgencyName:       c.AgencyName,
			AgencyDomainName: c.AgencyDomainName,
			DelegatedProject: c.DelegatedProject,
		}

		domainAuthOptions = golangsdk.AuthOptions{
			AgencyName:       c.AgencyName,
			AgencyDomainName: c.AgencyDomainName,
		}
	} else {
		projectAuthOptions = golangsdk.AuthOptions{
			DomainID:   c.DomainID,
			DomainName: c.DomainName,
			TenantID:   c.TenantID,
			TenantName: c.TenantName,
		}

		domainAuthOptions = golangsdk.AuthOptions{
			DomainID:   c.DomainID,
			DomainName: c.DomainName,
		}
	}

	for _, ao := range []*golangsdk.AuthOptions{&projectAuthOptions, &domainAuthOptions} {
		ao.IdentityEndpoint = c.IdentityEndpoint
		ao.TokenID = c.Token

	}
	return genClients(c, projectAuthOptions, domainAuthOptions)
}

func buildClientByAKSK(c *Config) error {
	var projectAuthOptions, domainAuthOptions golangsdk.AKSKAuthOptions

	if c.AgencyDomainName != "" && c.AgencyName != "" {
		projectAuthOptions = golangsdk.AKSKAuthOptions{
			DomainID:         c.DomainID,
			Domain:           c.DomainName,
			AgencyName:       c.AgencyName,
			AgencyDomainName: c.AgencyDomainName,
			DelegatedProject: c.DelegatedProject,
		}

		domainAuthOptions = golangsdk.AKSKAuthOptions{
			DomainID:         c.DomainID,
			Domain:           c.DomainName,
			AgencyName:       c.AgencyName,
			AgencyDomainName: c.AgencyDomainName,
		}
	} else {
		projectAuthOptions = golangsdk.AKSKAuthOptions{
			ProjectName: c.TenantName,
			ProjectId:   c.TenantID,
		}

		domainAuthOptions = golangsdk.AKSKAuthOptions{
			DomainID: c.DomainID,
			Domain:   c.DomainName,
		}
	}

	for _, ao := range []*golangsdk.AKSKAuthOptions{&projectAuthOptions, &domainAuthOptions} {
		ao.IdentityEndpoint = c.IdentityEndpoint
		ao.AccessKey = c.AccessKey
		ao.SecretKey = c.SecretKey
		ao.WithUserCatalog = true

		if c.Region != "" {
			ao.Region = c.Region
		}
		if c.SecurityToken != "" {
			ao.SecurityToken = c.SecurityToken
		}
	}
	return genClients(c, projectAuthOptions, domainAuthOptions)
}

func buildClientByPassword(c *Config) error {
	var projectAuthOptions, domainAuthOptions golangsdk.AuthOptions

	if c.AgencyDomainName != "" && c.AgencyName != "" {
		projectAuthOptions = golangsdk.AuthOptions{
			DomainID:         c.DomainID,
			DomainName:       c.DomainName,
			AgencyName:       c.AgencyName,
			AgencyDomainName: c.AgencyDomainName,
			DelegatedProject: c.DelegatedProject,
		}

		domainAuthOptions = golangsdk.AuthOptions{
			DomainID:         c.DomainID,
			DomainName:       c.DomainName,
			AgencyName:       c.AgencyName,
			AgencyDomainName: c.AgencyDomainName,
		}
	} else {
		projectAuthOptions = golangsdk.AuthOptions{
			DomainID:   c.DomainID,
			DomainName: c.DomainName,
			TenantID:   c.TenantID,
			TenantName: c.TenantName,
		}

		domainAuthOptions = golangsdk.AuthOptions{
			DomainID:   c.DomainID,
			DomainName: c.DomainName,
		}
	}

	for _, ao := range []*golangsdk.AuthOptions{&projectAuthOptions, &domainAuthOptions} {
		ao.IdentityEndpoint = c.IdentityEndpoint
		ao.Password = c.Password
		ao.Username = c.Username
		ao.UserID = c.UserID
	}
	return genClients(c, projectAuthOptions, domainAuthOptions)
}

func buildClientByAgency(c *Config) error {
	client, err := c.HcIamV3Client(c.Region)
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud IAM client: %s", err)
	}

	request := &iam_model.CreateTemporaryAccessKeyByAgencyRequest{}
	domainNameAssumeRoleIdentityAssumerole := c.AssumeRoleDomain
	durationSecondsAssumeRoleIdentityAssumerole := assumeRoleDuration
	assumeRoleIdentity := &iam_model.IdentityAssumerole{
		AgencyName:      c.AssumeRoleAgency,
		DomainName:      &domainNameAssumeRoleIdentityAssumerole,
		DurationSeconds: &durationSecondsAssumeRoleIdentityAssumerole,
	}
	var listMethodsIdentity = []iam_model.AgencyAuthIdentityMethods{
		iam_model.GetAgencyAuthIdentityMethodsEnum().ASSUME_ROLE,
	}
	identityAuth := &iam_model.AgencyAuthIdentity{
		Methods:    listMethodsIdentity,
		AssumeRole: assumeRoleIdentity,
	}
	authbody := &iam_model.AgencyAuth{
		Identity: identityAuth,
	}
	request.Body = &iam_model.CreateTemporaryAccessKeyByAgencyRequestBody{
		Auth: authbody,
	}
	response, err := client.CreateTemporaryAccessKeyByAgency(request)
	if err != nil {
		return fmt.Errorf("Error Creating temporary accesskey by agency: %s", err)
	}
	c.AccessKey, c.SecretKey, c.SecurityToken = response.Credential.Access, response.Credential.Secret, response.Credential.Securitytoken

	return buildClientByAKSK(c)
}

func buildClientByAgencyV5(c *Config) error {
	client, err := c.NewServiceClient("sts", c.Region)
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud IAM V5 client: %s", err)
	}

	createAssumeHttpUrl := "v5/agencies/assume"
	createAssumePath := client.Endpoint + createAssumeHttpUrl
	agencyUrn := "iam::" + c.AssumeRoleDomainID + ":agency:" + c.AssumeRoleAgency
	createAssumeOpts := map[string]interface{}{
		"duration_seconds":    assumeRoleDuration,
		"agency_urn":          agencyUrn,
		"agency_session_name": c.AssumeRoleAgency,
	}
	createAssumeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(createAssumeOpts),
	}
	createAssumeResp, err := client.Request("POST", createAssumePath, &createAssumeOpt)
	if err != nil {
		return fmt.Errorf("error creating IAM agency assume: %s", err)
	}
	createAssumeRespBody, err := utils.FlattenResponse(createAssumeResp)
	if err != nil {
		return fmt.Errorf("error extracting IAM agency assume response: %s", err)
	}

	accessKey, err := jmespath.Search("credentials.access_key_id", createAssumeRespBody)
	if err != nil {
		return fmt.Errorf("error fetching assume credentials: access_key_id is not found in API response")
	}
	secretKey, err := jmespath.Search("credentials.secret_access_key", createAssumeRespBody)
	if err != nil {
		return fmt.Errorf("error fetching assume credentials: secret_access_id is not found in API response")
	}
	securityToken, err := jmespath.Search("credentials.security_token", createAssumeRespBody)
	if err != nil {
		return fmt.Errorf("error fetching assume credentials: security_token is not found in API response")
	}
	c.AccessKey, c.SecretKey, c.SecurityToken = accessKey.(string), secretKey.(string), securityToken.(string)

	return buildClientByAKSK(c)
}

func (c *Config) reloadSecurityKey() error {
	err := getAuthConfigByMeta(c)
	if err != nil {
		return fmt.Errorf("Error reloading Auth credentials from ECS Metadata API: %s", err)
	}
	log.Printf("Successfully reload metadata security key, which will expire at: %s", c.SecurityKeyExpiresAt)
	return buildClientByAKSK(c)
}

func getAuthConfigByMeta(c *Config) error {
	req, err := http.NewRequest("GET", securityKeyURL, nil)
	if err != nil {
		return fmt.Errorf("Error building metadata API request: %s", err.Error())
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error requesting metadata API: %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error requesting metadata API: status code = %d", resp.StatusCode)
	}

	var parsedBody interface{}

	defer resp.Body.Close()
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error parsing metadata API response: %s", err.Error())
	}

	err = json.Unmarshal(rawBody, &parsedBody)
	if err != nil {
		return fmt.Errorf("Error unmarshal metadata API, agency_name is empty: %s", err.Error())
	}

	expiresAt, err := jmespath.Search("credential.expires_at", parsedBody)
	if err != nil {
		return fmt.Errorf("Error fetching metadata expires_at: %s", err.Error())
	}
	accessKey, err := jmespath.Search("credential.access", parsedBody)
	if err != nil {
		return fmt.Errorf("Error fetching metadata access: %s", err.Error())
	}
	secretKey, err := jmespath.Search("credential.secret", parsedBody)
	if err != nil {
		return fmt.Errorf("Error fetching metadata secret: %s", err.Error())
	}
	securityToken, err := jmespath.Search("credential.securitytoken", parsedBody)
	if err != nil {
		return fmt.Errorf("Error fetching metadata securitytoken: %s", err.Error())
	}

	if accessKey == nil || secretKey == nil || securityToken == nil || expiresAt == nil {
		return fmt.Errorf("Error fetching metadata authentication information")
	}
	expairesTime, err := time.Parse(time.RFC3339, expiresAt.(string))
	if err != nil {
		return err
	}
	c.AccessKey, c.SecretKey, c.SecurityToken, c.SecurityKeyExpiresAt = accessKey.(string), secretKey.(string), securityToken.(string), expairesTime

	return nil
}

func buildClientByMeta(c *Config) error {
	err := getAuthConfigByMeta(c)
	if err != nil {
		return fmt.Errorf("Error fetching Auth credentials from ECS Metadata API, AkSk or ECS agency must be provided: %s", err)
	}
	log.Printf("[DEBUG] Successfully got metadata security key, which will expire at: %s", c.SecurityKeyExpiresAt)
	return buildClientByAKSK(c)
}
