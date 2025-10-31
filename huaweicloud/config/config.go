package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3/domains"
	"github.com/chnsz/golangsdk/openstack/identity/v3/projects"
	"github.com/chnsz/golangsdk/openstack/identity/v3/users"
	"github.com/chnsz/golangsdk/openstack/obs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/mutexkv"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	providerUserAgent string = "terraform-provider-iac"
	InternationalSite string = "International"
)

var (
	// MutexKV is a global lock on all resources, it can lock the specified shared string (such as resource ID, resource
	// Name, port, etc.) to prevent other resources from using it, for concurrency control.
	// Usage: MutexKV.Lock({resource ID}) and MutexKV.Unlock({resource ID})
	MutexKV = mutexkv.NewMutexKV()
	// If an account sends a CBC request and it crosses websites, the CBC service will return a 403 error to indicate
	// attention.
	crossWebsiteErrs = []string{
		"CBC.0150",
		"CBC.0156",
	}
)

type Config struct {
	AccessKey             string
	SecretKey             string
	CACertFile            string
	ClientCertFile        string
	ClientKeyFile         string
	DomainID              string
	DomainName            string
	IdentityEndpoint      string
	Insecure              bool
	Region                string
	TenantID              string
	TenantName            string
	Token                 string
	SecurityToken         string
	AssumeRoleAgency      string
	AssumeRoleDomain      string
	AssumeRoleDomainID    string
	AssumeRoleDuration    int
	AssumeRoleList        []AssumeRole
	AssumeRoleIdpID       string
	AssumeRoleIdToken     string
	AssumeRoleIdTokenFile string
	Cloud                 string
	MaxRetries            int
	TerraformVersion      string
	RegionClient          bool
	EnterpriseProjectID   string
	SharedConfigFile      string
	Profile               string

	// metadata security key expires at
	SecurityKeyExpiresAt time.Time

	HwClient     *golangsdk.ProviderClient
	DomainClient *golangsdk.ProviderClient

	// websiteType is the site type of HuaweiCloud.
	// The value can be Chinese(default) and International.
	websiteType string
	// the custom endpoints used to override the default endpoint URL
	Endpoints map[string]string

	// RegionProjectIDMap is a map which stores the region-projectId pairs,
	// and region name will be the key and projectID will be the value in this map.
	RegionProjectIDMap map[string]string

	// RPLock is used to make the accessing of RegionProjectIDMap serial,
	// prevent sending duplicate query requests
	RPLock *sync.Mutex

	// SecurityKeyLock is used to make the accessing of SecurityKeyExpiresAt serial,
	// prevent sending duplicate query metadata api
	SecurityKeyLock *sync.Mutex

	// Legacy
	Username         string
	UserID           string
	Password         string
	AgencyName       string
	AgencyDomainName string
	DelegatedProject string

	// Metadata is used for extend
	Metadata any

	EnableForceNew bool

	// SigningAlgorithm is used to select encryption algorithm
	SigningAlgorithm string
	DefaultTags      map[string]interface{}
}

type AssumeRole struct {
	RoleAgency   string
	RoleDomain   string
	RoleDomainID string
	RoleDuration int
}

