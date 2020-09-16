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
	t.Logf("GeminiDB endpoint:\t %s", actualURL)

	// test the endpoint of gaussdb service
	serviceClient, err = config.initServiceClient("gaussdb", OS_REGION_NAME, "mysql/v3")
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
	serviceClient, err = config.initServiceClient("dcs", OS_REGION_NAME, "v2")
	if err != nil {
		t.Fatalf("Error creating HuaweiCloud dcs v2 client: %s", err)
	}
	expectedURL = fmt.Sprintf("https://dcs.%s.%s/v2/%s/", OS_REGION_NAME, config.Cloud, config.TenantID)
	actualURL = serviceClient.ResourceBaseURL()
	if actualURL != expectedURL {
		t.Fatalf("DCS v2 endpoint: expected %s but got %s", green(expectedURL), yellow(actualURL))
	}
	t.Logf("DCS v2 endpoint:\t %s", actualURL)
}
