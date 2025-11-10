package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceImageAssetStatistics_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_image_asset_statistics.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceImageAssetStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "local_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cicd_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "registry_num"),
				),
			},
		},
	})
}

func testAccDataSourceImageAssetStatistics_basic() string {
	return `
data "huaweicloud_hss_image_asset_statistics" "test" {
  enterprise_project_id = "0"
}
`
}
