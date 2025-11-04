package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceClusterProtectOverview_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_cluster_protect_overview.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceClusterProtectOverview_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "under_protect_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policy_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "event_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "block_num"),
				),
			},
		},
	})
}

func testDataSourceClusterProtectOverview_basic() string {
	return `
data "huaweicloud_hss_cluster_protect_overview" "test" {
  enterprise_project_id = "0"
}
`
}
