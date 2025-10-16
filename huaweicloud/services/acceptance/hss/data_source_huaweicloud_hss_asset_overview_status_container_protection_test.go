package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssetOverviewStatusContainerProtection_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_asset_overview_status_container_protection.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAssetOverviewStatusContainerProtection_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "no_risk"),
					resource.TestCheckResourceAttrSet(dataSourceName, "risk"),
					resource.TestCheckResourceAttrSet(dataSourceName, "no_protect"),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
				),
			},
		},
	})
}

func testAccDataSourceAssetOverviewStatusContainerProtection_basic() string {
	return `
data "huaweicloud_hss_asset_overview_status_container_protection" "test" {
  enterprise_project_id = "0"
}
`
}
