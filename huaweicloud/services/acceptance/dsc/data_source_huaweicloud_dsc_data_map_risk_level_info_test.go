package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDataMapRiskLevelInfo_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dsc_data_map_risk_level_info.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataMapRiskLevelInfo_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_level.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_level.0.total"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_level.0.risk_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_level.0.risk_list.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_level.0.risk_list.0.detail_list.#"),
				),
			},
		},
	})
}

const testAccDataSourceDataMapRiskLevelInfo_basic = `
data "huaweicloud_dsc_data_map_risk_level_info" "test" {}
`