func (c *Config) LoadAndValidate() error {
	if c.MaxRetries < 0 {
		return fmt.Errorf("max_retries should be a positive value")
	}

	// Assume role OIDC
	if c.AssumeRoleIdpID != "" {
		if c.Region == "" {
			return errors.New("region should be provided")
		}

		// Get token from file if not specified
		if c.AssumeRoleIdToken == "" && c.AssumeRoleIdTokenFile != "" {
			token, err := os.ReadFile(c.AssumeRoleIdTokenFile)
			if err != nil {
				return fmt.Errorf("Error reading id_token_file: %s", err)
			}
			tokenStr := string(token)
			c.AssumeRoleIdToken = strings.Trim(tokenStr, "\n")
		}

		subjectToken, err := getSubjectTokenByIdp(c)
		if err != nil {
			return err
		}
		err = getSecurityTokenByIdp(c, subjectToken)
		if err != nil {
			return err
		}
	}

	err := buildClient(c)
	if err != nil {
		return err
	}

	if c.Region == "" {
		return fmt.Errorf("region should be provided")
	}

	// Assume role
	if c.AssumeRoleAgency != "" {
		if c.AssumeRoleDomainID != "" {
			err = buildClientByAgencyV5(c)
		} else {
			err = buildClientByAgency(c)
		}
		if err != nil {
			return err
		}
	}

	// Assume role list
	if len(c.AssumeRoleList) != 0 {
		return buildClientByAgencyChain(c)
	}

	if c.HwClient != nil && c.HwClient.ProjectID != "" {
		c.RegionProjectIDMap[c.Region] = c.HwClient.ProjectID
	}
	log.Printf("[DEBUG] init region and project map: %#v", c.RegionProjectIDMap)

	// set DomainID for IAM resource
	if c.DomainID == "" {
		if domainID, err := c.getDomainID(); err == nil {
			c.DomainID = domainID

			// update DomainClient.AKSKAuthOptions
			if c.DomainClient.AKSKAuthOptions.AccessKey != "" {
				c.DomainClient.AKSKAuthOptions.DomainID = c.DomainID
			}
		} else {
			log.Printf("[WARN] get domain id failed: %s", err)
		}
	}

	if c.UserID == "" && c.Username != "" {
		if userID, err := c.getUserIDbyName(c.Username); err == nil {
			c.UserID = userID
		} else {
			log.Printf("[WARN] get user id failed: %s", err)
		}
	}

	if c.SigningAlgorithm != "" {
		c.HwClient.AKSKAuthOptions.SigningAlgorithm = c.SigningAlgorithm
		c.DomainClient.AKSKAuthOptions.SigningAlgorithm = c.SigningAlgorithm
	}

	return nil
}

// SetWebsiteType will update WebsiteType field by a probe API.
// we will get status code 403 and the following response body in International website with https://bss.myhuaweicloud.com
//
//	{
//	  "error_code": "CBC.0150",
//	  "error_msg": "Access denied. The customer does not belong to the website you are now at."
//	}
//
// In addition to the error code 'CBC.0150', some regions also return error code 'CBC.0156'.
// we can call the probe API and parse the response body to decide whether the account belongs to International website or not.
// we select https://support.huaweicloud.com/intl/zh-cn/api-oce/zh-cn_topic_0000001256679455.html as the probe API.
func (c *Config) SetWebsiteType() error {
	bssClient, err := c.NewServiceClient("bss", c.Region)
	if err != nil {
		return fmt.Errorf("error creating BSS client: %s", err)
	}

	probeUrlPath := bssClient.Endpoint + "v2/products/service-types?limit=1"
	probeRequestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = bssClient.Request("GET", probeUrlPath, &probeRequestOpt)
	if err != nil {
		if respErr, ok := err.(golangsdk.ErrDefault403); ok {
			resp := struct {
				ErrorCode string `json:"error_code"`
				ErrorMsg  string `json:"error_msg"`
			}{}

			if decodeErr := json.Unmarshal(respErr.Body, &resp); decodeErr != nil {
				log.Printf("[WARN] failed to unmarshal the response body: %s", decodeErr)
			}

			if utils.IsStrContainsSliceElement(resp.ErrorCode, crossWebsiteErrs, false, true) {
				log.Printf("[DEBUG] the current account belongs to %s website", InternationalSite)
				c.websiteType = InternationalSite
				return nil
			}
		}
		return err
	}

	return nil
}

func (c *Config) GetWebsiteType() string {
	return c.websiteType
}

func (c *Config) SetServiceEndpoint(service, endpoint string) {
	// only update the customizing service endpoint when it isn't specified
	if _, ok := c.Endpoints[service]; ok {
		return
	}

	c.Endpoints[service] = endpoint
	multiKeys := GetServiceDerivedCatalogKeys(service)
	for _, k := range multiKeys {
		c.Endpoints[k] = endpoint
	}
}

