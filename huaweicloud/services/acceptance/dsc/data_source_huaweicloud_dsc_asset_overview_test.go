package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscAssetOverview_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dsc_asset_overview.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscAssetOverview_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "asset_sensitive_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "asset_total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "high_level_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "middle_level_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "low_level_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "un_classed_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "asset_classification_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "asset_classification_list.0.color_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "asset_classification_list.0.level_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "asset_classification_list.0.sensitive_num"),
				),
			},
		},
	})
}

const testAccDataSourceDscAssetOverview_basic = `
data "huaweicloud_dsc_asset_overview" "test" {}
`
