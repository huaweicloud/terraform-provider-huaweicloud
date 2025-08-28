package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccPolicyAssignmentEvaluateResultUpdate_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRMSPolicyAssignmentEvaluationHash(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyAssignmentEvaluateResultUpdate__basic(),
			},
		},
	})
}

func testAccPolicyAssignmentEvaluateResultUpdate__basic() string {
	return fmt.Sprintf(
		`
data "huaweicloud_rms_resources" "test" {
  type = "kms.keys"
}

data "huaweicloud_rms_resource_policy_states" "test" {
  resource_id = data.huaweicloud_rms_resources.test.resources[0].id
}

resource "huaweicloud_rms_policy_assignment_evaluate_result_update" "test" {
  policy_assignment_id = data.huaweicloud_rms_resource_policy_states.test.value.0.policy_assignment_id
  trigger_type         = data.huaweicloud_rms_resource_policy_states.test.value.0.trigger_type
  compliance_state     = data.huaweicloud_rms_resource_policy_states.test.value.0.compliance_state
  evaluation_time      = data.huaweicloud_rms_resource_policy_states.test.value.0.evaluation_time
  evaluation_hash      = "%[1]s"

  policy_resource {
    resource_id       = data.huaweicloud_rms_resource_policy_states.test.value.0.resource_id
    resource_name     = data.huaweicloud_rms_resource_policy_states.test.value.0.resource_name
    resource_provider = data.huaweicloud_rms_resource_policy_states.test.value.0.resource_provider
    resource_type     = data.huaweicloud_rms_resource_policy_states.test.value.0.resource_type
    region_id         = data.huaweicloud_rms_resource_policy_states.test.value.0.region_id
    domain_id         = data.huaweicloud_rms_resource_policy_states.test.value.0.domain_id
  }
}
`, acceptance.HW_RMS_POLICY_ASSIGNMENT_EVALUATION_HASH)
}
