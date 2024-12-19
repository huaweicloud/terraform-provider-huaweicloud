package cceautopilot

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCceAutopilotAddonTemplates_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_cce_autopilot_addon_templates.basic"
	dataSource2 := "data.huaweicloud_cce_autopilot_addon_templates.filter_by_name"
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCceAutopilotAddonTemplates_basic,
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

const testDataSourceDataSourceCceAutopilotAddonTemplates_basic = `
data "huaweicloud_cce_autopilot_addon_templates" "basic" {}

data "huaweicloud_cce_autopilot_addon_templates" "filter_by_name" {
  addon_template_name = "log-agent"
}

locals {
  name_result = [for v in data.huaweicloud_cce_autopilot_addon_templates.filter_by_name.templates[*].name : v == "log-agent"]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_cce_autopilot_addon_templates.basic.templates) > 0
}

output "is_name_filter_useful" {
  value = alltrue(local.name_result) && length(local.name_result) > 0
}
`
