package openstack

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"github.com/huaweicloud/golangsdk"
	tokens2 "github.com/huaweicloud/golangsdk/openstack/identity/v2/tokens"
	tokens3 "github.com/huaweicloud/golangsdk/openstack/identity/v3/tokens"
	"github.com/huaweicloud/golangsdk/openstack/utils"
)

const (
	// v2 represents Keystone v2.
	// It should never increase beyond 2.0.
	v2 = "v2.0"

	// v3 represents Keystone v3.
	// The version can be anything from v3 to v3.x.
	v3 = "v3"
)

/*
NewClient prepares an unauthenticated ProviderClient instance.
Most users will probably prefer using the AuthenticatedClient function
instead.

This is useful if you wish to explicitly control the version of the identity
service that's used for authentication explicitly, for example.

A basic example of using this would be:

	ao, err := openstack.AuthOptionsFromEnv()
	provider, err := openstack.NewClient(ao.IdentityEndpoint)
	client, err := openstack.NewIdentityV3(provider, golangsdk.EndpointOpts{})
*/
func NewClient(endpoint string) (*golangsdk.ProviderClient, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	u.RawQuery, u.Fragment = "", ""

	var base string
	versionRe := regexp.MustCompile("v[0-9.]+/?")
	if version := versionRe.FindString(u.Path); version != "" {
		base = strings.Replace(u.String(), version, "", -1)
	} else {
		base = u.String()
	}

	endpoint = golangsdk.NormalizeURL(endpoint)
	base = golangsdk.NormalizeURL(base)

	p := new(golangsdk.ProviderClient)
	p.IdentityBase = base
	p.IdentityEndpoint = endpoint
	p.UseTokenLock()

	return p, nil
}

