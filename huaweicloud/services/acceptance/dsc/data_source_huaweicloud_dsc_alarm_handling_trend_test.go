package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscAlarmHandlingTrend_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dsc_alarm_handling_trend.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscAlarmHandlingTrend_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "time_axis.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_alarm_variation.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_event_variation.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "untreated_alarm_variation.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "untreated_event_variation.#"),
				),
			},
		},
	})
}

const testAccDataSourceDscAlarmHandlingTrend_basic = `
data "huaweicloud_dsc_alarm_handling_trend" "test" {}
`
