package config

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/chnsz/golangsdk"
	huaweisdk "github.com/chnsz/golangsdk/openstack"
	"github.com/chnsz/golangsdk/openstack/identity/v3/domains"
	"github.com/chnsz/golangsdk/openstack/identity/v3/projects"
	"github.com/chnsz/golangsdk/openstack/obs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/pathorcontents"
)

const (
	obsLogFile         string = "./.obs-sdk.log"
	obsLogFileSize10MB int64  = 1024 * 1024 * 10
)

type Config struct {
	AccessKey           string
	SecretKey           string
	CACertFile          string
	ClientCertFile      string
	ClientKeyFile       string
	DomainID            string
	DomainName          string
	IdentityEndpoint    string
	Insecure            bool
	Password            string
	Region              string
	TenantID            string
	TenantName          string
	Token               string
	SecurityToken       string
	Username            string
	UserID              string
	AgencyName          string
	AgencyDomainName    string
	DelegatedProject    string
	Cloud               string
	MaxRetries          int
	TerraformVersion    string
	RegionClient        bool
	EnterpriseProjectID string

	HwClient     *golangsdk.ProviderClient
	DomainClient *golangsdk.ProviderClient

	// the custom endpoints used to override the default endpoint URL
	Endpoints map[string]string

	// RegionProjectIDMap is a map which stores the region-projectId pairs,
	// and region name will be the key and projectID will be the value in this map.
	RegionProjectIDMap map[string]string

	// RPLock is used to make the accessing of RegionProjectIDMap serial,
	// prevent sending duplicate query requests
	RPLock *sync.Mutex
}

func (c *Config) LoadAndValidate() error {
	if c.MaxRetries < 0 {
		return fmt.Errorf("max_retries should be a positive value")
	}

	err := fmt.Errorf("Must config token or aksk or username password to be authorized")

	if c.Token != "" {
		err = buildClientByToken(c)

	} else if c.AccessKey != "" && c.SecretKey != "" {
		err = buildClientByAKSK(c)

	} else if c.Password != "" {
		if c.Username == "" && c.UserID == "" {
			err = fmt.Errorf("\"password\": one of `user_name, user_id` must be specified")
		} else {
			err = buildClientByPassword(c)
		}

	}
	if err != nil {
		return err
	}

	if c.HwClient != nil && c.HwClient.ProjectID != "" {
		c.RegionProjectIDMap[c.Region] = c.HwClient.ProjectID
	}

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

	return nil
}

