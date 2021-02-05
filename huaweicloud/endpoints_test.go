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

func testAccPreCheckServiceEndpoints(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set, skipping HuaweiCloud service endpoints test.")
	}

	projectID := os.Getenv("HW_PROJECT_ID")
	if projectID == "" {
		t.Fatalf(yellow("HW_PROJECT_ID must be set for service endpoint acceptance test"))
	}
}

func TestAccServiceEndpoints_Global(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider().(*schema.Provider)
	raw := make(map[string]interface{})
	err := testProvider.Configure(terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", err)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	config := testProvider.Meta().(*Config)

	// test the endpoint of IAM service
	serviceClient, err = config.IAMV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud IAM client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://iam.%s/v3.0/", config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("IAM endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("IAM endpoint:\t %s", actualURL)

	// test the endpoint of identity service
	serviceClient, err = config.IdentityV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud identity client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://iam.%s/v3/", config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("Identity endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("Identity endpoint:\t %s", actualURL)

	// test the endpoint of CDN service
	serviceClient, err = config.CdnV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud CDN client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cdn.%s/v1.0/", config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("CDN endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("CDN endpoint:\t %s", actualURL)

	// test the endpoint of bss v1 service
	serviceClient, err = config.BssV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud BSS v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://bss.%s/v1.0/", config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("BSS v1 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("BSS v1 endpoint:\t %s", actualURL)

	// test the endpoint of bss v2 service
	serviceClient, err = config.BssV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud BSS v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://bss.%s/v2/", config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("BSS v2 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("BSS v2 endpoint:\t %s", actualURL)
}

func TestAccServiceEndpoints_Management(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider().(*schema.Provider)
	raw := make(map[string]interface{})
	err := testProvider.Configure(terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", err)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	config := testProvider.Meta().(*Config)

	// test the endpoint of CTS service
	serviceClient, err = config.ctsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud CTS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cts.%s.%s/v1.0/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("CTS endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("CTS endpoint:\t %s", actualURL)

	// test the endpoint of LTS service
	serviceClient, err = config.ltsV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud LTS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://lts.%s.%s/v2/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("LTS endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("LTS endpoint:\t %s", actualURL)

	// test the endpoint of CES service
	serviceClient, err = config.newCESClient(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud CES client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://ces.%s.%s/V1.0/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("CES endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("CES endpoint:\t %s", actualURL)
}

func TestAccServiceEndpoints_Database(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider().(*schema.Provider)
	raw := make(map[string]interface{})
	err := testProvider.Configure(terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", err)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	config := testProvider.Meta().(*Config)

	// test the endpoint of RDS v1 service
	serviceClient, err = config.RdsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud RDS v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://rds.%s.%s/rds/v1/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("RDS v1 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("RDS v1 endpoint:\t %s", actualURL)

	// test the endpoint of RDS v3 service
	serviceClient, err = config.RdsV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud RDS v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://rds.%s.%s/v3/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("RDS v3 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("RDS v3 endpoint:\t %s", actualURL)

	// test the endpoint of DDS v3 service
	serviceClient, err = config.ddsV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud DDS v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dds.%s.%s/v3/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DDS v3 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DDS v3 endpoint:\t %s", actualURL)

	// test the endpoint of GeminiDB service
	serviceClient, err = config.GeminiDBV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud GeminiDB client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://gaussdb-nosql.%s.%s/v3/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("GeminiDB endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("GeminiDB/Cassandra endpoint:\t %s", actualURL)

	// test the endpoint of gaussdb service
	serviceClient, err = config.gaussdbV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud gaussdb client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://gaussdb.%s.%s/mysql/v3/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("gaussdb endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("gaussdb endpoint:\t %s", actualURL)

	// test the endpoint of openGauss service
	serviceClient, err = config.openGaussV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud openGauss client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://gaussdb.%s.%s/opengauss/v3/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("openGauss endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("openGauss endpoint:\t %s", actualURL)
}

func TestAccServiceEndpoints_Security(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider().(*schema.Provider)
	raw := make(map[string]interface{})
	err := testProvider.Configure(terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", err)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	config := testProvider.Meta().(*Config)

	// test the endpoint of anti-ddos service
	serviceClient, err = config.antiddosV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud anti-ddos client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://antiddos.%s.%s/v1/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("anti-ddos endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("anti-ddos endpoint:\t %s", actualURL)

	// test the endpoint of KMS service
	serviceClient, err = config.kmsKeyV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud KMS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://kms.%s.%s/v1.0/", HW_REGION_NAME, config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("KMS endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("KMS endpoint:\t %s", actualURL)
}

func TestAccServiceEndpoints_Application(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider().(*schema.Provider)
	raw := make(map[string]interface{})
	err := testProvider.Configure(terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", err)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	config := testProvider.Meta().(*Config)

	// test the endpoint of API-GW service
	serviceClient, err = config.apiGatewayV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud API-GW client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://apig.%s.%s/v1.0/apigw/", HW_REGION_NAME, config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("API-GW endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("API-GW endpoint:\t %s", actualURL)

	// test the endpoint of DCS v1 service
	serviceClient, err = config.dcsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud dcs v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dcs.%s.%s/v1.0/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DCS v1 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DCS v1 endpoint:\t %s", actualURL)

	// test the endpoint of DCS v2 service
	serviceClient, err = config.dcsV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud dcs v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dcs.%s.%s/v2/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DCS v2 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DCS v2 endpoint:\t %s", actualURL)

	// test the endpoint of DMS service
	serviceClient, err = config.dmsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud DMS v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dms.%s.%s/v1.0/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DMS v1 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DMS v1 endpoint:\t %s", actualURL)

	// test the endpoint of DMS v2 service
	serviceClient, err = config.dmsV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud DMS v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dms.%s.%s/v2/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DMS v2 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DMS v2 endpoint:\t %s", actualURL)
}

// TestAccServiceEndpoints_Compute test for endpoints of the clients used in ecs
// include computeV1Client,computeV11Client,computeV2Client,autoscalingV1Client,imageV2Client,
// cceV3Client,cceAddonV3Client,cciV1Client and FgsV2Client
func TestAccServiceEndpoints_Compute(t *testing.T) {

	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider().(*schema.Provider)
	raw := make(map[string]interface{})
	err := testProvider.Configure(terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", err)
	}

	config := testProvider.Meta().(*Config)
	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient

	// test for computeV1Client
	serviceClient, err = config.ComputeV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud ecs v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://ecs.%s.%s/v1/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "ecs", "v1", t)

	// test for computeV11Client
	serviceClient, err = nil, nil
	serviceClient, err = config.ComputeV11Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud ecs v1.1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://ecs.%s.%s/v1.1/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "ecs", "v1.1", t)

	// test for computeV2Client
	serviceClient, err = nil, nil
	serviceClient, err = config.ComputeV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud ecs v2.1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://ecs.%s.%s/v2.1/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "ecs", "v2.1", t)

	// test for autoscalingV1Client
	serviceClient, err = nil, nil
	serviceClient, err = config.AutoscalingV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud autoscaling v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://as.%s.%s/autoscaling-api/v1/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "autoscaling", "v1", t)

	// test for imageV2Client
	serviceClient, err = nil, nil
	serviceClient, err = config.ImageV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud image v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://ims.%s.%s/v2/", HW_REGION_NAME, config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "image", "v2", t)

	// test for cceV3Client
	serviceClient, err = nil, nil
	serviceClient, err = config.CceV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud cce v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cce.%s.%s/api/v3/projects/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "cce", "v3", t)

	// test for cceAddonV3Client
	serviceClient, err = nil, nil
	serviceClient, err = config.CceAddonV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud cceAddon v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cce.%s.%s/api/v3/", HW_REGION_NAME, config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "cceAddon", "v3", t)

	// test for cciV1Client
	serviceClient, err = nil, nil
	serviceClient, err = config.CciV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud cci v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cci.%s.%s/apis/networking.cci.io/v1beta1/", HW_REGION_NAME, config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "cci", "v1", t)

	// test for FgsV2Client
	serviceClient, err = nil, nil
	serviceClient, err = config.FgsV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud fgs v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://functiongraph.%s.%s/v2/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "fgs", "v2", t)
}

// TestAccServiceEndpoints_Storage test for the endpoints of the clients used in storage
// include blockStorageV2Client,blockStorageV3Client,sfsV2Client
// sfsV1Client,csbsV1Client and vbsV2Client
func TestAccServiceEndpoints_Storage(t *testing.T) {

	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider().(*schema.Provider)
	raw := make(map[string]interface{})
	err := testProvider.Configure(terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", err)
	}

	config := testProvider.Meta().(*Config)
	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient

	// test for blockStorageV2Client
	serviceClient, err = nil, nil
	serviceClient, err = config.BlockStorageV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud blockStorage v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://evs.%s.%s/v2/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "blockStorage", "v2", t)

	// test for blockStorageV3Client
	serviceClient, err = nil, nil
	serviceClient, err = config.BlockStorageV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud blockStorage v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://evs.%s.%s/v3/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "blockStorage", "v3", t)

	// test for	sfsV2Client
	serviceClient, err = nil, nil
	serviceClient, err = config.SfsV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud sfsV2 v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://sfs.%s.%s/v2/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "sfsV2", "v2", t)

	// test for sfsV1Client
	serviceClient, err = nil, nil
	serviceClient, err = config.SfsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud sfsV1 v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://sfs-turbo.%s.%s/v1/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "sfsV1", "v1", t)

	// test for csbsV1Client
	serviceClient, err = nil, nil
	serviceClient, err = config.CsbsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud csbsV1 v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://csbs.%s.%s/v1/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "csbsV1", "v1", t)

	// test for vbsV2Client
	serviceClient, err = nil, nil
	serviceClient, err = config.VbsV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud vbsV2 v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://vbs.%s.%s/v2/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "vbsV2", "v2", t)
}

// TestAccServiceEndpoints_Network test for the endpoints of the clients used in network
func TestAccServiceEndpoints_Network(t *testing.T) {

	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider().(*schema.Provider)
	raw := make(map[string]interface{})
	err := testProvider.Configure(terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", err)
	}

	config := testProvider.Meta().(*Config)
	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient

	// test endpoint of network v1 service
	serviceClient, err = config.NetworkingV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud networking v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://vpc.%s.%s/v1/", HW_REGION_NAME, config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "vpc", "v1", t)

	// test endpoint of network v2 service
	serviceClient, err = nil, nil
	serviceClient, err = config.NetworkingV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud networking v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://vpc.%s.%s/v2.0/", HW_REGION_NAME, config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "networking", "v2.0", t)

	// test endpoint of nat gateway
	serviceClient, err = nil, nil
	serviceClient, err = config.NatGatewayClient(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud nat gateway client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://nat.%s.%s/v2/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "nat", "v2", t)

	// test endpoint of secgroup v1
	serviceClient, err = nil, nil
	serviceClient, err = config.SecurityGroupV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud security_group v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://vpc.%s.%s/v1/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "security group", "v1", t)

	// test endpoint of elb v1.0
	serviceClient, err = nil, nil
	serviceClient, err = config.elasticLBClient(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud ELB v1.0 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://elb.%s.%s/v1.0/", HW_REGION_NAME, config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "elb", "v1.0", t)

	// test endpoint of elb v2.0
	serviceClient, err = nil, nil
	serviceClient, err = config.ElbV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud ELB v2.0 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://elb.%s.%s/v2.0/", HW_REGION_NAME, config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "elb", "v2.0", t)

	// test the endpoint of fw v2 service
	serviceClient, err = config.FwV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud fw v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://vpc.%s.%s/v2.0/", HW_REGION_NAME, config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "fw", "v2.0", t)

	// test the endpoint of DNS service
	serviceClient, err = config.DnsV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud DNS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dns.%s/v2/", config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DNS endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DNS endpoint:\t %s", actualURL)

	// test the endpoint of DNS service (with region)
	serviceClient, err = config.DnsWithRegionClient(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud DNS region client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dns.%s.%s/v2/", HW_REGION_NAME, config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DNS region endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DNS region endpoint:\t %s", actualURL)

	// test the endpoint of VPC endpoint
	serviceClient, err = config.VPCEPClient(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud VPC endpoint client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://vpcep.%s.%s/v1/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("VPCEP endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("VPCEP endpoint:\t %s", actualURL)
}

func TestAccServiceEndpoints_EnterpriseIntelligence(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider().(*schema.Provider)
	raw := make(map[string]interface{})
	err := testProvider.Configure(terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", err)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	config := testProvider.Meta().(*Config)

	// test the endpoint of MRS service
	serviceClient, err = config.MrsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud MRS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://mrs.%s.%s/v1.1/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("MRS endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("MRS endpoint:\t %s", actualURL)

	// test the endpoint of SMN service
	serviceClient, err = config.SmnV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud SMN client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://smn.%s.%s/v2/%s/notifications/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("SMN endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("SMN endpoint:\t %s", actualURL)

	serviceClient, err = config.cdmV11Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud cdm client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cdm.%s.%s/v1.1/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("cdm endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("cdm endpoint:\t %s", actualURL)

	serviceClient, err = config.disV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud dis client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dis.%s.%s/v2/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("dis endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("dis endpoint:\t %s", actualURL)

	serviceClient, err = config.cloudtableV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud cloudtable client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cloudtable.%s.%s/v2/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("cloudtable endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("cloudtable endpoint:\t %s", actualURL)

	serviceClient, err = config.cloudStreamV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud cloudStream client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cs.%s.%s/v1.0/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("cloudStream endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("cloudStream endpoint:\t %s", actualURL)

	serviceClient, err = config.cssV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud css client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://css.%s.%s/v1.0/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("css endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("css endpoint:\t %s", actualURL)

	serviceClient, err = config.dliV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating dli css client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dli.%s.%s/v1.0/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("dli endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("dli endpoint:\t %s", actualURL)

	serviceClient, err = config.dwsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating dws css client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dws.%s.%s/v1.0/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("dws endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("dws endpoint:\t %s", actualURL)

	serviceClient, err = config.gesV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating ges css client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://ges.%s.%s/v1.0/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("ges endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("ges endpoint:\t %s", actualURL)

	serviceClient, err = config.mlsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating mls css client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://mls.%s.%s/v1.0/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("mls endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("mls endpoint:\t %s", actualURL)
}

func TestAccServiceEndpoints_Edge(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider().(*schema.Provider)
	raw := make(map[string]interface{})
	err := testProvider.Configure(terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", err)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	config := testProvider.Meta().(*Config)

	// test the endpoint of iec service
	serviceClient, err = config.IECV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud IEC client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://iecs.%s/v1/", config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("IEC endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("IEC endpoint:\t %s", actualURL)
}

func TestAccServiceEndpoints_Others(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider().(*schema.Provider)
	raw := make(map[string]interface{})
	err := testProvider.Configure(terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", err)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	config := testProvider.Meta().(*Config)

	// test the endpoint of MAAS service
	serviceClient, err = config.maasV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud MAAS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://oms.%s.%s/v1/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("MAAS endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("MAAS endpoint:\t %s", actualURL)

	// test the endpoint of RTS service
	serviceClient, err = config.orchestrationV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud RTS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://rts.%s.%s/v1/%s/", HW_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("RTS endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("RTS endpoint:\t %s", actualURL)
}

func compareURL(expectedURL, actualURL, client, version string, t *testing.T) {
	if actualURL != expectedURL {
		t.Fatalf("%s %s endpoint: expected %s but got %s", client, version, green(expectedURL), yellow(actualURL))
	}
	t.Logf("%s %s endpoint:\t %s", client, version, actualURL)
}
