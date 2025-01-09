package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAggregatorPolicyAssignments_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_rms_resource_aggregator_policy_assignments.basic"
	dataSource2 := "data.huaweicloud_rms_resource_aggregator_policy_assignments.filter_by_compliance_state"
	dataSource3 := "data.huaweicloud_rms_resource_aggregator_policy_assignments.filter_by_policy_assignment_name"
	rName := acceptance.RandomAccResourceName()
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)
	dc3 := acceptance.InitDataSourceCheck(dataSource3)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAggregatorPolicyAssignments_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					dc3.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_compliance_state_filter_useful", "true"),
					resource.TestCheckOutput("is_policy_assignment_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceAggregatorPolicyAssignments_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rms_resource_aggregator" "test" {
  name        = "%[1]s"
  type        = "ACCOUNT"
  account_ids = ["%[2]s"]

  # wait 30 seconds to let the policies evaluate
  provisioner "local-exec" {
    command = "sleep 30"
  }
}
`, name, acceptance.HW_DOMAIN_ID)
}

func testDataSourceAggregatorPolicyAssignments_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_resource_aggregator_policy_assignments" "basic" {
  aggregator_id = huaweicloud_rms_resource_aggregator.test.id
}

data "huaweicloud_rms_resource_aggregator_policy_assignments" "filter_by_compliance_state" {
  aggregator_id = huaweicloud_rms_resource_aggregator.test.id
  
  filter {
    compliance_state = "Compliant"
  }
}

data "huaweicloud_rms_resource_aggregator_policy_assignments" "filter_by_policy_assignment_name" {
  aggregator_id = huaweicloud_rms_resource_aggregator.test.id

  filter {
    policy_assignment_name = "iam-password-policy"
  }
}

locals {
  compliance_state_filter_result = [
    for v in data.huaweicloud_rms_resource_aggregator_policy_assignments.filter_by_compliance_state.assignments[*].compliance.0.compliance_state :
    v == "Compliant"
  ]
  policy_assignment_name_filter_result = [
    for v in data.huaweicloud_rms_resource_aggregator_policy_assignments.filter_by_policy_assignment_name.assignments[*].policy_assignment_name :
    v == "iam-password-policy"
  ]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_rms_resource_aggregator_policy_assignments.basic.assignments) > 0
}

output "is_compliance_state_filter_useful" {
  value = alltrue(local.compliance_state_filter_result) && length(local.compliance_state_filter_result) > 0
}

output "is_policy_assignment_name_filter_useful" {
  value = alltrue(local.policy_assignment_name_filter_result) && length(local.policy_assignment_name_filter_result) > 0
}
`, testDataSourceAggregatorPolicyAssignments_base(name), name)
}