func retryBackoffFunc(ctx context.Context, respErr *golangsdk.ErrUnexpectedResponseCode, e error, retries uint) error {
	minutes := int(math.Pow(2, float64(retries)))
	if minutes > 30 { // won't wait more than 30 minutes
		minutes = 30
	}

	log.Printf("[WARN] Received StatusTooManyRequests response code, try to sleep %d minutes", minutes)
	sleep := time.Duration(minutes) * time.Minute

	if ctx != nil {
		select {
		case <-time.After(sleep):
		case <-ctx.Done():
			return e
		}
	} else {
		//lintignore:R018
		time.Sleep(sleep)
	}

	return nil
}

func getObsEndpoint(c *Config, region string) string {
	if endpoint, ok := c.Endpoints["obs"]; ok {
		// replace the region in customizing OBS endpoint
		subparts := strings.Split(endpoint, ".")
		if len(subparts) >= 3 && subparts[1] != region {
			subparts[1] = region
			return strings.Join(subparts, ".")
		}
		return endpoint
	}
	return fmt.Sprintf("https://obs.%s.%s/", region, c.Cloud)
}

func (c *Config) ObjectStorageClientWithSignature(region string) (*obs.ObsClient, error) {
	if c.AccessKey == "" || c.SecretKey == "" {
		return nil, fmt.Errorf("missing credentials for OBS, need access_key and secret_key values for provider")
	}

	clientConfigure := obs.WithHttpClient(&c.DomainClient.HTTPClient)
	userAgentConfigure := obs.WithUserAgent(buildObsUserAgent())
	envProxyConfigure := obs.WithProxyFromEnv(true)
	obsEndpoint := getObsEndpoint(c, region)
	if c.SecurityToken != "" {
		return obs.New(c.AccessKey, c.SecretKey, obsEndpoint,
			obs.WithSignature("OBS"), obs.WithSecurityToken(c.SecurityToken), clientConfigure,
			userAgentConfigure, envProxyConfigure)
	}
	return obs.New(c.AccessKey, c.SecretKey, obsEndpoint, obs.WithSignature("OBS"), clientConfigure,
		userAgentConfigure, envProxyConfigure)
}

func (c *Config) ObjectStorageClient(region string) (*obs.ObsClient, error) {
	if c.AccessKey == "" || c.SecretKey == "" {
		return nil, fmt.Errorf("missing credentials for OBS, need access_key and secret_key values for provider")
	}

	if !c.SecurityKeyExpiresAt.IsZero() {
		c.SecurityKeyLock.Lock()
		defer c.SecurityKeyLock.Unlock()
		timeNow := time.Now().Unix()
		expairesAtInt := c.SecurityKeyExpiresAt.Unix()
		if timeNow+keyExpiresDuration > expairesAtInt {
			err := c.reloadSecurityKey()
			if err != nil {
				return nil, err
			}
		}
	}

	clientConfigure := obs.WithHttpClient(&c.DomainClient.HTTPClient)
	userAgentConfigure := obs.WithUserAgent(buildObsUserAgent())
	envProxyConfigure := obs.WithProxyFromEnv(true)
	obsEndpoint := getObsEndpoint(c, region)
	if c.SecurityToken != "" {
		return obs.New(c.AccessKey, c.SecretKey, obsEndpoint, obs.WithSecurityToken(c.SecurityToken), clientConfigure,
			userAgentConfigure, envProxyConfigure)
	}
	return obs.New(c.AccessKey, c.SecretKey, obsEndpoint, clientConfigure, userAgentConfigure, envProxyConfigure)
}

func buildObsUserAgent() string {
	var agent string = providerUserAgent
	if customUserAgent := os.Getenv("HW_TF_CUSTOM_UA"); customUserAgent != "" {
		agent = fmt.Sprintf("%s %s", customUserAgent, providerUserAgent)
	}

	return agent
}

// When performing derivative calculations, the derivative service name must exist.
// Currently only iotda service supports derivative algorithms.
// When a new service requires derivative algorithm calculation, please maintain this configuration.
var derivedAuthServiceNameMap = map[string]string{
	"iotda": "iotdm",
}

