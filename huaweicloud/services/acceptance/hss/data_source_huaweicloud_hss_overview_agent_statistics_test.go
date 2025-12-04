package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOverviewAgentStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_overview_agent_statistics.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOverviewAgentStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "wait_upgrade_num"),
					resource.TestCheckResourceAttrSet(dataSource, "online_num"),
					resource.TestCheckResourceAttrSet(dataSource, "not_online_num"),
					resource.TestCheckResourceAttrSet(dataSource, "offline_num"),
					resource.TestCheckResourceAttrSet(dataSource, "incluster_num"),
					resource.TestCheckResourceAttrSet(dataSource, "not_incluster_num"),
				),
			},
		},
	})
}

func testDataSourceOverviewAgentStatistics_basic() string {
	return `
data "huaweicloud_hss_overview_agent_statistics" "test" {
  enterprise_project_id = "0"
  container_type        = 1
}
`
}
