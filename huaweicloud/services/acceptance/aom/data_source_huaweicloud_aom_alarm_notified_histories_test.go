package aom

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAlarmNotifiedHistories_basic(t *testing.T) {
	resourceName := "data.huaweicloud_aom_alarm_notified_histories.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAomAlarmEventSn(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAlarmNotifiedHistories_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestMatchResourceAttr(resourceName, "notified_histories.#", regexp.MustCompile("^[1-9]([0-9]+)?$")),
					resource.TestCheckResourceAttr(resourceName, "notified_histories.0.event_sn", acceptance.HW_AOM_ALARM_EVENT_SN),
					resource.TestMatchResourceAttr(resourceName, "notified_histories.0.notifications.#", regexp.MustCompile("^[1-9]([0-9]+)?$")),
					resource.TestCheckResourceAttrSet(resourceName, "notified_histories.0.notifications.0.action_rule"),
					resource.TestCheckResourceAttrSet(resourceName, "notified_histories.0.notifications.0.notifier_channel"),
					resource.TestMatchResourceAttr(resourceName, "notified_histories.0.notifications.0.smn_channel.#",
						regexp.MustCompile("^[1-9]([0-9]+)?$")),
					resource.TestCheckResourceAttrSet(resourceName, "notified_histories.0.notifications.0.smn_channel.0.sent_time"),
					resource.TestCheckResourceAttrSet(resourceName, "notified_histories.0.notifications.0.smn_channel.0.smn_request_id"),
					resource.TestCheckResourceAttrSet(resourceName, "notified_histories.0.notifications.0.smn_channel.0.smn_response_body"),
					resource.TestCheckResourceAttrSet(resourceName, "notified_histories.0.notifications.0.smn_channel.0.smn_response_code"),
					resource.TestCheckResourceAttrSet(resourceName, "notified_histories.0.notifications.0.smn_channel.0.smn_topic"),
					resource.TestMatchResourceAttr(resourceName, "notified_histories.0.notifications.0.smn_channel.0.smn_notified_history.#",
						regexp.MustCompile("^[1-9]([0-9]+)?$")),
					resource.TestCheckResourceAttrSet(resourceName,
						"notified_histories.0.notifications.0.smn_channel.0.smn_notified_history.0.smn_notified_content"),
					resource.TestCheckResourceAttrSet(resourceName,
						"notified_histories.0.notifications.0.smn_channel.0.smn_notified_history.0.smn_subscription_status"),
					resource.TestCheckResourceAttrSet(resourceName,
						"notified_histories.0.notifications.0.smn_channel.0.smn_notified_history.0.smn_subscription_type"),
				),
			},
		},
	})
}

func testAccDataAlarmNotifiedHistories_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_aom_alarm_notified_histories" "test" {
  event_sn = "%[1]s"
}
`, acceptance.HW_AOM_ALARM_EVENT_SN)
}