// NewServiceClientWithDerivedAuth create a ServiceClient that performs derivative algorithm authentication.
// Please confirm that the `derivedAuthServiceNameMap` configuration contains the derived service name of the service.
func (c *Config) NewServiceClientWithDerivedAuth(srv, region string, isDerived bool) (*golangsdk.ServiceClient, error) {
	client, err := c.NewServiceClient(srv, region)
	if err != nil {
		return nil, err
	}

	// Used to enable derivative calculations of AK/SK.
	if isDerived {
		client.AKSKAuthOptions.IsDerived = true
		client.AKSKAuthOptions.DerivedAuthServiceName = derivedAuthServiceNameMap[srv]
	}

	return client, nil
}

// NewServiceClient create a ServiceClient which was assembled from ServiceCatalog.
// If you want to add new ServiceClient, please make sure the catalog was already in allServiceCatalog.
// the endpoint likes https://{Name}.{Region}.myhuaweicloud.com/{Version}/{project_id}/{ResourceBase}
func (c *Config) NewServiceClient(srv, region string) (*golangsdk.ServiceClient, error) {
	serviceCatalog, ok := allServiceCatalog[srv]
	if !ok {
		return nil, fmt.Errorf("service type %s is invalid or not supportted", srv)
	}
	// update the service catalog name if necessary
	if name := getServiceCatalogNameByRegion(srv, region); name != "" {
		serviceCatalog.Name = name
	}

	if !c.SecurityKeyExpiresAt.IsZero() {
		c.SecurityKeyLock.Lock()
		defer c.SecurityKeyLock.Unlock()
		timeNow := time.Now().Unix()
		expairesAtInt := c.SecurityKeyExpiresAt.Unix()
		if timeNow+keyExpiresDuration > expairesAtInt {
			err := c.reloadSecurityKey()
			if err != nil {
				return nil, err
			}
		}
	}

	client := c.HwClient
	if serviceCatalog.Admin {
		client = c.DomainClient
	}

	if endpoint, ok := c.Endpoints[srv]; ok {
		if region != "" && region != c.Region {
			return nil, fmt.Errorf("Resource-level region must be the same as Provider-level region when using customizing endpoints")
		}
		return c.newServiceClientByEndpoint(client, srv, endpoint)
	}
	return c.newServiceClientByName(client, serviceCatalog, region)
}

func (c *Config) newServiceClientByName(client *golangsdk.ProviderClient, catalog ServiceCatalog, region string) (*golangsdk.ServiceClient, error) {
	if catalog.Name == "" {
		return nil, fmt.Errorf("must specify the service name")
	}

	// Custom Resource-level region only supports AK/SK authentication.
	// If set it when using non AK/SK authentication, then it must be the same as Provider-level region.
	if region != c.Region && (c.AccessKey == "" || c.SecretKey == "") {
		return nil, fmt.Errorf("Resource-level region must be the same as Provider-level region when using non AK/SK authentication if Resource-level region set")
	}

	c.RPLock.Lock()
	defer c.RPLock.Unlock()
	projectID, ok := c.RegionProjectIDMap[region]
	if !ok {
		// Not find in the map, then try to query and store.
		err := c.loadUserProjects(client, region)
		if err != nil {
			return nil, err
		}
		projectID = c.RegionProjectIDMap[region]
	}

	// update ProjectID and region in ProviderClient
	clone := new(golangsdk.ProviderClient)
	*clone = *client
	clone.ProjectID = projectID
	clone.AKSKAuthOptions.ProjectId = projectID
	clone.AKSKAuthOptions.Region = region

	sc := &golangsdk.ServiceClient{
		ProviderClient: clone,
	}

	if catalog.Scope == "global" && !c.RegionClient {
		sc.Endpoint = fmt.Sprintf("https://%s.%s/", catalog.Name, c.Cloud)
	} else {
		sc.Endpoint = fmt.Sprintf("https://%s.%s.%s/", catalog.Name, region, c.Cloud)
	}

	sc.ResourceBase = sc.Endpoint
	if catalog.Version != "" {
		sc.ResourceBase = sc.ResourceBase + catalog.Version + "/"
	}
	if !catalog.WithOutProjectID {
		sc.ResourceBase = sc.ResourceBase + projectID + "/"
	}
	if catalog.ResourceBase != "" {
		sc.ResourceBase = sc.ResourceBase + catalog.ResourceBase + "/"
	}

	return sc, nil
}

