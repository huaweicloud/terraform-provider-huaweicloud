package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBatchDeleteAlarmNotifications_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running the test case, please ensure that there is at least one WAF instance in the current region.
			// Prepare a WAF policy with a WAF CC protection rule.
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPrecheckWafAlarmNotificationId(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchDeleteAlarmNotifications_basic(),
			},
		},
	})
}

func testAccBatchDeleteAlarmNotifications_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_batch_delete_alarm_notifications" "test" {
  enterprise_project_id = "%[1]s"

  alert_notice_configs {
    id = "%[2]s"
  }
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, acceptance.HW_WAF_ALARM_NOTIFICATION_ID)
}
