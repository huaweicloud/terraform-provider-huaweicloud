package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssetOverviewStatusOs_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_asset_overview_status_os.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAssetOverviewStatusOs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "win_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "linux_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "os_list.#"),
				),
			},
		},
	})
}

func testAccDataSourceAssetOverviewStatusOs_basic() string {
	return `
data "huaweicloud_hss_asset_overview_status_os" "test" {}
`
}
