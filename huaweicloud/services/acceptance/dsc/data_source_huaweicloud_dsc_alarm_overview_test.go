package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscAlarmOverview_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dsc_alarm_overview.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscAlarmOverview_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "turn_off_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "turn_on_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "alarm_source_info.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_alarm.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "untreated_alarm.#"),
				),
			},
		},
	})
}

const testAccDataSourceDscAlarmOverview_basic = `
data "huaweicloud_dsc_alarm_overview" "test" {}
`
