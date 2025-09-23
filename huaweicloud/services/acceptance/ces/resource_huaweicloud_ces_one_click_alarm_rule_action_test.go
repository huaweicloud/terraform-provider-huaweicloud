package ces

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCesOneClickAlarmRuleAction_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCesOneClickAlarmRuleAction_basic(name),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testCesOneClickAlarmRuleAction_basic(name string) string {
	return fmt.Sprintf(`
%s

locals {
  alarm_id = [for v in data.huaweicloud_ces_one_click_alarm_rules.test.alarms[*].alarm_id : v if v != ""][0]
}

resource "huaweicloud_ces_one_click_alarm_rule_action" "test" {
  one_click_alarm_id = huaweicloud_ces_one_click_alarm.test.id
  alarm_ids          = [local.alarm_id]
  alarm_enabled      = false
}
`, testDataSourceCesOneClickAlarmRules_basic(name))
}
