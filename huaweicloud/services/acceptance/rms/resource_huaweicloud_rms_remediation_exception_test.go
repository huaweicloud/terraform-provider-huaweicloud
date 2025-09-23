package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rms"
)

func getResourceRemediationExceptionFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("rms", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating RMS client: %s", err)
	}
	policyAssignmentID := state.Primary.ID
	return rms.ListRemediationExceptionInfo(client, cfg.DomainID, policyAssignmentID)
}

func TestAccResourceRemediationException_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_rms_remediation_exception.test"
	name := acceptance.RandomAccResourceName()
	baseConfig := testResourceRemediationException_base(name)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceRemediationExceptionFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRMSTargetIDForFGS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceRemediationException_basic(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "policy_assignment_id", "huaweicloud_rms_policy_assignment.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "exceptions.#", "1"),
				),
			},
			{
				Config: testResourceRemediationException_update(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "policy_assignment_id", "huaweicloud_rms_policy_assignment.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "exceptions.#", "2"),
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

func testResourceRemediationException_basic(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rms_remediation_exception" "test" {
  policy_assignment_id = huaweicloud_rms_policy_assignment.test.id

  exceptions {
    resource_id = huaweicloud_vpc.test[0].id
    message     = "test 0"
  }

  depends_on = [huaweicloud_rms_remediation_configuration.test]
}
`, baseConfig)
}

func testResourceRemediationException_update(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rms_remediation_exception" "test" {
  policy_assignment_id = huaweicloud_rms_policy_assignment.test.id

  exceptions {
    resource_id = huaweicloud_vpc.test[1].id
    message     = "test 1"
  }

  exceptions {
    resource_id = huaweicloud_vpc.test[2].id
    message     = "test 2"
  }

  depends_on = [huaweicloud_rms_remediation_configuration.test]
}
`, baseConfig)
}

func testResourceRemediationException_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  count       = 3
  name        = "vpc-bo-%[1]s${count.index}"
  cidr        = "192.168.0.0/16"
  description = "created by acc test"
  
  tags = {
    name = "bo"
  }
}

data "huaweicloud_rms_policy_definitions" "test" {
  name = "regular-matching-of-names"
}
  
resource "huaweicloud_rms_policy_assignment" "test" {
  name                 = "%[1]s"
  description          = "A resource name that does not match the regular expression is considered 'non-compliant'."
  policy_definition_id = try(data.huaweicloud_rms_policy_definitions.test.definitions[0].id, "")
  status               = "Enabled"
  
  policy_filter {
    region            = "%[2]s"
    resource_provider = "vpc"
    resource_type     = "vpcs"
    tag_key           = "name"
    tag_value         = "bo"
  }
  
  parameters = {
    regularExpression = jsonencode("bo-form_")
  }

  depends_on = [huaweicloud_vpc.test]

  lifecycle {
    ignore_changes = [
      status
    ]
  }
}

resource "huaweicloud_rms_remediation_configuration" "test" {
  policy_assignment_id = huaweicloud_rms_policy_assignment.test.id
  target_type          = "fgs"
  target_id            = "%[3]s"

  resource_parameter {
    resource_id = "file_prefix"
  }
	
  maximum_attempts      = 5  
  retry_attempt_seconds = 60

  provisioner "local-exec" {
    command = "sleep 10"
  }
}`, name, acceptance.HW_REGION_NAME, acceptance.HW_RMS_TARGET_ID_FOR_FGS)
}
