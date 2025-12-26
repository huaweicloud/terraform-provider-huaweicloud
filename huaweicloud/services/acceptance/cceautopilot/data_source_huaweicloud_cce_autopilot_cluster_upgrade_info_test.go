package cceautopilot

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCceAutopilotClusterUpgradeInfo_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_autopilot_cluster_upgrade_info.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCceAutopilotClusterUpgradeInfo_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
				),
			},
		},
	})
}

const testDataSourceCceAutopilotClusterUpgradeInfo_basic = `
data "huaweicloud_cce_autopilot_cluster_upgrade_info" "test" {
  cluster_id = "83a085fe-d4f1-11f0-80d0-0255ac10178d"
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_cce_autopilot_cluster_upgrade_info.test.spec) > 0
}
`
