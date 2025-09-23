package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocAlarmActionHistories_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_alarm_action_histories.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocAlarmID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocAlarmActionHistories_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_handle_histories.#"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_handle_histories.0.work_order_id"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_handle_histories.0.create_name"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_handle_histories.0.task_type"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_handle_histories.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_handle_histories.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_handle_histories.0.duration"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_handle_histories.0.status"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocAlarmActionHistories_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_coc_alarm_action_histories" "test" {
  alarm_id = "%s"
}
`, acceptance.HW_COC_ALARM_ID)
}
