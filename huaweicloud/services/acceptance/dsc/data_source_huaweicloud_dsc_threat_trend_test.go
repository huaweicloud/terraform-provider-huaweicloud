package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscThreatTrend_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dsc_threat_trend.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscThreatTrend_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "api_attacked_variation.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "database_attacked_variation.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "time_axis.#"),
				),
			},
		},
	})
}

const testAccDataSourceDscThreatTrend_basic = `
data "huaweicloud_dsc_threat_trend" "test" {}
`