// newServiceClientByEndpoint returns a ServiceClient which the endpoint was initialized by customer
// the format of customer endpoint likes https://{Name}.{Region}.xxxx.com
func (c *Config) newServiceClientByEndpoint(client *golangsdk.ProviderClient, srv, endpoint string) (*golangsdk.ServiceClient, error) {
	catalog, ok := allServiceCatalog[srv]
	if !ok {
		return nil, fmt.Errorf("service type %s is invalid or not supportted", srv)
	}

	// Copy the client to prevent interference with the original data.
	clone := new(golangsdk.ProviderClient)
	*clone = *client
	sc := &golangsdk.ServiceClient{
		ProviderClient: clone,
		Endpoint:       endpoint,
	}

	sc.ResourceBase = sc.Endpoint
	if catalog.Version != "" {
		sc.ResourceBase = sc.ResourceBase + catalog.Version + "/"
	}
	if !catalog.WithOutProjectID {
		sc.ResourceBase = sc.ResourceBase + client.ProjectID + "/"
	}
	if catalog.ResourceBase != "" {
		sc.ResourceBase = sc.ResourceBase + catalog.ResourceBase + "/"
	}
	return sc, nil
}

func (c *Config) getDomainID() (string, error) {
	identityClient, err := c.IdentityV3Client(c.Region)
	if err != nil {
		return "", fmt.Errorf("Error creating IAM client: %s", err)
	}
	// ResourceBase: https://iam.{CLOUD}/v3/auth/
	identityClient.ResourceBase += "auth/"

	// the List request does not support query options
	allPages, err := domains.List(identityClient, nil).AllPages()
	if err != nil {
		return "", fmt.Errorf("List domains failed, err=%s", err)
	}

	all, err := domains.ExtractDomains(allPages)
	if err != nil {
		return "", fmt.Errorf("Extract domains failed, err=%s", err)
	}

	if len(all) == 0 {
		return "", fmt.Errorf("domain was not found")
	}

	if c.DomainName != "" && c.DomainName != all[0].Name {
		return "", fmt.Errorf("domain %s was not found, got %s", c.DomainName, all[0].Name)
	}

	return all[0].ID, nil
}

func (c *Config) getUserIDbyName(name string) (string, error) {
	identityClient, err := c.IdentityV3Client(c.Region)
	if err != nil {
		return "", fmt.Errorf("Error creating IAM client: %s", err)
	}

	opts := users.ListOpts{
		Name: name,
	}
	allPages, err := users.List(identityClient, opts).AllPages()
	if err != nil {
		return "", fmt.Errorf("query IAM user %s failed, err=%s", name, err)
	}

	all, err := users.ExtractUsers(allPages)
	if err != nil {
		return "", fmt.Errorf("Extract users failed, err=%s", err)
	}

	if len(all) == 0 {
		return "", fmt.Errorf("IAM user %s was not found", name)
	}

	if name != "" && name != all[0].Name {
		return "", fmt.Errorf("IAM user %s was not found, got %s", name, all[0].Name)
	}

	return all[0].ID, nil
}

