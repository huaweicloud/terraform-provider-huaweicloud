package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceRemediationExecution_basic(t *testing.T) {
	basicConfig := testResourceRmsRemediationExecution_base()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRMSTargetIDForFGS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testResourceRemediationExecution_basic(basicConfig),
			},
		},
	})
}

func TestAccResourceRemediationExecution_specifyResources(t *testing.T) {
	basicConfig := testResourceRmsRemediationExecution_base()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRMSTargetIDForFGS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testResourceRemediationExecution_specifyResources(basicConfig),
			},
		},
	})
}

func testResourceRemediationExecution_basic(basicConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rms_remediation_execution" "test" {
  policy_assignment_id = huaweicloud_rms_policy_assignment.test.id   
  all_supported        = true

  depends_on = [huaweicloud_rms_remediation_configuration.test]
}
`, basicConfig)
}

func testResourceRemediationExecution_specifyResources(basicConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rms_remediation_execution" "test" {
  policy_assignment_id = huaweicloud_rms_policy_assignment.test.id   
  all_supported        = false

  resource_ids = [
    huaweicloud_vpc.test[0].id,
    huaweicloud_vpc.test[1].id,
  ]   

  depends_on = [huaweicloud_rms_remediation_configuration.test]
}
`, basicConfig)
}

func testResourceRmsRemediationExecution_base() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  count       = 3
  name        = "vpc-bo-%[1]s${count.index}"
  cidr        = "192.168.0.0/16"
  description = "created by acc test"
  
  tags = {
    name = "bo"
  }

  lifecycle {
    ignore_changes = [
      name
    ]
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
    command = "sleep 15"
  }
}`, name, acceptance.HW_REGION_NAME, acceptance.HW_RMS_TARGET_ID_FOR_FGS)
}
