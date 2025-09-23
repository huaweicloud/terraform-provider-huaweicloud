package rms

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getResourceRmsRemediationConfigurationFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("rms", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating RMS client: %s", err)
	}

	uri := "v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/remediation-configuration"
	uri = strings.ReplaceAll(uri, "{domain_id}", cfg.DomainID)
	uri = strings.ReplaceAll(uri, "{policy_assignment_id}", state.Primary.ID)

	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Data()
}

func TestAccResourceRmsRemediationConfiguration_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_rms_remediation_configuration.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceRmsRemediationConfigurationFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRMSTargetIDForRFS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceRmsRemediationConfiguration_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "target_type", "rfs"),
					resource.TestCheckResourceAttr(resourceName, "target_id", acceptance.HW_RMS_TARGET_ID_FOR_RFS),
					resource.TestCheckResourceAttr(resourceName, "resource_parameter.0.resource_id", "file_prefix"),
					resource.TestCheckResourceAttr(resourceName, "static_parameter.0.var_key", "bucket_name"),
					resource.TestCheckResourceAttr(resourceName, "auth_type", "agency"),
					resource.TestCheckResourceAttr(resourceName, "auth_value", "test_RFS_CTS"),
					resource.TestCheckResourceAttr(resourceName, "maximum_attempts", "5"),
					resource.TestCheckResourceAttr(resourceName, "retry_attempt_seconds", "3600"),
					resource.TestCheckResourceAttrSet(resourceName, "static_parameter.0.var_value"),
				),
			},
			{
				Config: testResourceRmsRemediationConfiguration_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "target_type", "rfs"),
					resource.TestCheckResourceAttr(resourceName, "target_id", acceptance.HW_RMS_TARGET_ID_FOR_RFS),
					resource.TestCheckResourceAttr(resourceName, "resource_parameter.0.resource_id", "file_prefix"),
					resource.TestCheckResourceAttr(resourceName, "static_parameter.0.var_key", "bucket_name"),
					resource.TestCheckResourceAttr(resourceName, "static_parameter.1.var_key", "compress_type"),
					resource.TestCheckResourceAttr(resourceName, "static_parameter.1.var_value", "\"json\""),
					resource.TestCheckResourceAttr(resourceName, "auth_type", "agency"),
					resource.TestCheckResourceAttr(resourceName, "auth_value", "test_RFS_CTS"),
					resource.TestCheckResourceAttr(resourceName, "maximum_attempts", "6"),
					resource.TestCheckResourceAttr(resourceName, "retry_attempt_seconds", "60"),
					resource.TestCheckResourceAttrSet(resourceName, "static_parameter.0.var_value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testResourceRmsRemediationConfiguration_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rms_remediation_configuration" "test" {
  policy_assignment_id = huaweicloud_rms_policy_assignment.test.id
  target_type          = "rfs"
  target_id            = "%[2]s"
  resource_parameter {
    resource_id = "file_prefix"
  }
  
  static_parameter {
    var_key   = "bucket_name"
    var_value = "\"${huaweicloud_obs_bucket.test.bucket}\""
  }
  
  auth_type             = "agency" 
  auth_value            = "test_RFS_CTS" 
  maximum_attempts      = 5  
  retry_attempt_seconds = 3600  
}`, testResourceRmsRemediationConfiguration_base(), acceptance.HW_RMS_TARGET_ID_FOR_RFS)
}

func testResourceRmsRemediationConfiguration_update() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rms_remediation_configuration" "test" {
  policy_assignment_id = huaweicloud_rms_policy_assignment.test.id
  target_type          = "rfs"
  target_id            = "%[2]s"
  resource_parameter {
    resource_id = "file_prefix"
  }
  
  static_parameter {
    var_key   = "bucket_name"
    var_value = "\"${huaweicloud_obs_bucket.test.bucket}\""
  }

  static_parameter {
    var_key   = "compress_type"
    var_value = "\"json\""
  }
  
  auth_type             = "agency" 
  auth_value            = "test_RFS_CTS" 
  maximum_attempts      = 6  
  retry_attempt_seconds = 60  
}`, testResourceRmsRemediationConfiguration_base(), acceptance.HW_RMS_TARGET_ID_FOR_RFS)
}

func testResourceRmsRemediationConfiguration_base() string {
	name := acceptance.RandomAccResourceNameWithDash()
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[1]s"
  storage_class = "STANDARD"
  acl           = "private"
  force_destroy = true
}
  
resource "huaweicloud_cts_tracker" "test" {
  bucket_name = huaweicloud_obs_bucket.test.bucket
  file_prefix = "cts-updated"
  lts_enabled = false
  enabled     = false
}

data "huaweicloud_rms_policy_definitions" "test" {
  name = "multi-region-cts-tracker-exists"
}
  
resource "huaweicloud_rms_policy_assignment" "test" {
  name                 = "multi-region-cts-tracker-exists"
  description          = "An account is noncompliant if it does not have a CTS tracker in specified regions."
  period               = "One_Hour"
  policy_definition_id = try(data.huaweicloud_rms_policy_definitions.test.definitions[0].id, "")
  
  parameters = {
    regionList = jsonencode(["cn-north-4"])
  }
}`, name)
}