/*
AuthenticatedClient logs in to an OpenStack cloud found at the identity endpoint
specified by the options, acquires a token, and returns a Provider Client
instance that's ready to operate.

If the full path to a versioned identity endpoint was specified  (example:
http://example.com:5000/v3), that path will be used as the endpoint to query.

If a versionless endpoint was specified (example: http://example.com:5000/),
the endpoint will be queried to determine which versions of the identity service
are available, then chooses the most recent or most supported version.

Example:

	ao, err := openstack.AuthOptionsFromEnv()
	provider, err := openstack.AuthenticatedClient(ao)
	client, err := openstack.NewNetworkV2(client, golangsdk.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
*/
func AuthenticatedClient(options golangsdk.AuthOptions) (*golangsdk.ProviderClient, error) {
	client, err := NewClient(options.IdentityEndpoint)
	if err != nil {
		return nil, err
	}

	err = Authenticate(client, options)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Authenticate or re-authenticate against the most recent identity service
// supported at the provided endpoint.
func Authenticate(client *golangsdk.ProviderClient, options golangsdk.AuthOptions) error {
	versions := []*utils.Version{
		{ID: v2, Priority: 20, Suffix: "/v2.0/"},
		{ID: v3, Priority: 30, Suffix: "/v3/"},
	}

	chosen, endpoint, err := utils.ChooseVersion(client, versions)
	if err != nil {
		return err
	}

	switch chosen.ID {
	case v2:
		return v2auth(client, endpoint, options, golangsdk.EndpointOpts{})
	case v3:
		return v3auth(client, endpoint, &options, golangsdk.EndpointOpts{})
	default:
		// The switch statement must be out of date from the versions list.
		return fmt.Errorf("Unrecognized identity version: %s", chosen.ID)
	}
}

// AuthenticateV2 explicitly authenticates against the identity v2 endpoint.
func AuthenticateV2(client *golangsdk.ProviderClient, options golangsdk.AuthOptions, eo golangsdk.EndpointOpts) error {
	return v2auth(client, "", options, eo)
}

func v2auth(client *golangsdk.ProviderClient, endpoint string, options golangsdk.AuthOptions, eo golangsdk.EndpointOpts) error {
	v2Client, err := NewIdentityV2(client, eo)
	if err != nil {
		return err
	}

	if endpoint != "" {
		v2Client.Endpoint = endpoint
	}

	v2Opts := tokens2.AuthOptions{
		IdentityEndpoint: options.IdentityEndpoint,
		Username:         options.Username,
		Password:         options.Password,
		TenantID:         options.TenantID,
		TenantName:       options.TenantName,
		AllowReauth:      options.AllowReauth,
		TokenID:          options.TokenID,
	}

	result := tokens2.Create(v2Client, v2Opts)

	token, err := result.ExtractToken()
	if err != nil {
		return err
	}

	catalog, err := result.ExtractServiceCatalog()
	if err != nil {
		return err
	}

	if options.AllowReauth {
		client.ReauthFunc = func() error {
			client.TokenID = ""
			return v2auth(client, endpoint, options, eo)
		}
	}
	client.TokenID = token.ID
	client.ProjectID = token.Tenant.ID
	client.EndpointLocator = func(opts golangsdk.EndpointOpts) (string, error) {
		return V2EndpointURL(catalog, opts)
	}

	return nil
}

// AuthenticateV3 explicitly authenticates against the identity v3 service.
func AuthenticateV3(client *golangsdk.ProviderClient, options tokens3.AuthOptionsBuilder, eo golangsdk.EndpointOpts) error {
	return v3auth(client, "", options, eo)
}

func v3auth(client *golangsdk.ProviderClient, endpoint string, opts tokens3.AuthOptionsBuilder, eo golangsdk.EndpointOpts) error {
	// Override the generated service endpoint with the one returned by the version endpoint.
	v3Client, err := NewIdentityV3(client, eo)
	if err != nil {
		return err
	}

	if endpoint != "" {
		v3Client.Endpoint = endpoint
	}

	result := tokens3.Create(v3Client, opts)

	token, err := result.ExtractToken()
	if err != nil {
		return err
	}

	project, err := result.ExtractProject()
	if err != nil {
		return err
	}

	catalog, err := result.ExtractServiceCatalog()
	if err != nil {
		return err
	}

	client.TokenID = token.ID
	client.ProjectID = project.ID

	if opts.CanReauth() {
		client.ReauthFunc = func() error {
			client.TokenID = ""
			return v3auth(client, endpoint, opts, eo)
		}
	}
	client.EndpointLocator = func(opts golangsdk.EndpointOpts) (string, error) {
		return V3EndpointURL(catalog, opts)
	}

	return nil
}

// NewIdentityV2 creates a ServiceClient that may be used to interact with the
// v2 identity service.
func NewIdentityV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	endpoint := client.IdentityBase + "v2.0/"
	clientType := "identity"
	var err error
	if !reflect.DeepEqual(eo, golangsdk.EndpointOpts{}) {
		eo.ApplyDefaults(clientType)
		endpoint, err = client.EndpointLocator(eo)
		if err != nil {
			return nil, err
		}
	}

	return &golangsdk.ServiceClient{
		ProviderClient: client,
		Endpoint:       endpoint,
		Type:           clientType,
	}, nil
}

// NewIdentityV3 creates a ServiceClient that may be used to access the v3
// identity service.
func NewIdentityV3(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	endpoint := client.IdentityBase + "v3/"
	clientType := "identity"
	var err error
	if !reflect.DeepEqual(eo, golangsdk.EndpointOpts{}) {
		eo.ApplyDefaults(clientType)
		endpoint, err = client.EndpointLocator(eo)
		if err != nil {
			return nil, err
		}
	}

	// Ensure endpoint still has a suffix of v3.
	// This is because EndpointLocator might have found a versionless
	// endpoint and requests will fail unless targeted at /v3.
	if !strings.HasSuffix(endpoint, "v3/") {
		endpoint = endpoint + "v3/"
	}

	return &golangsdk.ServiceClient{
		ProviderClient: client,
		Endpoint:       endpoint,
		Type:           clientType,
	}, nil
}

func initClientOpts(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts, clientType string) (*golangsdk.ServiceClient, error) {
	sc := new(golangsdk.ServiceClient)
	eo.ApplyDefaults(clientType)
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return sc, err
	}
	sc.ProviderClient = client
	sc.Endpoint = url
	sc.Type = clientType
	return sc, nil
}