func generateTLSConfig(c *Config) (*tls.Config, error) {
	config := &tls.Config{}
	if c.CACertFile != "" {
		caCert, _, err := pathorcontents.Read(c.CACertFile)
		if err != nil {
			return nil, fmt.Errorf("Error reading CA Cert: %s", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))
		config.RootCAs = caCertPool
	}

	if c.Insecure {
		config.InsecureSkipVerify = true
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

func genClient(c *Config, ao golangsdk.AuthOptionsProvider) (*golangsdk.ProviderClient, error) {
	client, err := huaweisdk.NewClient(ao.GetIdentityEndpoint())
	if err != nil {
		return nil, err
	}

	// Set UserAgent
	client.UserAgent.Prepend("terraform-provider-iac")

	config, err := generateTLSConfig(c)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: config}

	client.HTTPClient = http.Client{
		Transport: &LogRoundTripper{
			Rt:         transport,
			OsDebug:    logging.IsDebugOrHigher(),
			MaxRetries: c.MaxRetries,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if client.AKSKAuthOptions.AccessKey != "" {
				golangsdk.ReSign(req, golangsdk.SignOptions{
					AccessKey: client.AKSKAuthOptions.AccessKey,
					SecretKey: client.AKSKAuthOptions.SecretKey,
				})
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

func buildClientByToken(c *Config) error {
	var pao, dao golangsdk.AuthOptions

	if c.AgencyDomainName != "" && c.AgencyName != "" {
		pao = golangsdk.AuthOptions{
			AgencyName:       c.AgencyName,
			AgencyDomainName: c.AgencyDomainName,
			DelegatedProject: c.DelegatedProject,
		}

		dao = golangsdk.AuthOptions{
			AgencyName:       c.AgencyName,
			AgencyDomainName: c.AgencyDomainName,
		}
	} else {
		pao = golangsdk.AuthOptions{
			DomainID:   c.DomainID,
			DomainName: c.DomainName,
			TenantID:   c.TenantID,
			TenantName: c.TenantName,
		}

		dao = golangsdk.AuthOptions{
			DomainID:   c.DomainID,
			DomainName: c.DomainName,
		}
	}

	for _, ao := range []*golangsdk.AuthOptions{&pao, &dao} {
		ao.IdentityEndpoint = c.IdentityEndpoint
		ao.TokenID = c.Token

	}
	return genClients(c, pao, dao)
}

func buildClientByAKSK(c *Config) error {
	var pao, dao golangsdk.AKSKAuthOptions

	if c.AgencyDomainName != "" && c.AgencyName != "" {
		pao = golangsdk.AKSKAuthOptions{
			DomainID:         c.DomainID,
			Domain:           c.DomainName,
			AgencyName:       c.AgencyName,
			AgencyDomainName: c.AgencyDomainName,
			DelegatedProject: c.DelegatedProject,
		}

		dao = golangsdk.AKSKAuthOptions{
			DomainID:         c.DomainID,
			Domain:           c.DomainName,
			AgencyName:       c.AgencyName,
			AgencyDomainName: c.AgencyDomainName,
		}
	} else {
		pao = golangsdk.AKSKAuthOptions{
			BssDomainID: c.DomainID,
			BssDomain:   c.DomainName,
			ProjectName: c.TenantName,
			ProjectId:   c.TenantID,
		}

		dao = golangsdk.AKSKAuthOptions{
			DomainID: c.DomainID,
			Domain:   c.DomainName,
		}
	}

	for _, ao := range []*golangsdk.AKSKAuthOptions{&pao, &dao} {
		ao.IdentityEndpoint = c.IdentityEndpoint
		ao.AccessKey = c.AccessKey
		ao.SecretKey = c.SecretKey
		if c.SecurityToken != "" {
			ao.SecurityToken = c.SecurityToken
			ao.WithUserCatalog = true
		}
	}
	return genClients(c, pao, dao)
}

func buildClientByPassword(c *Config) error {
	var pao, dao golangsdk.AuthOptions

	if c.AgencyDomainName != "" && c.AgencyName != "" {
		pao = golangsdk.AuthOptions{
			DomainID:         c.DomainID,
			DomainName:       c.DomainName,
			AgencyName:       c.AgencyName,
			AgencyDomainName: c.AgencyDomainName,
			DelegatedProject: c.DelegatedProject,
		}

		dao = golangsdk.AuthOptions{
			DomainID:         c.DomainID,
			DomainName:       c.DomainName,
			AgencyName:       c.AgencyName,
			AgencyDomainName: c.AgencyDomainName,
		}
	} else {
		pao = golangsdk.AuthOptions{
			DomainID:   c.DomainID,
			DomainName: c.DomainName,
			TenantID:   c.TenantID,
			TenantName: c.TenantName,
		}

		dao = golangsdk.AuthOptions{
			DomainID:   c.DomainID,
			DomainName: c.DomainName,
		}
	}

	for _, ao := range []*golangsdk.AuthOptions{&pao, &dao} {
		ao.IdentityEndpoint = c.IdentityEndpoint
		ao.Password = c.Password
		ao.Username = c.Username
		ao.UserID = c.UserID
	}
	return genClients(c, pao, dao)
}

func genClients(c *Config, pao, dao golangsdk.AuthOptionsProvider) error {
	client, err := genClient(c, pao)
	if err != nil {
		return err
	}
	c.HwClient = client

	client, err = genClient(c, dao)
	if err == nil {
		c.DomainClient = client
	}
	return err
}

func getObsEndpoint(c *Config, region string) string {
	if endpoint, ok := c.Endpoints["obs"]; ok {
		return endpoint
	}
	return fmt.Sprintf("https://obs.%s.%s/", region, c.Cloud)
}

func (c *Config) ObjectStorageClientWithSignature(region string) (*obs.ObsClient, error) {
	if c.AccessKey == "" || c.SecretKey == "" {
		return nil, fmt.Errorf("missing credentials for OBS, need access_key and secret_key values for provider")
	}

	// init log
	if logging.IsDebugOrHigher() {
		if err := obs.InitLog(obsLogFile, obsLogFileSize10MB, 10, obs.LEVEL_DEBUG, false); err != nil {
			log.Printf("[WARN] initial obs sdk log failed: %s", err)
		}
	}

	obsEndpoint := getObsEndpoint(c, region)
	if c.SecurityToken != "" {
		return obs.New(c.AccessKey, c.SecretKey, obsEndpoint,
			obs.WithSignature("OBS"), obs.WithSecurityToken(c.SecurityToken))
	}
	return obs.New(c.AccessKey, c.SecretKey, obsEndpoint, obs.WithSignature("OBS"))
}

func (c *Config) ObjectStorageClient(region string) (*obs.ObsClient, error) {
	if c.AccessKey == "" || c.SecretKey == "" {
		return nil, fmt.Errorf("missing credentials for OBS, need access_key and secret_key values for provider")
	}

	// init log
	if logging.IsDebugOrHigher() {
		if err := obs.InitLog(obsLogFile, obsLogFileSize10MB, 10, obs.LEVEL_DEBUG, false); err != nil {
			log.Printf("[WARN] initial obs sdk log failed: %s", err)
		}
	}

	obsEndpoint := getObsEndpoint(c, region)
	if c.SecurityToken != "" {
		return obs.New(c.AccessKey, c.SecretKey, obsEndpoint, obs.WithSecurityToken(c.SecurityToken))
	}
	return obs.New(c.AccessKey, c.SecretKey, obsEndpoint)
}

// NewServiceClient create a ServiceClient which was assembled from ServiceCatalog.
// If you want to add new ServiceClient, please make sure the catalog was already in allServiceCatalog.
// the endpoint likes https://{Name}.{Region}.myhuaweicloud.com/{Version}/{project_id}/{ResourceBase}
func (c *Config) NewServiceClient(srv, region string) (*golangsdk.ServiceClient, error) {
	serviceCatalog, ok := allServiceCatalog[srv]
	if !ok {
		return nil, fmt.Errorf("service type %s is invalid or not supportted", srv)
	}

	client := c.HwClient
	if serviceCatalog.Admin {
		client = c.DomainClient
	}

	if endpoint, ok := c.Endpoints[srv]; ok {
		return c.newServiceClientByEndpoint(client, srv, endpoint)
	}
	return c.newServiceClientByName(client, serviceCatalog, region)
}

func (c *Config) newServiceClientByName(client *golangsdk.ProviderClient, catalog ServiceCatalog, region string) (*golangsdk.ServiceClient, error) {
	if catalog.Name == "" || catalog.Version == "" {
		return nil, fmt.Errorf("must specify the service name and api version")
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
		projectID, _ = c.RegionProjectIDMap[region]
	}

	sc := new(golangsdk.ServiceClient)

	clone := new(golangsdk.ProviderClient)
	*clone = *client
	clone.ProjectID = projectID
	clone.AKSKAuthOptions.ProjectId = projectID
	clone.AKSKAuthOptions.Region = region
	sc.ProviderClient = clone

	if catalog.Scope == "global" && !c.RegionClient {
		sc.Endpoint = fmt.Sprintf("https://%s.%s/", catalog.Name, c.Cloud)
	} else {
		sc.Endpoint = fmt.Sprintf("https://%s.%s.%s/", catalog.Name, region, c.Cloud)
	}

	sc.ResourceBase = sc.Endpoint + catalog.Version + "/"
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

	sc := &golangsdk.ServiceClient{
		ProviderClient: client,
		Endpoint:       endpoint,
	}
	sc.ResourceBase = sc.Endpoint + catalog.Version + "/"
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

// loadUserProjects will query the region-projectId pair and store it into RegionProjectIDMap
func (c *Config) loadUserProjects(client *golangsdk.ProviderClient, region string) error {

	log.Printf("Load projectID for region: %s", region)
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
		c.RegionProjectIDMap[item.Name] = item.ID
	}
	return nil
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
func (c *Config) GetEnterpriseProjectID(d *schema.ResourceData) string {
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		return v.(string)
	}

	return c.EnterpriseProjectID
}

// ********** client for Global Service **********
func (c *Config) IAMV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("iam", region)
}

func (c *Config) IdentityV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("identity", region)
}

func (c *Config) CdnV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cdn", region)
}

func (c *Config) EnterpriseProjectClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("eps", region)
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

func (c *Config) ImageV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("ims", region)
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
	return c.NewServiceClient("cciv1", region)
}

func (c *Config) FgsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("fgsv2", region)
}

