package acceptance

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud"
)

var (
	HW_REGION_NAME       = os.Getenv("HW_REGION_NAME")
	HW_AVAILABILITY_ZONE = os.Getenv("HW_AVAILABILITY_ZONE")
	HW_ACCESS_KEY        = os.Getenv("HW_ACCESS_KEY")
	HW_SECRET_KEY        = os.Getenv("HW_SECRET_KEY")
	HW_PROJECT_ID        = os.Getenv("HW_PROJECT_ID")
	HW_DOMAIN_ID         = os.Getenv("HW_DOMAIN_ID")
	HW_DOMAIN_NAME       = os.Getenv("HW_DOMAIN_NAME")

	HW_FLAVOR_ID   = os.Getenv("HW_FLAVOR_ID")
	HW_FLAVOR_NAME = os.Getenv("HW_FLAVOR_NAME")
	HW_IMAGE_ID    = os.Getenv("HW_IMAGE_ID")
	HW_IMAGE_NAME  = os.Getenv("HW_IMAGE_NAME")
	HW_VPC_ID      = os.Getenv("HW_VPC_ID")
	HW_NETWORK_ID  = os.Getenv("HW_NETWORK_ID")
	HW_SUBNET_ID   = os.Getenv("HW_SUBNET_ID")

	HW_DEPRECATED_ENVIRONMENT = os.Getenv("HW_DEPRECATED_ENVIRONMENT")
)

var TestAccProviders map[string]terraform.ResourceProvider
var TestAccProvider *schema.Provider

func init() {
	TestAccProvider = huaweicloud.Provider().(*schema.Provider)
	TestAccProviders = map[string]terraform.ResourceProvider{
		"huaweicloud": TestAccProvider,
	}
}

func TestAccPreCheckRequiredEnvVars(t *testing.T) {
	if HW_REGION_NAME == "" {
		t.Fatal("HW_REGION_NAME must be set for acceptance tests")
	}
}

func TestAccPreCheck(t *testing.T) {
	// Do not run the test if this is a deprecated testing environment.
	if HW_DEPRECATED_ENVIRONMENT != "" {
		t.Skip("This environment only runs deprecated tests")
	}

	TestAccPreCheckRequiredEnvVars(t)
}

func TestAccPreCheckDeprecated(t *testing.T) {
	if HW_DEPRECATED_ENVIRONMENT == "" {
		t.Skip("This environment does not support deprecated tests")
	}

	TestAccPreCheckRequiredEnvVars(t)
}
