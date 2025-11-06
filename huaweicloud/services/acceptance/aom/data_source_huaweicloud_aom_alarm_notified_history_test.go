package aom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAlarmNotifiedHistory_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_aom_alarm_notified_history.test"
	timeRange := "-1.-1.60" // Last 60 minutes

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAlarmNotifiedHistory_basic(timeRange),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "time_range", timeRange),
					resource.TestCheckResourceAttrSet(dataSourceName, "events.#"),
				),
			},
		},
	})
}

func testAccDataSourceAlarmNotifiedHistory_basic(timeRange string) string {
	return fmt.Sprintf(`
data "huaweicloud_aom_alarm_notified_history" "test" {
  time_range = "%s"
}
`, timeRange)
}
