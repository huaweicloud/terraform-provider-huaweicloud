package rms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsResourcePolicyStates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rms_resource_policy_states.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRmsResourcePolicyStates_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "value.#"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.resource_provider"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.trigger_type"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.compliance_state"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.policy_assignment_id"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.policy_assignment_name"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.policy_definition_id"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.evaluation_time"),
					resource.TestCheckOutput("compliance_state_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRmsResourcePolicyStates_basic() string {
	return `
data "huaweicloud_rms_resources" "test" {
  type = "kms.keys"
}

data "huaweicloud_rms_resource_policy_states" "test" {
  resource_id = data.huaweicloud_rms_resources.test.resources[0].id
}

locals {
  compliance_state = data.huaweicloud_rms_resource_policy_states.test.value[0].compliance_state
}
data "huaweicloud_rms_resource_policy_states" "compliance_state_filter" {
  resource_id      = data.huaweicloud_rms_resources.test.resources[0].id
  compliance_state = data.huaweicloud_rms_resource_policy_states.test.value[0].compliance_state
}
output "compliance_state_filter_is_useful" {
  value = length(data.huaweicloud_rms_resource_policy_states.compliance_state_filter.value) > 0 && alltrue(
  [for v in data.huaweicloud_rms_resource_policy_states.compliance_state_filter.value[*].compliance_state :
  v == local.compliance_state]
  )
}
`
}
