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

	projectID := os.Getenv("OS_PROJECT_ID")
	if projectID == "" {
		t.Fatalf(yellow("OS_PROJECT_ID must be set for service endpoint acceptance test"))
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
	serviceClient, err = config.IdentityV3Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud IAM client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://iam.%s/v3/", config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("IAM endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("IAM endpoint:\t %s", actualURL)

	// test the endpoint of CDN service
	serviceClient, err = config.CdnV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud CDN client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cdn.%s/v1.0/", config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("CDN endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("CDN endpoint:\t %s", actualURL)

	// test the endpoint of DNS service
	serviceClient, err = config.DnsV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud DNS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dns.%s/v2/", config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DNS endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DNS endpoint:\t %s", actualURL)
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
	serviceClient, err = config.ctsV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud CTS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cts.%s.%s/v1.0/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("CTS endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("CTS endpoint:\t %s", actualURL)

	// test the endpoint of LTS service
	serviceClient, err = config.ltsV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud LTS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://lts.%s.%s/v2/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("LTS endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("LTS endpoint:\t %s", actualURL)

	// test the endpoint of CES service
	serviceClient, err = config.loadCESClient(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud CES client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://ces.%s.%s/V1.0/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
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
	serviceClient, err = config.RdsV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud RDS v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://rds.%s.%s/rds/v1/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("RDS v1 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("RDS v1 endpoint:\t %s", actualURL)

	// test the endpoint of RDS v3 service
	serviceClient, err = config.RdsV3Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud RDS v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://rds.%s.%s/v3/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("RDS v3 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("RDS v3 endpoint:\t %s", actualURL)

	// test the endpoint of DDS v3 service
	serviceClient, err = config.ddsV3Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud DDS v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dds.%s.%s/v3/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DDS v3 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DDS v3 endpoint:\t %s", actualURL)

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
	serviceClient, err = config.gaussdbV3Client(OS_REGION_NAME)
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
	serviceClient, err = config.antiddosV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud anti-ddos client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://antiddos.%s.%s/v1/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("anti-ddos endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("anti-ddos endpoint:\t %s", actualURL)

	// test the endpoint of KMS service
	serviceClient, err = config.kmsKeyV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud KMS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://kms.%s.%s/v1.0/", OS_REGION_NAME, config.Cloud)
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
	serviceClient, err = config.apiGatewayV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud API-GW client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://apig.%s.%s/v1.0/apigw/", OS_REGION_NAME, config.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("API-GW endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("API-GW endpoint:\t %s", actualURL)

	// test the endpoint of DCS v1 service
	serviceClient, err = config.dcsV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud dcs v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dcs.%s.%s/v1.0/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DCS v1 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DCS v1 endpoint:\t %s", actualURL)

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

	// test the endpoint of DMS service
	serviceClient, err = config.dmsV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud DMS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dms.%s.%s/v1.0/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DMS endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DMS endpoint:\t %s", actualURL)
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
	serviceClient, err = config.MrsV1Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud MRS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://mrs.%s.%s/v1.1/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("MRS endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("MRS endpoint:\t %s", actualURL)

	// test the endpoint of SMN service
	serviceClient, err = config.SmnV2Client(OS_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud SMN client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://smn.%s.%s/v2/%s/notifications/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("SMN endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("SMN endpoint:\t %s", actualURL)
}
