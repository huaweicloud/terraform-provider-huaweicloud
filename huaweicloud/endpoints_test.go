package huaweicloud

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
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

func TestAccServiceEndpoints_IAM(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", diags[0].Summary)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error
	cfg := testProvider.Meta().(*config.Config)

	// test the endpoint of IAM service
	serviceClient, err = cfg.IAMV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud IAM client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://iam.%s.%s/v3.0/", HW_REGION_NAME, cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("IAM endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("IAM endpoint:\t %s", actualURL)

	// test the endpoint of identity service
	serviceClient, err = cfg.IdentityV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud identity client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://iam.%s.%s/v3/", HW_REGION_NAME, cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("Identity endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("Identity endpoint:\t %s", actualURL)

	// test the endpoint of IAM service without version number
	serviceClient, err = cfg.IAMNoVersionClient(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud IAM client without version number: %s", err)
	}
	expectedURL = fmt.Sprintf("https://iam.%s.%s/", HW_REGION_NAME, cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("Identity endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("IAM endpoint without version number:\t %s", actualURL)
}

func TestAccServiceEndpoints_Global(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", diags[0].Summary)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error
	cfg := testProvider.Meta().(*config.Config)

	// test the endpoint of CDN service
	serviceClient, err = cfg.CdnV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud CDN client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cdn.%s/v1.0/", cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("CDN endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("CDN endpoint:\t %s", actualURL)

	// test the endpoint of bss v1 service
	serviceClient, err = cfg.BssV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud BSS v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://bss.%s/v1.0/", cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("BSS v1 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("BSS v1 endpoint:\t %s", actualURL)

	// test the endpoint of bss v2 service
	serviceClient, err = cfg.BssV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud BSS v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://bss.%s/v2/", cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("BSS v2 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("BSS v2 endpoint:\t %s", actualURL)

	// test the endpoint of EPS service
	serviceClient, err = cfg.EnterpriseProjectClient(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud EPS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://eps.%s/v1.0/", cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("EPS endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("EPS endpoint:\t %s", actualURL)

	// test the endpoint of RMS service
	serviceClient, err = cfg.RmsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud RMS v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://rms.%s/v1/", cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("RMS endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("RMS endpoint:\t %s", actualURL)
}

func TestAccServiceEndpoints_Management(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", diags[0].Summary)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error
	cfg := testProvider.Meta().(*config.Config)

	// test the endpoint of CTS service
	serviceClient, err = cfg.CtsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud CTS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cts.%s.%s/v1.0/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("CTS endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("CTS endpoint:\t %s", actualURL)

	// test the endpoint of LTS service
	serviceClient, err = cfg.LtsV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud LTS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://lts.%s.%s/v2/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("LTS endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("LTS endpoint:\t %s", actualURL)

	// test the endpoint of CES service
	serviceClient, err = cfg.CesV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud CES client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://ces.%s.%s/V1.0/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("CES endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("CES endpoint:\t %s", actualURL)

	// test the endpoint of TMS service
	serviceClient, err = cfg.TmsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("error creating TMS v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://tms.%s/v1.0/", cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("TMS endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("TMS endpoint:\t %s", actualURL)

	serviceClient, err = cfg.TmsV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("error creating TMS v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://tms.%s/v2.0/", cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("TMS endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("TMS endpoint:\t %s", actualURL)
}

func TestAccServiceEndpoints_Database(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", diags[0].Summary)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error
	cfg := testProvider.Meta().(*config.Config)

	// test the endpoint of RDS v1 service
	serviceClient, err = cfg.RdsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud RDS v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://rds.%s.%s/rds/v1/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("RDS v1 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("RDS v1 endpoint:\t %s", actualURL)

	// test the endpoint of RDS v3 service
	serviceClient, err = cfg.RdsV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud RDS v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://rds.%s.%s/v3/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("RDS v3 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("RDS v3 endpoint:\t %s", actualURL)

	// test the endpoint of DDS v3 service
	serviceClient, err = cfg.DdsV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud DDS v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dds.%s.%s/v3/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DDS v3 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DDS v3 endpoint:\t %s", actualURL)

	// test the endpoint of GeminiDB service
	serviceClient, err = cfg.GeminiDBV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud GeminiDB client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://gaussdb-nosql.%s.%s/v3/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("GeminiDB endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("GeminiDB/Cassandra endpoint:\t %s", actualURL)

	// test the endpoint of GeminiDB V3.1 service
	serviceClient, err = cfg.GeminiDBV31Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud GeminiDBV31 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://gaussdb-nosql.%s.%s/v3.1/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("GeminiDBV31 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("GeminiDB/CassandraV31 endpoint:\t %s", actualURL)

	// test the endpoint of gaussdb service
	serviceClient, err = cfg.GaussdbV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud gaussdb client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://gaussdb.%s.%s/v3/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("gaussdb endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("gaussdb endpoint:\t %s", actualURL)

	// test the endpoint of openGauss service
	serviceClient, err = cfg.OpenGaussV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud openGauss client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://gaussdb-opengauss.%s.%s/v3/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("openGauss endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("openGauss endpoint:\t %s", actualURL)

	// test the endpoint of DRS
	serviceClient, err = cfg.DrsV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud DRS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://drs.%s.%s/v3/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DRS endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DRS endpoint:\t %s", actualURL)
}

func TestAccServiceEndpoints_GaussDBForMySQL(t *testing.T) {
	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", diags[0].Summary)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error
	cfg := testProvider.Meta().(*config.Config)

	// test the endpoint of gaussdb service in cn-north-4
	region := "cn-north-4"
	serviceClient, err = cfg.GaussdbV3Client(region)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud gaussdb client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://gaussdb.%s.%s/", region, cfg.Cloud)
	actualURL = serviceClient.Endpoint
	if actualURL != expectedURL {
		t.Fatalf("gaussdb endpoint in %s: expected %s but got %s", region, green(expectedURL), yellow(actualURL))
	}
	t.Logf("gaussdb endpoint in %s:\t %s", region, actualURL)

	// test the endpoint of gaussdb service in ap-southeast-1
	region = "ap-southeast-1"
	serviceClient, err = cfg.GaussdbV3Client(region)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud gaussdb client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://gaussdbformysql.%s.%s/", region, cfg.Cloud)
	actualURL = serviceClient.Endpoint
	if actualURL != expectedURL {
		t.Fatalf("gaussdb endpoint in %s: expected %s but got %s", region, green(expectedURL), yellow(actualURL))
	}
	t.Logf("gaussdb endpoint in %s:\t %s", region, actualURL)
}

func TestAccServiceEndpoints_Security(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", diags[0].Summary)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error
	cfg := testProvider.Meta().(*config.Config)

	// test the endpoint of anti-ddos service
	serviceClient, err = cfg.AntiDDosV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud anti-ddos client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://antiddos.%s.%s/v1/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("anti-ddos endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("anti-ddos endpoint:\t %s", actualURL)

	// test the endpoint of KMS service v1.0
	serviceClient, err = cfg.KmsKeyV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud KMS(v1.0) client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://kms.%s.%s/v1.0/", HW_REGION_NAME, cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("KMS(v1.0) endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("KMS(v1.0) endpoint:\t %s", actualURL)

	// test the endpoint of KMS service v1
	serviceClient, err = cfg.KmsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud KMS(v1) client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://kms.%s.%s/v1/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("KMS(v1) endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("KMS(v1) endpoint:\t %s", actualURL)

	// test the endpoint of KMS service v3
	serviceClient, err = cfg.KmsV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud KMS(v3) client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://kms.%s.%s/v3/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("KMS(v3) endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("KMS(v3) endpoint:\t %s", actualURL)

	// test the endpoint of SCM service
	serviceClient, err = cfg.ScmV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud SCM client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://scm.%s.%s/v3/", HW_REGION_NAME, cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("SCM endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("SCM endpoint:\t %s", actualURL)

	// test the endpoint of WAF service
	serviceClient, err = cfg.WafV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud WAF client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://waf.%s.%s/v1/%s/waf/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("WAF endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("WAF endpoint:\t %s", actualURL)

	// test the endpoint of WAF Dedicated service
	serviceClient, err = cfg.WafDedicatedV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud WAF dedicated client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://waf.%s.%s/v1/%s/premium-waf/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("WAF dedicated endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("WAF dedicated endpoint:\t %s", actualURL)
}

func TestAccServiceEndpoints_Application(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", diags[0].Summary)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error
	cfg := testProvider.Meta().(*config.Config)

	// test the endpoint of API-GW service
	serviceClient, err = cfg.ApiGatewayV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud API-GW client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://apig.%s.%s/v1.0/apigw/", HW_REGION_NAME, cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("API-GW endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("API-GW endpoint:\t %s", actualURL)

	// test the endpoint of API-GW v2 service
	serviceClient, err = cfg.ApigV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud API-GW v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://apig.%s.%s/v2/%s/apigw/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("API-GW v2 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("API-GW v2 endpoint:\t %s", actualURL)

	// test the endpoint of BCS v2 service
	serviceClient, err = cfg.BcsV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud BCS v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://bcs.%s.%s/v2/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("BCS v2 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("BCS v2 endpoint:\t %s", actualURL)

	// test the endpoint of CSE v2 service
	serviceClient, err = cfg.CseV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud CSE v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cse.%s.%s/v2/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("CSE v2 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("CSE v2 endpoint:\t %s", actualURL)

	// test the endpoint of DCS v1 service
	serviceClient, err = cfg.DcsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud dcs v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dcs.%s.%s/v1.0/", HW_REGION_NAME, cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DCS v1 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DCS v1 endpoint:\t %s", actualURL)

	// test the endpoint of DCS v2 service
	serviceClient, err = cfg.DcsV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud dcs v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dcs.%s.%s/v2/", HW_REGION_NAME, cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DCS v2 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DCS v2 endpoint:\t %s", actualURL)

	// test the endpoint of DMS service
	serviceClient, err = cfg.DmsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud DMS v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dms.%s.%s/v1.0/", HW_REGION_NAME, cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DMS v1 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DMS v1 endpoint:\t %s", actualURL)

	// test the endpoint of DMS v2 service
	serviceClient, err = cfg.DmsV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud DMS v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dms.%s.%s/v2/", HW_REGION_NAME, cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DMS v2 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DMS v2 endpoint:\t %s", actualURL)

	// test the endpoint of ServiceStage V1 service
	serviceClient, err = cfg.ServiceStageV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud ServiceStage v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://servicestage.%s.%s/v1/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("ServiceStage v1 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("ServiceStage v1 endpoint:\t %s", actualURL)

	// test the endpoint of ServiceStage v2 service
	serviceClient, err = cfg.ServiceStageV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud ServiceStage v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://servicestage.%s.%s/v2/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("ServiceStage v2 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("ServiceStage v2 endpoint:\t %s", actualURL)

	// test the endpoint of EG service
	serviceClient, err = cfg.EgV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("error creating EG v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://eg.%s.%s/v1/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("EG endpoint: expected %s, but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("EG endpoint:\t %s", actualURL)
}

// TestAccServiceEndpoints_Compute test for endpoints of the clients used in ecs
// include computeV1Client,computeV11Client,computeV2Client,autoscalingV1Client,imageV2Client,
// cceV3Client,cceAddonV3Client,cciV1Client,cciV1BetaClient and FgsV2Client
func TestAccServiceEndpoints_Compute(t *testing.T) {

	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", diags[0].Summary)
	}

	cfg := testProvider.Meta().(*config.Config)
	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error

	// test for computeV1Client
	serviceClient, err = cfg.ComputeV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud ecs v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://ecs.%s.%s/v1/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "ecs", "v1", t)

	// test for computeV11Client
	serviceClient, err = cfg.ComputeV11Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud ecs v1.1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://ecs.%s.%s/v1.1/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "ecs", "v1.1", t)

	// test for computeV2Client
	serviceClient, err = cfg.ComputeV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud ecs v2.1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://ecs.%s.%s/v2.1/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "ecs", "v2.1", t)

	// test for autoscalingV1Client
	serviceClient, err = cfg.AutoscalingV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud autoscaling v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://as.%s.%s/autoscaling-api/v1/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "autoscaling", "v1", t)

	// test for imageV2Client
	serviceClient, err = cfg.ImageV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud image v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://ims.%s.%s/v2/", HW_REGION_NAME, cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "image", "v2", t)

	// test for cceV3Client
	serviceClient, err = cfg.CceV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud cce v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cce.%s.%s/api/v3/projects/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "cce", "v3", t)

	// test for cceAddonV3Client
	serviceClient, err = cfg.CceAddonV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud cceAddon v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cce.%s.%s/api/v3/", HW_REGION_NAME, cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "cceAddon", "v3", t)

	// test for cciV1Client
	serviceClient, err = cfg.CciV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud cci v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cci.%s.%s/api/v1/", HW_REGION_NAME, cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "cci", "v1", t)

	// test for cciV1BetaClient
	serviceClient, err = cfg.CciV1BetaClient(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud cci v1 beta client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cci.%s.%s/apis/networking.cci.io/v1beta1/", HW_REGION_NAME, cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "cci", "v1 beta", t)

	// test for FgsV2Client
	serviceClient, err = cfg.FgsV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud fgs v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://functiongraph.%s.%s/v2/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "fgs", "v2", t)

	// test for swrV2Client
	serviceClient, err = cfg.SwrV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud swr v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://swr-api.%s.%s/v2/", HW_REGION_NAME, cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "swr", "v2", t)

	// test for BmsV1Client
	serviceClient, err = cfg.BmsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud BMS v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://bms.%s.%s/v1/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "bms", "v1", t)
}

// TestAccServiceEndpoints_Storage test for the endpoints of the clients used in storage
// include blockStorageV2Client, BlockStorageV21Client, sfsV2Client, sfsV1Client,
// csbsV1Client, vbsV2Client and cbrV3Client
func TestAccServiceEndpoints_Storage(t *testing.T) {

	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", diags[0].Summary)
	}

	cfg := testProvider.Meta().(*config.Config)
	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error

	// test for blockStorageV2Client
	serviceClient, err = cfg.BlockStorageV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud blockStorage v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://evs.%s.%s/v2/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "blockStorage", "v2", t)

	// test for blockStorageV21Client
	serviceClient, err = cfg.BlockStorageV21Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud blockStorage v2.1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://evs.%s.%s/v2.1/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "blockStorage", "v2.1", t)

	// test for cbrV3Client
	serviceClient, err = cfg.CbrV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud cbr v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cbr.%s.%s/v3/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "cbr", "v3", t)

	// test for	sfsV2Client
	serviceClient, err = cfg.SfsV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud sfs v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://sfs.%s.%s/v2/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "sfs", "v2", t)

	// test for sfsV1Client
	serviceClient, err = cfg.SfsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud sfs v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://sfs-turbo.%s.%s/v1/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "sfs", "v1", t)

	// test for csbsV1Client
	serviceClient, err = cfg.CsbsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud csbs v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://csbs.%s.%s/v1/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "csbs", "v1", t)

	// test for vbsV2Client
	serviceClient, err = cfg.VbsV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud vbs v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://vbs.%s.%s/v2/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "vbs", "v2", t)
}

// TestAccServiceEndpoints_Network test for the endpoints of the clients used in network
func TestAccServiceEndpoints_Network(t *testing.T) {

	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", diags[0].Summary)
	}

	cfg := testProvider.Meta().(*config.Config)
	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error

	// test endpoint of vpc v1 service
	serviceClient, err = cfg.NetworkingV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud vpc v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://vpc.%s.%s/v1/", HW_REGION_NAME, cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "vpc", "v1", t)

	// test endpoint of vpc v3
	serviceClient, err = cfg.NetworkingV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud vpc v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://vpc.%s.%s/v3/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "vpc", "v3", t)

	// test endpoint of network v2 service
	serviceClient, err = cfg.NetworkingV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud networking v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://vpc.%s.%s/v2.0/", HW_REGION_NAME, cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "networking", "v2.0", t)

	// test endpoint of nat gateway
	serviceClient, err = cfg.NatGatewayClient(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud nat gateway client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://nat.%s.%s/v2/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "nat", "v2", t)

	// test endpoint of nat gateway v2.0
	serviceClient, err = cfg.NatV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud nat gateway v2.0 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://nat.%s.%s/v2.0/", HW_REGION_NAME, cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "nat", "v2.0", t)

	// test endpoint of elb v2.0
	serviceClient, err = cfg.ElbV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud ELB v2.0 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://elb.%s.%s/v2.0/", HW_REGION_NAME, cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "elb", "v2.0", t)

	// test endpoint of elb v3
	serviceClient, err = cfg.ElbV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud ELB v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://elb.%s.%s/v3/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "elb", "v3", t)

	// test endpoint of loadbalancer(elb v2)
	serviceClient, err = cfg.LoadBalancerClient(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud ELB v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://elb.%s.%s/v2/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "elb", "v2", t)

	// test the endpoint of fw v2 service
	serviceClient, err = cfg.FwV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud fw v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://vpc.%s.%s/v2.0/", HW_REGION_NAME, cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	compareURL(expectedURL, actualURL, "fw", "v2.0", t)

	// test the endpoint of DNS service
	serviceClient, err = cfg.DnsV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud DNS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dns.%s/v2/", cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DNS endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DNS endpoint:\t %s", actualURL)

	// test the endpoint of DNS service (with region)
	serviceClient, err = cfg.DnsWithRegionClient(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud DNS region client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dns.%s.%s/v2/", HW_REGION_NAME, cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DNS region endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DNS region endpoint:\t %s", actualURL)

	// test the endpoint of VPC endpoint
	serviceClient, err = cfg.VPCEPClient(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud VPC endpoint client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://vpcep.%s.%s/v1/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("VPCEP endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("VPCEP endpoint:\t %s", actualURL)

	// Test the endpoint of ER endpoint (ver.3)
	serviceClient, err = cfg.ErV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating ER v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://er.%s.%s/v3/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("ER endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("ER endpoint:\t %s", actualURL)

	// Test the endpoint of DC endpoint (ver.3)
	serviceClient, err = cfg.DcV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating DC v3 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dcaas.%s.%s/v3/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DC endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DC endpoint:\t %s", actualURL)
}

func TestAccServiceEndpoints_EnterpriseIntelligence(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", diags[0].Summary)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error
	cfg := testProvider.Meta().(*config.Config)

	// test the endpoint of MRS v1.1 service
	serviceClient, err = cfg.MrsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud MRS v1.1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://mrs.%s.%s/v1.1/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("MRS v1.1 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("MRS v1.1 endpoint:\t %s", actualURL)

	// test the endpoint of MRS v2 service
	serviceClient, err = cfg.MrsV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud MRS v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://mrs.%s.%s/v2/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("MRS v2 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("MRS v2 endpoint:\t %s", actualURL)

	// test the endpoint of SMN service
	serviceClient, err = cfg.SmnV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud SMN client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://smn.%s.%s/v2/%s/notifications/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("SMN endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("SMN endpoint:\t %s", actualURL)

	serviceClient, err = cfg.CdmV11Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud cdm client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cdm.%s.%s/v1.1/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("cdm endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("cdm endpoint:\t %s", actualURL)

	serviceClient, err = cfg.DisV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud dis client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dis.%s.%s/v2/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("dis endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("dis endpoint:\t %s", actualURL)

	serviceClient, err = cfg.DisV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud dis client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dis.%s.%s/v3/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("dis v3 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("dis v3 endpoint:\t %s", actualURL)

	serviceClient, err = cfg.CloudtableV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud cloudtable client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cloudtable.%s.%s/v2/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("cloudtable endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("cloudtable endpoint:\t %s", actualURL)

	serviceClient, err = cfg.CloudStreamV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud cloudStream client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://cs.%s.%s/v1.0/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("cloudStream endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("cloudStream endpoint:\t %s", actualURL)

	serviceClient, err = cfg.CssV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud css client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://css.%s.%s/v1.0/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("css endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("css endpoint:\t %s", actualURL)

	serviceClient, err = cfg.DliV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating dli css client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dli.%s.%s/v1.0/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("dli endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("dli endpoint:\t %s", actualURL)

	// test the endpoint of DLI v2.0 service
	serviceClient, err = cfg.DliV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating dli v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dli.%s.%s/v2.0/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("dli endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("dli endpoint:\t %s", actualURL)

	serviceClient, err = cfg.DwsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating dws css client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dws.%s.%s/v1.0/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("dws endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("dws endpoint:\t %s", actualURL)

	serviceClient, err = cfg.DwsV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating dws v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dws.%s.%s/v2/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("dws v2 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("dws v2 endpoint:\t %s", actualURL)

	serviceClient, err = cfg.GesV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating ges css client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://ges.%s.%s/v1.0/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("ges endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("ges endpoint:\t %s", actualURL)

	serviceClient, err = cfg.MlsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating mls css client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://mls.%s.%s/v1.0/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("mls endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("mls endpoint:\t %s", actualURL)

	serviceClient, err = cfg.ModelArtsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating ModelArts v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://modelarts.%s.%s/v1/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("ModelArts v1 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("ModelArts v1 endpoint:\t %s", actualURL)

	serviceClient, err = cfg.ModelArtsV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating ModelArts v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://modelarts.%s.%s/v2/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("ModelArts v2 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("ModelArts v2 endpoint:\t %s", actualURL)

	// test the endpoint of DataArts Studio v1 service
	serviceClient, err = cfg.DataArtsV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud DataArts v1 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dayu.%s.%s/v1/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DataArts v1 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DataArts v1 endpoint:\t %s", actualURL)

	// test the endpoint of Workspace service (with region)
	serviceClient, err = cfg.WorkspaceV2Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud Workspace V2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://workspace.%s.%s/v2/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("Workspace region endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("Workspace region endpoint:\t %s", actualURL)

}

func TestAccServiceEndpoints_Edge(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", diags[0].Summary)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error
	cfg := testProvider.Meta().(*config.Config)

	// test the endpoint of iec service
	serviceClient, err = cfg.IECV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud IEC client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://iecs.%s/v1/", cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("IEC endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("IEC endpoint:\t %s", actualURL)
}

func TestAccServiceEndpoints_Others(t *testing.T) {
	testAccPreCheckServiceEndpoints(t)

	testProvider := Provider()
	raw := make(map[string]interface{})
	diags := testProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("Unexpected error when configure HuaweiCloud provider: %s", diags[0].Summary)
	}

	var expectedURL, actualURL string
	var serviceClient *golangsdk.ServiceClient
	var err error
	cfg := testProvider.Meta().(*config.Config)

	// test the endpoint of MAAS service
	serviceClient, err = cfg.MaasV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud MAAS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://oms.%s.%s/v1/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("MAAS endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("MAAS endpoint:\t %s", actualURL)

	// test the endpoint of SMS service
	serviceClient, err = cfg.SmsV3Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud SMS client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://sms.ap-southeast-1.%s/v3/", cfg.Cloud)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("SMS endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("SMS endpoint:\t %s", actualURL)

	// test the endpoint of AOM service
	serviceClient, err = cfg.AomV1Client(HW_REGION_NAME)
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud AOM client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://aom.%s.%s/svcstg/icmgr/v1/%s/", HW_REGION_NAME, cfg.Cloud, cfg.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("AOM endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("AOM endpoint:\t %s", actualURL)
}

func compareURL(expectedURL, actualURL, client, version string, t *testing.T) {
	if actualURL != expectedURL {
		t.Fatalf("%s %s endpoint: expected %s but got %s", client, version, green(expectedURL), yellow(actualURL))
	}
	t.Logf("%s %s endpoint:\t %s", client, version, actualURL)
}