// loadUserProjects will query the region-projectId pair and store it into RegionProjectIDMap
func (c *Config) loadUserProjects(client *golangsdk.ProviderClient, region string) error {

	log.Printf("[DEBUG] Load project ID for region: %s", region)
	domainID := client.DomainID
	opts := projects.ListOpts{
		DomainID: domainID,
		Name:     region,
	}
	sc := new(golangsdk.ServiceClient)
	sc.Endpoint = c.IdentityEndpoint + "/"
	sc.ProviderClient = client
	allPages, err := projects.List(sc, &opts).AllPages()
	if err != nil {
		return fmt.Errorf("List projects failed, err=%s", err)
	}

	all, err := projects.ExtractProjects(allPages)
	if err != nil {
		return fmt.Errorf("Extract projects failed, err=%s", err)
	}

	if len(all) == 0 {
		return fmt.Errorf("Wrong name or no access to the region: %s", region)
	}

	for _, item := range all {
		log.Printf("[DEBUG] add %s/%s to region and project map", item.Name, item.ID)
		c.RegionProjectIDMap[item.Name] = item.ID
	}
	return nil
}

// GetProjectID is used to get the project ID for services
func (c *Config) GetProjectID(region string) string {
	c.RPLock.Lock()
	defer c.RPLock.Unlock()

	projectID, ok := c.RegionProjectIDMap[region]
	if !ok {
		// Not find in the map, then try to query and store.
		if err := c.loadUserProjects(c.DomainClient, region); err != nil {
			log.Printf("[WARN] can not find the project ID of %s: %s", region, err)
			return ""
		}
		projectID = c.RegionProjectIDMap[region]
	}

	return projectID
}

// GetRegion returns the region that was specified in the resource. If a
// region was not set, the provider-level region is checked. The provider-level
// region can either be set by the region argument or by HW_REGION_NAME.
func (c *Config) GetRegion(d *schema.ResourceData) string {
	if v, ok := d.GetOk("region"); ok {
		return v.(string)
	}

	return c.Region
}

// GetEnterpriseProjectID returns the enterprise_project_id that was specified in the resource.
// If it was not set, the provider-level value is checked. The provider-level value can
// either be set by the `enterprise_project_id` argument or by HW_ENTERPRISE_PROJECT_ID.
// If the provider-level value
func (c *Config) GetEnterpriseProjectID(d *schema.ResourceData, defaultEps ...string) string {
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		return v.(string)
	}

	if c.EnterpriseProjectID != "" {
		return c.EnterpriseProjectID
	}
	if len(defaultEps) > 0 {
		return defaultEps[0]
	}
	return ""
}

// CheckValueInterchange checks if the new value of key1 is equal to the old value of key2,
// and the new value of key2 is equal to the old value of key1.
func CheckValueInterchange(d *schema.ResourceDiff, key1, key2 string) (isKey1NewEqualKey2Old bool, isKey2NewEqualKey1Old bool) {
	oldKey1Value, newKey1Value := d.GetChange(key1)
	oldKey2Value, newKey2Value := d.GetChange(key2)

	// Check if any of the values are empty strings, in which case we return false for both checks.
	if oldKey1Value.(string) == "" || newKey1Value.(string) == "" ||
		oldKey2Value.(string) == "" || newKey2Value.(string) == "" {
		return false, false
	}

	isKey1NewEqualKey2Old = newKey1Value.(string) == oldKey2Value.(string)
	isKey2NewEqualKey1Old = newKey2Value.(string) == oldKey1Value.(string)

	return isKey1NewEqualKey2Old, isKey2NewEqualKey1Old
}

// GetForceNew returns the enable_force_new that was specified in the resource.
// If it was not set, the provider-level value is checked. The provider-level value can
// either be set by the `enable_force_new` argument or by HW_ENABLE_FORCE_NEW.
func (c *Config) GetForceNew(d *schema.ResourceDiff) bool {
	if v, ok := d.GetOk("enable_force_new"); ok {
		res, err := strconv.ParseBool(v.(string))
		if err != nil {
			return false
		}
		return res
	}

	return c.EnableForceNew
}

// ********** client for Global Service **********
func (c *Config) IAMV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("iam", region)
}

func (c *Config) IAMV3GlobalClient() (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("iam_global", "")
}

func (c *Config) IdentityV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("identity", region)
}

func (c *Config) IAMNoVersionClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("iam_no_version", region)
}

