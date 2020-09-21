package huaweicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
)

const (
	greenCode  = "\033[0m\033[1;32m"
	yellowCode = "\033[0m\033[1;33m"
	resetCode  = "\033[0m\033[1;31m"
)

func green(str interface{}) string {
	return fmt.Sprintf("%s%#v%s", greenCode, str, resetCode)
}

func yellow(str interface{}) string {
	return fmt.Sprintf("%s%#v%s", yellowCode, str, resetCode)
}

func TestAccServiceEndpoints(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set, skipping HuaweiCloud SSL test.")
	}

	projectID := os.Getenv("OS_PROJECT_ID")
	if projectID == "" {
		t.Fatalf(yellow("OS_PROJECT_ID must be set for service endpoint acceptance test"))
	}

	testProvider := Provider().(*schema.Provider)
	raw := make(map[string]interface{})
	err := testProvider.Configure(terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", err)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	config := testProvider.Meta().(*Config)

	// test the endpoint of RDS v3 service
	serviceClient, err = config.RdsV3Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud RDS v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://rds.%s.%s/v3/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("rds v3 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("rds v3 endpoint:\t %s", actualURL)

	// test the endpoint of sfs-turbo service
	serviceClient, err = config.sfsV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud sfs-turbo client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://sfs-turbo.%s.%s/v1/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("sfs-turbo endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("sfs-turbo endpoint:\t %s", actualURL)

	// test the endpoint of GeminiDB service
	serviceClient, err = config.GeminiDBV3Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud GeminiDB client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://gaussdb-nosql.%s.%s/v3/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("GeminiDB endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("GeminiDB/Cassandra endpoint:\t %s", actualURL)

	// test the endpoint of gaussdb service
	serviceClient, err = config.NewServiceClient("gaussdb", OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud gaussdb client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://gaussdb.%s.%s/mysql/v3/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("gaussdb endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("gaussdb endpoint:\t %s", actualURL)

	// test the endpoint of openGauss service
	serviceClient, err = config.openGaussV3Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud openGauss client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://gaussdb.%s.%s/opengauss/v3/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("openGauss endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("openGauss endpoint:\t %s", actualURL)

	// test the endpoint of DCS v2 service
	serviceClient, err = config.NewServiceClient("dcsv2", OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud dcs v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dcs.%s.%s/v2/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DCS v2 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DCS v2 endpoint:\t %s", actualURL)

	// test the endpoint of bss v1 service
	serviceClient, err = config.BssV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud bss v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://bss.%s.%s/v1.0/", OS_REGION_NAME, config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("Bss v1 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("Bss v1 endpoint:\t %s", actualURL)

	testEndpointOfCompute(config, t)

	testEndpointOfStorage(config, t)

	testEndpointOfNetWork(config, t)
}

// testEndpointOfCompute test for endpoints of the clients used in ecs
// include computeV1Client,computeV11Client,computeV2Client,autoscalingV1Client,imageV2Client,
// cceV3Client,cceAddonV3Client,cciV1Client and FgsV2Client
func testEndpointOfCompute(config *Config, t *testing.T) {

	var expectedURL, actualURL string
	var (
		serviceClient *golangsdk.ServiceClient
		err           error
	)
	// test for computeV1Client
	serviceClient, err = config.computeV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud ecs v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://ecs.%s.%s/v1/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "ecs", "v1", t)

	// test for computeV11Client
	serviceClient, err = nil, nil
	serviceClient, err = config.computeV11Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud ecs v1.1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://ecs.%s.%s/v1.1/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "ecs", "v1.1", t)

	// test for computeV2Client
	serviceClient, err = nil, nil
	serviceClient, err = config.computeV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud ecs v2.1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://ecs.%s.%s/v2.1/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "ecs", "v2.1", t)

	// test for autoscalingV1Client
	serviceClient, err = nil, nil
	serviceClient, err = config.autoscalingV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud autoscaling v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://as.%s.%s/autoscaling-api/v1/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "autoscaling", "v1", t)

	// test for imageV2Client
	serviceClient, err = nil, nil
	serviceClient, err = config.imageV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud image v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://ims.%s.%s/v2/", OS_REGION_NAME, config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "image", "v2", t)

	// test for cceV3Client
	serviceClient, err = nil, nil
	serviceClient, err = config.cceV3Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud cce v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cce.%s.%s/api/v3/projects/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "cce", "v3", t)

	// test for cceAddonV3Client
	serviceClient, err = nil, nil
	serviceClient, err = config.cceAddonV3Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud cceAddon v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cce.%s.%s/api/v3/", OS_REGION_NAME, config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "cceAddon", "v3", t)

	// test for cciV1Client
	serviceClient, err = nil, nil
	serviceClient, err = config.cciV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud cci v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cci.%s.%s/apis/networking.cci.io/v1beta1/", OS_REGION_NAME, config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "cci", "v1", t)

	// test for FgsV2Client
	serviceClient, err = nil, nil
	serviceClient, err = config.FgsV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud fgs v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://functiongraph.%s.%s/v2/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "fgs", "v2", t)
}

// testEndpointOfStorage test for the endpoints of the clients used in storage
// include blockStorageV2Client,blockStorageV3Client,loadEVSV2Client,sfsV2Client
// sfsV1Client,csbsV1Client and vbsV2Client
func testEndpointOfStorage(config *Config, t *testing.T) {

	var expectedURL, actualURL string
	var (
		serviceClient *golangsdk.ServiceClient
		err           error
	)

	// test for blockStorageV2Client
	serviceClient, err = nil, nil
	serviceClient, err = config.blockStorageV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud blockStorage v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://evs.%s.%s/v2/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "blockStorage", "v2", t)

	// test for blockStorageV3Client
	serviceClient, err = nil, nil
	serviceClient, err = config.blockStorageV3Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud blockStorage v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://evs.%s.%s/v3/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "blockStorage", "v3", t)

	// test for loadEVSV2Client
	serviceClient, err = nil, nil
	serviceClient, err = config.loadEVSV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud loadEVS v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://evs.%s.%s/v2/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "loadEVS", "v2", t)

	// test for	sfsV2Client
	serviceClient, err = nil, nil
	serviceClient, err = config.sfsV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud sfsV2 v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://sfs.%s.%s/v2/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "sfsV2", "v2", t)

	// test for sfsV1Client
	serviceClient, err = nil, nil
	serviceClient, err = config.sfsV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud sfsV1 v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://sfs-turbo.%s.%s/v1/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "sfsV1", "v1", t)

	// test for csbsV1Client
	serviceClient, err = nil, nil
	serviceClient, err = config.csbsV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud csbsV1 v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://csbs.%s.%s/v1/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "csbsV1", "v1", t)

	// test for vbsV2Client
	serviceClient, err = nil, nil
	serviceClient, err = config.vbsV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud vbsV2 v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://vbs.%s.%s/v2/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "vbsV2", "v2", t)

}

// testEndpointOfNetWork test for the endpoints of the clients used in network
// include networkingV1Client, networkingV2Client, networkingHwV2Client, natV2Client, loadElasticLoadBalancerClient and fwV2Client
func testEndpointOfNetWork(config *Config, t *testing.T) {

	var expectedURL, actualURL string
	var (
		serviceClient *golangsdk.ServiceClient
		err           error
	)

	// test endpoint of network v1 service
	serviceClient, err = config.networkingV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud networking v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://vpc.%s.%s/v1/", OS_REGION_NAME, config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "networking", "v1", t)

	// test endpoint of network v2 service
	serviceClient, err = nil, nil
	serviceClient, err = config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud networking v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://vpc.%s.%s/v2.0/", OS_REGION_NAME, config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "networking", "v2.0", t)

	// test endpoint of networkingHw v2
	serviceClient, err = nil, nil
	serviceClient, err = config.networkingHwV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud networkingHw v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://vpc.%s.%s/v2.0/", OS_REGION_NAME, config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "networkingHw", "v2.0", t)

	// test endpoint of nat v2
	serviceClient, err = nil, nil
	serviceClient, err = config.natV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud nat v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://nat.%s.%s/v2.0/", OS_REGION_NAME, config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "nat", "v2.0", t)

	// test endpoint of loadElasticLoadBalancer v1.0
	serviceClient, err = nil, nil
	serviceClient, err = config.loadElasticLoadBalancerClient(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud loadElasticLoadBalancer v1.0 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://elb.%s.%s/v1.0/", OS_REGION_NAME, config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "loadElasticLoadBalancer", "v1.0", t)

	// test the endpoint of fw v2 service
	serviceClient, err = config.fwV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud fw v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://vpc.%s.%s/v2.0/", OS_REGION_NAME, config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "fw", "v2.0", t)
}

func compareURL(expectedURL, actualURL, client, version string, t *testing.T) {

	if actualURL != expectedURL {
		t.Fatalf("%s %s endpoint: expected %s but got %s", client, version, green(expectedURL), yellow(actualURL))
	}
	t.Logf("%s %s endpoint:\t %s", client, version, actualURL)
}
