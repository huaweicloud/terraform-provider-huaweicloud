package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssetOverviewStatusAgent_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_asset_overview_status_agent.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAssetOverviewStatusAgent_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "online"),
					resource.TestCheckResourceAttrSet(dataSourceName, "offline"),
					resource.TestCheckResourceAttrSet(dataSourceName, "not_installed"),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
				),
			},
		},
	})
}

func testAccDataSourceAssetOverviewStatusAgent_basic() string {
	return `
data "huaweicloud_hss_asset_overview_status_agent" "test" {
  enterprise_project_id = "0"
}
`
}
