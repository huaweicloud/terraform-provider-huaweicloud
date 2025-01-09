package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAggregatorPolicyStates_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_rms_resource_aggregator_policy_states.basic"
	dataSource2 := "data.huaweicloud_rms_resource_aggregator_policy_states.filter_by_compliance_state"
	dataSource3 := "data.huaweicloud_rms_resource_aggregator_policy_states.filter_by_policy_assignment_name"
	dataSource4 := "data.huaweicloud_rms_resource_aggregator_policy_states.filter_by_resource_id"
	rName := acceptance.RandomAccResourceName()
	password := acceptance.RandomPassword()
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)
	dc3 := acceptance.InitDataSourceCheck(dataSource3)
	dc4 := acceptance.InitDataSourceCheck(dataSource4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAggregatorPolicyStates_basic(rName, password),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					dc3.CheckResourceExists(),
					dc4.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_compliance_state_filter_useful", "true"),
					resource.TestCheckOutput("is_policy_assignment_name_filter_useful", "true"),
					resource.TestCheckOutput("is_resource_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceAggregatorPolicyStates_base(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user" "test" {
  name     = "%[1]s"
  password = "%[2]s"
  enabled  = true
  email    = "%[1]s@abc.com"
}

resource "huaweicloud_rms_resource_aggregator" "test" {
  name        = "%[1]s"
  type        = "ACCOUNT"
  account_ids = ["%[3]s"]

  depends_on = [huaweicloud_identity_user.test]

  # wait 40 seconds to let the policies evaluate
  provisioner "local-exec" {
    command = "sleep 40"
  }
}
`, name, password, acceptance.HW_DOMAIN_ID)
}

func testDataSourceAggregatorPolicyStates_basic(name, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_resource_aggregator_policy_states" "basic" {
  aggregator_id = huaweicloud_rms_resource_aggregator.test.id
}

data "huaweicloud_rms_resource_aggregator_policy_states" "filter_by_compliance_state" {
  aggregator_id     = huaweicloud_rms_resource_aggregator.test.id
  compliance_state  = "Compliant"
}

data "huaweicloud_rms_resource_aggregator_policy_states" "filter_by_policy_assignment_name" {
  aggregator_id          = huaweicloud_rms_resource_aggregator.test.id
  policy_assignment_name = "iam-password-policy"
}

data "huaweicloud_rms_resource_aggregator_policy_states" "filter_by_resource_id" {
  aggregator_id = huaweicloud_rms_resource_aggregator.test.id
  resource_id   = huaweicloud_identity_user.test.id
}

locals {
  compliance_state_filter_result = [
    for v in data.huaweicloud_rms_resource_aggregator_policy_states.filter_by_compliance_state.states[*].compliance_state :
    v == "Compliant"
  ]
  policy_assignment_name_filter_result = [
    for v in data.huaweicloud_rms_resource_aggregator_policy_states.filter_by_policy_assignment_name.states[*].policy_assignment_name :
    v == "iam-password-policy"
  ]
  resource_id_filter_result = [
    for v in data.huaweicloud_rms_resource_aggregator_policy_states.filter_by_resource_id.states[*].resource_id :
    v == huaweicloud_identity_user.test.id
  ]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_rms_resource_aggregator_policy_states.basic.states) > 0
}

output "is_compliance_state_filter_useful" {
  value = alltrue(local.compliance_state_filter_result) && length(local.compliance_state_filter_result) > 0
}

output "is_policy_assignment_name_filter_useful" {
  value = alltrue(local.policy_assignment_name_filter_result) && length(local.policy_assignment_name_filter_result) > 0
}

output "is_resource_id_filter_useful" {
  value = alltrue(local.resource_id_filter_result) && length(local.resource_id_filter_result) > 0
}
`, testDataSourceAggregatorPolicyStates_base(name, password), name)
}
