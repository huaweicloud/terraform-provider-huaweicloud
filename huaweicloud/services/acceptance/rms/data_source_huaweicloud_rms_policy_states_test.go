package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsPolicyStates_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_rms_policy_states.basic"
	dataSource2 := "data.huaweicloud_rms_policy_states.filter_by_compliance_state"
	dataSource3 := "data.huaweicloud_rms_policy_states.filter_by_resource_name"
	dataSource4 := "data.huaweicloud_rms_policy_states.filter_by_resource_id"
	dataSource5 := "data.huaweicloud_rms_policy_states.filter_by_assignment_id"
	rName := acceptance.RandomAccResourceName()
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)
	dc3 := acceptance.InitDataSourceCheck(dataSource3)
	dc4 := acceptance.InitDataSourceCheck(dataSource4)
	dc5 := acceptance.InitDataSourceCheck(dataSource5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRmsPolicyStates_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					dc3.CheckResourceExists(),
					dc4.CheckResourceExists(),
					dc5.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_compliance_state_filter_useful", "true"),
					resource.TestCheckOutput("is_resource_name_filter_useful", "true"),
					resource.TestCheckOutput("is_resource_id_filter_useful", "true"),
					resource.TestCheckOutput("is_assignment_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRmsPolicyStates_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_policy_definitions" "test" {
  name = "allowed-ecs-flavors"
}

resource "huaweicloud_rms_policy_assignment" "test" {
  name                 = "%[2]s"
  description          = "An ECS is noncompliant if its flavor is not in the specified flavor list (filter by resource ID)."
  policy_definition_id = try(data.huaweicloud_rms_policy_definitions.test.definitions[0].id, "")

  policy_filter {
    region            = "%[3]s"
    resource_provider = "ecs"
    resource_type     = "cloudservers"
    resource_id       = huaweicloud_compute_instance.test.id
  }

  parameters = {
    listOfAllowedFlavors = "[\"${data.huaweicloud_compute_flavors.test.ids[0]}\"]"
  }
}
`, testAccPolicyAssignment_ecsConfig(name), name, acceptance.HW_REGION_NAME)
}

func testDataSourceDataSourceRmsPolicyStates_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_policy_states" "basic" {}

data "huaweicloud_rms_policy_states" "filter_by_compliance_state" {
  compliance_state = "Compliant"
}

data "huaweicloud_rms_policy_states" "filter_by_resource_name" {
  resource_name = "%[2]s"

  depends_on = [huaweicloud_compute_instance.test]
}

data "huaweicloud_rms_policy_states" "filter_by_resource_id" {
  resource_id = huaweicloud_compute_instance.test.id
}

data "huaweicloud_rms_policy_states" "filter_by_assignment_id" {
  policy_assignment_id = huaweicloud_rms_policy_assignment.test.id
}

locals {
  compliance_state_result = [for v in data.huaweicloud_rms_policy_states.filter_by_compliance_state.states[*].compliance_state : v == "Compliant"]

  resource_name_filter_result = [for v in data.huaweicloud_rms_policy_states.filter_by_resource_name.states[*].resource_name : v == "%[2]s"]

  resource_id_filter_result = [
    for v in data.huaweicloud_rms_policy_states.filter_by_resource_id.states[*].resource_id : v == huaweicloud_compute_instance.test.id
  ]

  assignment_id_filter_result = [
    for v in data.huaweicloud_rms_policy_states.filter_by_assignment_id.states[*].policy_assignment_id :
	v == huaweicloud_rms_policy_assignment.test.id
  ]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_rms_policy_states.basic.states) > 0
}

output "is_compliance_state_filter_useful" {
  value = alltrue(local.compliance_state_result) && length(local.compliance_state_result) > 0
}

output "is_resource_name_filter_useful" {
  value = alltrue(local.resource_name_filter_result) && length(local.resource_name_filter_result) > 0
}

output "is_resource_id_filter_useful" {
  value = alltrue(local.resource_id_filter_result) && length(local.resource_id_filter_result) > 0
}

output "is_assignment_id_filter_useful" {
  value = alltrue(local.assignment_id_filter_result) && length(local.assignment_id_filter_result) > 0
}

`, testDataSourceDataSourceRmsPolicyStates_base(name), name)
}