func (c *Config) SwrV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("swr", region)
}

func (c *Config) BmsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("bms", region)
}

// ********** client for Storage **********
func (c *Config) BlockStorageV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("volumev2", region)
}

func (c *Config) BlockStorageV3Client(region string) (*golangsdk.ServiceClient, error) {
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

// ********** client for Network **********
func (c *Config) NetworkingV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("vpc", region)
}

// NetworkingV2Client returns a ServiceClient for neutron APIs
// the endpoint likes: https://vpc.{region}.myhuaweicloud.com/v2.0/
func (c *Config) NetworkingV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("networkv2", region)
}

func (c *Config) SecurityGroupV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("security_group", region)
}

// VPCEPClient returns a ServiceClient for VPC Endpoint APIs
// the endpoint likes: https://vpcep.{region}.myhuaweicloud.com/v1/{project_id}/
func (c *Config) VPCEPClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("vpcep", region)
}

func (c *Config) NatGatewayClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("nat", region)
}

func (c *Config) ElasticLBClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("elb", region)
}

// client for v2.0 api
func (c *Config) ElbV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("elbv2", region)
}

// client for v3 api
func (c *Config) ElbV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("elbv3", region)
}

// client for v2 api
func (c *Config) LoadBalancerClient(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("loadbalancer", region)
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

// ********** client for Management **********
func (c *Config) CtsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("cts", region)
}

