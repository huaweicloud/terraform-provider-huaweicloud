package coc

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceAlarmClear_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocAlarmID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAlarmClear_basic(),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testAlarmClear_basic() string {
	currentTime := time.Now()
	tenMinutesAgo := currentTime.Add(-10*time.Minute).Unix() * 1e3
	fiveMinutesAgo := currentTime.Add(-5*time.Minute).Unix() * 1e3
	return fmt.Sprintf(`
resource "huaweicloud_coc_alarm_clear" "test" {
  alarm_ids            = "%[1]s"
  remarks              = "this is remark"
  is_service_interrupt = false
  start_time           = %[2]v
  fault_recovery_time  = %[3]v
}
`, acceptance.HW_COC_ALARM_ID, fiveMinutesAgo, tenMinutesAgo)
}
