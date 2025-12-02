package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssetOverviewStatisticsTopPorts_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_asset_overview_statistics_top_ports.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case need a host with host protection enabled.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAssetOverviewStatisticsTopPorts_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.percentage"),
				),
			},
		},
	})
}

const testAccDataSourceAssetOverviewStatisticsTopPorts_basic = `
data "huaweicloud_hss_asset_overview_statistics_top_ports" "test" {
  enterprise_project_id = "0"
}
`
