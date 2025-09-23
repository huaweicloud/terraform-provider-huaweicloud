package rms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsResourcePolicyStatesSummary_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rms_resource_policy_states_summary.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRmsResourcePolicyStatesSummary_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "value.#"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.compliance_state"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.results.#"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.results.0.resource_details.#"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.results.0.assignment_details.#"),
					resource.TestCheckOutput("resource_id_filter_is_useful", "true"),
					resource.TestCheckOutput("resource_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRmsResourcePolicyStatesSummary_basic() string {
	return `
data "huaweicloud_rms_resources" "test" {
  type = "kms.keys"
}

data "huaweicloud_rms_resource_policy_states_summary" "test" {}

data "huaweicloud_rms_resource_policy_states_summary" "resource_id_filter" {
  resource_id = data.huaweicloud_rms_resources.test.resources[0].id
}
locals {
  resource_id = data.huaweicloud_rms_resources.test.resources[0].id
}
output "resource_id_filter_is_useful" {
  value = length(data.huaweicloud_rms_resource_policy_states_summary.resource_id_filter.value) > 0 
}

data "huaweicloud_rms_resource_policy_states_summary" "resource_name_filter" {
  resource_name = data.huaweicloud_rms_resources.test.resources[0].name
}
locals {
  resource_name = data.huaweicloud_rms_resources.test.resources[0].name
}
output "resource_name_filter_is_useful" {
  value = length(data.huaweicloud_rms_resource_policy_states_summary.resource_name_filter.value) > 0 
}
`
}