func (c *Config) IdentityV3ExtClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("identity_ext", region)
}

func (c *Config) CdnV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cdn", region)
}

func (c *Config) EnterpriseProjectClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("eps", region)
}

func (c *Config) TmsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("tms", region)
}

func (c *Config) TmsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("tmsv2", region)
}

// ********** client for Compute **********
func (c *Config) ComputeV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("ecs", region)
}

func (c *Config) ComputeV11Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("ecsv11", region)
}

func (c *Config) ComputeV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("ecsv21", region)
}

func (c *Config) AutoscalingV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("autoscaling", region)
}

func (c *Config) ImageV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("imsv1", region)
}

func (c *Config) ImageV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("ims", region)
}

func (c *Config) CceV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("ccev1", region)
}

func (c *Config) CceV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cce", region)
}

func (c *Config) CceAddonV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cce_addon", region)
}

func (c *Config) AomV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("aom", region)
}

func (c *Config) CciV1BetaClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cciv1_bata", region)
}

func (c *Config) CciV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cci", region)
}

func (c *Config) FgsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("fgs", region)
}

func (c *Config) SwrV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("swr", region)
}

func (c *Config) BmsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("bms", region)
}

// ********** client for Storage **********
func (c *Config) BlockStorageV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("evsv1", region)
}

func (c *Config) BlockStorageV21Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("evsv21", region)
}

func (c *Config) BlockStorageV5Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("evsv5", region)
}

func (c *Config) BlockStorageV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("evs", region)
}

func (c *Config) SfsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("sfs", region)
}

func (c *Config) SfsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("sfs-turbo", region)
}

func (c *Config) CbrV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cbr", region)
}

func (c *Config) CsbsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("csbs", region)
}

func (c *Config) VbsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("vbs", region)
}

func (c *Config) SdrsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("sdrs", region)
}

// ********** client for Network **********
func (c *Config) NetworkingV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("vpc", region)
}

// NetworkingV2Client returns a ServiceClient for neutron APIs
// the endpoint likes: https://vpc.{region}.myhuaweicloud.com/v2.0/
func (c *Config) NetworkingV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("networkv2", region)
}

func (c *Config) NetworkingV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("vpcv3", region)
}

// VPCEPClient returns a ServiceClient for VPC Endpoint APIs
// the endpoint likes: https://vpcep.{region}.myhuaweicloud.com/v1/{project_id}/
func (c *Config) VPCEPClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("vpcep", region)
}

func (c *Config) NatGatewayClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("nat", region)
}

// NatV3Client has the endpoint: https://nat.{{region}}/{{cloud}}/v3/
func (c *Config) NatV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("natv3", region)
}

// ElbV2Client is the client for elb v2.0 (openstack) api
func (c *Config) ElbV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("elbv2", region)
}

// ElbV3Client is the client for elb v3 api
func (c *Config) ElbV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("elbv3", region)
}

// LoadBalancerClient is the client for elb v2 api
func (c *Config) LoadBalancerClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("elb", region)
}

func (c *Config) FwV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("networkv2", region)
}

func (c *Config) DnsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dns", region)
}

func (c *Config) DnsWithRegionClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dns_region", region)
}

func (c *Config) DNSV21Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dnsv21", region)
}

func (c *Config) ErV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("er", region)
}

func (c *Config) DcV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dc", region)
}

// ********** client for Management **********
func (c *Config) CtsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cts", region)
}

func (c *Config) CesV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("ces", region)
}

func (c *Config) CesV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cesv2", region)
}

func (c *Config) LtsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("lts", region)
}

func (c *Config) SmnV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("smn", region)
}

func (c *Config) SmnV2TagClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("smn-tag", region)
}

func (c *Config) RmsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("rms", region)
}

// ********** client for Security **********
func (c *Config) AntiDDosV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("anti-ddos", region)
}

func (c *Config) AntiDDosV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("anti-ddosv2", region)
}

func (c *Config) AadV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("aad", region)
}

func (c *Config) KmsKeyV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("kms", region)
}

