package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssetOverviewAccountTop_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_asset_overview_account_top.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAssetOverviewAccountTop_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "total_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.percentage"),
				),
			},
		},
	})
}

func testDataSourceAssetOverviewAccountTop_basic() string {
	return `
data "huaweicloud_hss_asset_overview_account_top" "test" {
  enterprise_project_id = "0"
}
`
}
