package config

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/auth"
	"github.com/chnsz/golangsdk/auth/core/signer"
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
				if c.SigningAlgorithm == "" || c.SigningAlgorithm == signer.HmacSHA256 {
					return auth.Sign(req, client.AKSKAuthOptions.AccessKey, client.AKSKAuthOptions.SecretKey)
				}

				sn, err := signer.GetSigner(signer.SigningAlgorithm(c.SigningAlgorithm))
				if err != nil {
					return err
				}

				_, err = sn.Sign(req, client.AKSKAuthOptions.AccessKey, client.AKSKAuthOptions.SecretKey)
				if err != nil {
					return err
				}
				return nil
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
			SigningAlgorithm: c.SigningAlgorithm,
		}

		domainAuthOptions = golangsdk.AKSKAuthOptions{
			DomainID:         c.DomainID,
			Domain:           c.DomainName,
			AgencyName:       c.AgencyName,
			AgencyDomainName: c.AgencyDomainName,
			SigningAlgorithm: c.SigningAlgorithm,
		}
	} else {
		projectAuthOptions = golangsdk.AKSKAuthOptions{
			ProjectName:      c.TenantName,
			ProjectId:        c.TenantID,
			SigningAlgorithm: c.SigningAlgorithm,
		}

		domainAuthOptions = golangsdk.AKSKAuthOptions{
			DomainID:         c.DomainID,
			Domain:           c.DomainName,
			SigningAlgorithm: c.SigningAlgorithm,
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
		"agency_urn":          agencyUrn,
		"agency_session_name": c.AssumeRoleAgency,
	}
	if c.AssumeRoleDuration != 0 {
		createAssumeOpts["duration_seconds"] = c.AssumeRoleDuration
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

	accessKey := utils.PathSearch("credentials.access_key_id", createAssumeRespBody, "").(string)
	if accessKey == "" {
		log.Printf("[DEBUG] unable to find the access key ID of the assume credential from the API response")
	}
	secretKey := utils.PathSearch("credentials.secret_access_key", createAssumeRespBody, "").(string)
	if secretKey == "" {
		log.Printf("[DEBUG] unable to find the secret access ID of the assume credential from the API response")
	}
	securityToken := utils.PathSearch("credentials.security_token", createAssumeRespBody, "").(string)
	if securityToken == "" {
		log.Printf("[DEBUG] unable to find the security token of the assume credential from the API response")
	}
	c.AccessKey, c.SecretKey, c.SecurityToken = accessKey, secretKey, securityToken

	return buildClientByAKSK(c)
}

func buildClientByAgencyChain(c *Config) error {
	var err error
	assumeRoleList := c.AssumeRoleList
	for _, role := range assumeRoleList {
		if role.RoleDomainID != "" {
			err = getTemporaryAKSKByAgencyV5(c, role)
		} else {
			err = getTemporaryAKSKByAgency(c, role)
		}
		if err != nil {
			return err
		}
	}

	return buildClientByAKSK(c)
}

func getTemporaryAKSKByAgency(c *Config, role AssumeRole) error {
	client, err := c.HcIamV3Client(c.Region)
	if err != nil {
		return fmt.Errorf("error creating Huaweicloud IAM client: %s", err)
	}

	request := &iam_model.CreateTemporaryAccessKeyByAgencyRequest{}
	domainNameAssumeRoleIdentityAssumerole := role.RoleDomain
	durationSecondsAssumeRoleIdentityAssumerole := assumeRoleDuration
	assumeRoleIdentity := &iam_model.IdentityAssumerole{
		AgencyName:      role.RoleAgency,
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
		return fmt.Errorf("error Creating temporary accesskey by agency: %s", err)
	}
	c.AccessKey, c.SecretKey, c.SecurityToken = response.Credential.Access, response.Credential.Secret, response.Credential.Securitytoken

	return nil
}

func getTemporaryAKSKByAgencyV5(c *Config, role AssumeRole) error {
	client, err := c.NewServiceClient("sts", c.Region)
	if err != nil {
		return fmt.Errorf("error creating Huaweicloud IAM V5 client: %s", err)
	}

	createAssumeHttpUrl := "v5/agencies/assume"
	createAssumePath := client.Endpoint + createAssumeHttpUrl
	agencyUrn := "iam::" + role.RoleDomainID + ":agency:" + role.RoleAgency
	createAssumeOpts := map[string]interface{}{
		"agency_urn":          agencyUrn,
		"agency_session_name": role.RoleAgency,
	}
	if role.RoleDuration != 0 {
		createAssumeOpts["duration_seconds"] = role.RoleDuration
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

	accessKey := utils.PathSearch("credentials.access_key_id", createAssumeRespBody, "").(string)
	if accessKey == "" {
		log.Printf("[DEBUG] unable to find the access key ID of the assume credential from the API response")
	}
	secretKey := utils.PathSearch("credentials.secret_access_key", createAssumeRespBody, "").(string)
	if secretKey == "" {
		log.Printf("[DEBUG] unable to find the secret access ID of the assume credential from the API response")
	}
	securityToken := utils.PathSearch("credentials.security_token", createAssumeRespBody, "").(string)
	if securityToken == "" {
		log.Printf("[DEBUG] unable to find the security token of the assume credential from the API response")
	}
	c.AccessKey, c.SecretKey, c.SecurityToken = accessKey, secretKey, securityToken

	// set project map to empty, to use the project id of another account
	c.RegionProjectIDMap = map[string]string{}

	// rebuild the client to use new AK, SK and security_token
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

	expiresAt := utils.PathSearch("credential.expires_at", parsedBody, "").(string)
	accessKey := utils.PathSearch("credential.access", parsedBody, "").(string)
	secretKey := utils.PathSearch("credential.secret", parsedBody, "").(string)
	securityToken := utils.PathSearch("credential.securitytoken", parsedBody, "").(string)

	if accessKey == "" || secretKey == "" || securityToken == "" || expiresAt == "" {
		return fmt.Errorf("Error fetching metadata authentication information")
	}
	expairesTime, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil {
		return err
	}
	c.AccessKey, c.SecretKey, c.SecurityToken, c.SecurityKeyExpiresAt = accessKey, secretKey, securityToken, expairesTime

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

// getSubjectTokenByIdp gets the subject token using ID token
func getSubjectTokenByIdp(c *Config) (string, error) {
	iamBaseURL := fmt.Sprintf("https://iam.%s.myhuaweicloud.com", c.Region)
	idTokenURL := fmt.Sprintf("%s/v3.0/OS-AUTH/id-token/tokens", iamBaseURL)

	// Create HTTP Client with logging
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}
	client := &http.Client{
		Transport: &LogRoundTripper{
			Rt:         transport,
			MaxRetries: c.MaxRetries,
		},
	}

	// Prepare request body
	body, err := json.Marshal(map[string]interface{}{
		"auth": map[string]interface{}{
			"id_token": map[string]string{
				"id": c.AssumeRoleIdToken,
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("Error marshal Id Token request body: %s", err.Error())
	}

	// Create POST request
	req, err := http.NewRequest("POST", idTokenURL, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("Error building Id Token API request: %s", err.Error())
	}

	// Set request headers
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("X-Idp-Id", c.AssumeRoleIdpID)

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error requesting Id Token API: %s", err.Error())
	}

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("Error requesting Id Token API: status code = %d", resp.StatusCode)
	}

	// Extract subject token from response header
	subjectToken := resp.Header.Get("X-Subject-Token")
	if subjectToken == "" {
		return "", fmt.Errorf("Error fetching X-Subject-Token from Idp auth response: %s", err.Error())
	}

	return subjectToken, nil
}

// getSecurityTokenByIdp gets the security token using subject token
func getSecurityTokenByIdp(c *Config, subjectToken string) error {
	iamBaseURL := fmt.Sprintf("https://iam.%s.myhuaweicloud.com", c.Region)
	securityTokenURL := fmt.Sprintf("%s/v3.0/OS-CREDENTIAL/securitytokens", iamBaseURL)

	// Create HTTP Client with logging
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}
	client := &http.Client{
		Transport: &LogRoundTripper{
			Rt:         transport,
			MaxRetries: c.MaxRetries,
		},
	}

	// Prepare request body
	body, err := json.Marshal(map[string]interface{}{
		"auth": map[string]interface{}{
			"identity": map[string]interface{}{
				"methods": []string{"token"},
				"token": map[string]interface{}{
					"id": subjectToken,
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("Error marshal Security Token request body: %s", err.Error())
	}

	// Create POST request
	req, err := http.NewRequest("POST", securityTokenURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Error building Security Token API request: %s", err.Error())
	}

	// Set request headers
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error requesting Security Token API: %s", err.Error())
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Error requesting Security Token API: status code = %d", resp.StatusCode)
	}

	// Read response body
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading security token API response: %s", err.Error())
	}

	var parsedBody interface{}
	err = json.Unmarshal(rawBody, &parsedBody)
	if err != nil {
		return fmt.Errorf("Error unmarshal security token API response: %s", err.Error())
	}

	// Extract security token
	accessKey := utils.PathSearch("credential.access", parsedBody, "").(string)
	secretKey := utils.PathSearch("credential.secret", parsedBody, "").(string)
	securityToken := utils.PathSearch("credential.securitytoken", parsedBody, "").(string)

	if accessKey == "" || secretKey == "" || securityToken == "" {
		return errors.New("Error fetching security token from API response")
	}
	c.AccessKey, c.SecretKey, c.SecurityToken = accessKey, secretKey, securityToken

	return nil
}