// NewObjectStorageV1 creates a ServiceClient that may be used with the v1
// object storage package.
func NewObjectStorageV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	return initClientOpts(client, eo, "object-store")
}

// NewComputeV2 creates a ServiceClient that may be used with the v2 compute
// package.
func NewComputeV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	return initClientOpts(client, eo, "compute")
}

// NewNetworkV2 creates a ServiceClient that may be used with the v2 network
// package.
func NewNetworkV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "network")
	sc.ResourceBase = sc.Endpoint + "v2.0/"
	return sc, err
}

// NewBlockStorageV1 creates a ServiceClient that may be used to access the v1
// block storage service.
func NewBlockStorageV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	return initClientOpts(client, eo, "volume")
}

// NewBlockStorageV2 creates a ServiceClient that may be used to access the v2
// block storage service.
func NewBlockStorageV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	return initClientOpts(client, eo, "volumev2")
}

// NewBlockStorageV3 creates a ServiceClient that may be used to access the v3 block storage service.
func NewBlockStorageV3(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	return initClientOpts(client, eo, "volumev3")
}

// NewSharedFileSystemV2 creates a ServiceClient that may be used to access the v2 shared file system service.
func NewSharedFileSystemV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	return initClientOpts(client, eo, "sharev2")
}

// NewCDNV1 creates a ServiceClient that may be used to access the OpenStack v1
// CDN service.
func NewCDNV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	return initClientOpts(client, eo, "cdn")
}

// NewOrchestrationV1 creates a ServiceClient that may be used to access the v1
// orchestration service.
func NewOrchestrationV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	return initClientOpts(client, eo, "orchestration")
}

// NewDBV1 creates a ServiceClient that may be used to access the v1 DB service.
func NewDBV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	return initClientOpts(client, eo, "database")
}

// NewDNSV2 creates a ServiceClient that may be used to access the v2 DNS
// service.
func NewDNSV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "dns")
	sc.ResourceBase = sc.Endpoint + "v2/"
	return sc, err
}

// NewImageServiceV2 creates a ServiceClient that may be used to access the v2
// image service.
func NewImageServiceV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "image")
	sc.ResourceBase = sc.Endpoint + "v2/"
	return sc, err
}

// NewLoadBalancerV2 creates a ServiceClient that may be used to access the v2
// load balancer service.
func NewLoadBalancerV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "load-balancer")
	sc.ResourceBase = sc.Endpoint + "v2.0/"
	return sc, err
}

// NewOtcV1 creates a ServiceClient that may be used with the v1 network package.
func NewElbV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts, otctype string) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "compute")
	//fmt.Printf("client=%+v.\n", sc)
	sc.Endpoint = strings.Replace(strings.Replace(sc.Endpoint, "ecs", otctype, 1), "/v2/", "/v1.0/", 1)
	//fmt.Printf("url=%s.\n", sc.Endpoint)
	sc.ResourceBase = sc.Endpoint
	sc.Type = otctype
	return sc, err
}

// NewSmnServiceV2 creates a ServiceClient that may be used to access the v2 Simple Message Notification service.
func NewSmnServiceV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {

	sc, err := initClientOpts(client, eo, "compute")
	sc.Endpoint = strings.Replace(sc.Endpoint, "ecs", "smn", 1)
	sc.ResourceBase = sc.Endpoint + "notifications/"
	sc.Type = "smn"
	return sc, err
}

//NewRdsServiceV1 creates the a ServiceClient that may be used to access the v1
//rds service which is a service of db instances management.
func NewRdsServiceV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	newsc, err := initClientOpts(client, eo, "compute")
	rdsendpoint := strings.Replace(strings.Replace(newsc.Endpoint, "ecs", "rds", 1), "/v2/", "/rds/v1/", 1)
	newsc.Endpoint = rdsendpoint
	newsc.ResourceBase = rdsendpoint
	newsc.Type = "rds"
	return newsc, err
}

func NewCESClient(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "ces")
	if err != nil {
		return nil, err
	}
	sc.ResourceBase = sc.Endpoint
	return sc, err
}
