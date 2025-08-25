package rms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePolicyStatesSummary_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rms_policy_states_summary.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePolicyStatesSummary_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "results.#"),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.resource_details.#"),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.resource_details.0.compliant_count"),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.resource_details.0.non_compliant_count"),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.assignment_details.#"),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.assignment_details.0.compliant_count"),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.assignment_details.0.non_compliant_count"),
					resource.TestCheckOutput("tags_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourcePolicyStatesSummary_basic() string {
	return `
data "huaweicloud_rms_policy_states_summary" "test" {}

data "huaweicloud_rms_policy_assignments" "test" {}

locals {
  tag_key = [for v in data.huaweicloud_rms_policy_assignments.test.assignments : v.policy_filter if
  length(v.policy_filter) > 0 && alltrue([for vv in v.policy_filter : vv.tag_key != ""])][0][0].tag_key
}

data "huaweicloud_rms_policy_states_summary" "tags_filter" {
  tags =[local.tag_key]
}
output "tags_filter_is_useful" {
  value = length(data.huaweicloud_rms_policy_states_summary.tags_filter.results) > 0 
}
`
}