func (c *Config) KmsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("kmsv1", region)
}

func (c *Config) KmsV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("kmsv3", region)
}

// WafV1Client is not avaliable in HuaweiCloud, will be imported by other clouds
func (c *Config) WafV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("waf", region)
}

func (c *Config) WafDedicatedV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("waf-dedicated", region)
}

// ********** client for Enterprise Intelligence **********
func (c *Config) MrsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("mrs", region)
}

func (c *Config) MrsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("mrsv2", region)
}

func (c *Config) DwsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dws", region)
}

func (c *Config) DwsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dwsv2", region)
}

func (c *Config) DliV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dli", region)
}

func (c *Config) DliV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dliv2", region)
}

func (c *Config) DliV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dliv3", region)
}

func (c *Config) DisV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dis", region)
}

func (c *Config) DisV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("disv3", region)
}

func (c *Config) CssV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("css", region)
}

func (c *Config) CssV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cssv2", region)
}

func (c *Config) CloudStreamV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cs", region)
}

func (c *Config) CloudtableV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cloudtable", region)
}

func (c *Config) CdmV11Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cdm", region)
}

func (c *Config) GesV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("ges", region)
}

func (c *Config) ModelArtsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("modelarts", region)
}

func (c *Config) ModelArtsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("modelartsv2", region)
}

func (c *Config) DataArtsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dataarts", region)
}

func (c *Config) WorkspaceV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("workspace", region)
}

// ********** client for Application **********
func (c *Config) ApiGatewayV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("apig", region)
}

func (c *Config) ApigV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("apigv2", region)
}

func (c *Config) BcsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("bcs", region)
}

func (c *Config) CseV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cse", region)
}

func (c *Config) DcsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dcsv1", region)
}

func (c *Config) DcsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dcs", region)
}

func (c *Config) DmsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dms", region)
}

func (c *Config) DmsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dmsv2", region)
}

func (c *Config) ServiceStageV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("servicestage", region)
}

func (c *Config) ServiceStageV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("servicestagev2", region)
}

func (c *Config) EgV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("eg", region)
}

// ********** client for Database **********
func (c *Config) RdsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("rdsv1", region)
}

func (c *Config) RdsV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("rds", region)
}

func (c *Config) RdsV31Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("rdsv31", region)
}

func (c *Config) DdsV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dds", region)
}

func (c *Config) GeminiDBV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("geminidb", region)
}

func (c *Config) GeminiDBV31Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("geminidbv31", region)
}

func (c *Config) OpenGaussV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("opengauss", region)
}

func (c *Config) GaussdbV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("gaussdb", region)
}

func (c *Config) DrsV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("drs", region)
}

func (c *Config) DrsV5Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("drsv5", region)
}

// ********** client for edge / IoT **********

// IECV1Client returns a ServiceClient for IEC Endpoint APIs
func (c *Config) IECV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("iec", region)
}

// ********** client for Others **********
func (c *Config) BssV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("bss", region)
}

func (c *Config) BssV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("bssv2", region)
}

func (c *Config) MaasV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("oms", region)
}

func (c *Config) SmsV3Client(region string) (*golangsdk.ServiceClient, error) {
	if c.GetWebsiteType() == InternationalSite {
		return c.NewServiceClient("sms-intl", region)
	}
	return c.NewServiceClient("sms", region)
}

func (c *Config) ScmV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("scm", region)
}

// the following clients are used for Joint-Operation Cloud only

// MlsV1Client has the endpoint: https://mls.{{region}}/{{cloud}}/v1.0/{{project_id}}
func (c *Config) MlsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("mls", region)
}

// NatV2Client has the endpoint: https://nat.{{region}}/{{cloud}}/v2.0/
func (c *Config) NatV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("natv2", region)
}

// KooGalleryV1Client has the endpoint: https://mkt.{{cloud}}/v1/
func (c *Config) KooGalleryV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("mkt", region)
}

func (c *Config) VpnV5Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("vpn", region)
}

func (c *Config) StsClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("sts", region)
}
