package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccModifyAlarmNotification_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running the test case, please ensure that there is at least one WAF instance in the current region.
			// Prepare a WAF alarm notification.
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckWafAlertId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccModifyAlarmNotification_basic(name),
			},
		},
	})
}

func testAccModifyAlarmNotification_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name                  = "%[2]s"
  enterprise_project_id = "0"
}

resource "huaweicloud_waf_modify_alarm_notification" "test" {
  alert_id                  = "%[1]s"
  name                      = "%[2]s"
  topic_urn                 = huaweicloud_smn_topic.test.topic_urn
  notice_class              = "threat_alert_notice"
  enabled                   = true
  sendfreq                  = 120
  locale                    = "en-us"
  times                     = 5
  threat                    = ["anticrawler","cc"]
  is_all_enterprise_project =  false
}
`, acceptance.HW_WAF_ALERT_ID, name)
}