func (c *Config) CesV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("ces", region)
}

func (c *Config) LtsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("lts", region)
}

func (c *Config) SmnV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("smn", region)
}

// ********** client for Security **********
func (c *Config) AntiDDosV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("anti-ddos", region)
}

func (c *Config) KmsKeyV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("kms", region)
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

func (c *Config) DliV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dli", region)
}

func (c *Config) DisV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("disv2", region)
}

func (c *Config) CssV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("css", region)
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

// ********** client for Application **********
func (c *Config) ApiGatewayV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("apig", region)
}

func (c *Config) ApigV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("apig_v2", region)
}

func (c *Config) BcsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("bcs", region)
}

func (c *Config) DcsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dcsv1", region)
}

func (c *Config) DcsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dcsv2", region)
}

func (c *Config) DmsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dms", region)
}

func (c *Config) DmsV2Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dmsv2", region)
}

// ********** client for Database **********
func (c *Config) RdsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("rdsv1", region)
}

func (c *Config) RdsV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("rds", region)
}

func (c *Config) DdsV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("dds", region)
}

func (c *Config) GeminiDBV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("geminidb", region)
}

func (c *Config) OpenGaussV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("opengauss", region)
}

func (c *Config) GaussdbV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("gaussdb", region)
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

func (c *Config) OrchestrationV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("rts", region)
}

func (c *Config) MlsV1Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("mls", region)
}

func (c *Config) ScmV3Client(region string) (*golangsdk.ServiceClient, error) {
	return c.NewServiceClient("scm", region)
}
