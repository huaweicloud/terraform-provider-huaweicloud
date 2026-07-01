package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscEventOverview_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dsc_event_overview.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscEventOverview_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "block_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "turn_off_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "turn_on_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "not_overdue_event.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "not_overdue_event.0.fatal_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "not_overdue_event.0.high_risk_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "not_overdue_event.0.middle_risk_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "not_overdue_event.0.low_risk_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "not_overdue_event.0.notice_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "overdue_event.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "overdue_event.0.fatal_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "overdue_event.0.high_risk_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "overdue_event.0.middle_risk_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "overdue_event.0.low_risk_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "overdue_event.0.notice_num"),
				),
			},
		},
	})
}

func testAccDataSourceDscEventOverview_basic() string {
	return `
data "huaweicloud_dsc_event_overview" "test" {}
`
}
