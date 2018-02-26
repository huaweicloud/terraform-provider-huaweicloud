package openstack

import (
	"fmt"
	"strings"

	"github.com/huaweicloud/golangsdk"
	tokens2 "github.com/huaweicloud/golangsdk/openstack/identity/v2/tokens"
	tokens3 "github.com/huaweicloud/golangsdk/openstack/identity/v3/tokens"
	"github.com/huaweicloud/golangsdk/openstack/utils"
)

func GetProjectId(client *golangsdk.ProviderClient) (string, error) {
	versions := []*utils.Version{
		{ID: v2, Priority: 20, Suffix: "/v2.0/"},
		{ID: v3, Priority: 30, Suffix: "/v3/"},
	}

	chosen, endpoint, err := utils.ChooseVersion(client, versions)
	if err != nil {
		return "", err
	}

	switch chosen.ID {
	case v2:
		return getV2ProjectId(client, endpoint)
	case v3:
		return getV3ProjectId(client, endpoint)
	default:
		return "", fmt.Errorf("Unrecognized identity version: %s", chosen.ID)
	}
}

func getV2ProjectId(client *golangsdk.ProviderClient, endpoint string) (string, error) {
	v2Client, err := NewIdentityV2(client, golangsdk.EndpointOpts{})
	if err != nil {
		return "", err
	}

	if endpoint != "" {
		v2Client.Endpoint = endpoint
	}

	result := tokens2.Get(v2Client, client.TokenID)
	token, err := result.ExtractToken()
	if err != nil {
		return "", err
	}

	return token.Tenant.ID, nil
}

func getV3ProjectId(client *golangsdk.ProviderClient, endpoint string) (string, error) {
	v3Client, err := NewIdentityV3(client, golangsdk.EndpointOpts{})
	if err != nil {
		return "", err
	}

	if endpoint != "" {
		v3Client.Endpoint = endpoint
	}

	result := tokens3.Get(v3Client, client.TokenID)
	project, err := result.ExtractProject()
	if err != nil {
		return "", err
	}

	return project.ID, nil
}

func initClientOptsExtension(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts, clientType string) (*golangsdk.ServiceClientExtension, error) {
	pid, e := GetProjectId(client)
	if e != nil {
		return nil, e
	}

	c, e := initClientOpts(client, eo, clientType)
	if e != nil {
		return nil, e
	}

	sc := new(golangsdk.ServiceClientExtension)
	sc.ServiceClient = c
	sc.ProjectID = pid
	return sc, nil
}

//NewAutoScalingService creates a ServiceClient that may be used to access the
//auto-scaling service of huawei public cloud
func NewAutoScalingService(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "as")
	return sc, err
}

// NewKmsKeyV1 creates a ServiceClient that may be used to access the v1
// kms key service.
func NewKmsKeyV1(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "kms")
	return sc, err
}

func NewElasticLoadBalancer(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	//sc, err := initClientOpts1(client, eo, "elb")
	sc, err := initClientOpts(client, eo, "compute")
	if err != nil {
		return sc, err
	}
	sc.Endpoint = strings.Replace(sc.Endpoint, "ecs", "elb", 1)
	sc.Endpoint = sc.Endpoint[:strings.LastIndex(sc.Endpoint, "v2")+3]
	sc.Endpoint = strings.Replace(sc.Endpoint, "v2", "v1.0", 1)
	sc.ResourceBase = sc.Endpoint
	return sc, err
}

// NewNatV2 creates a ServiceClient that may be used with the v2 nat package.
func NewNatV2(client *golangsdk.ProviderClient, eo golangsdk.EndpointOpts) (*golangsdk.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "network")
	sc.Endpoint = strings.Replace(sc.Endpoint, "vpc", "nat", 1)
	sc.Endpoint = strings.Replace(sc.Endpoint, "myhwclouds", "myhuaweicloud", 1)
	sc.ResourceBase = sc.Endpoint + "v2.0/"
	return sc, err
}
