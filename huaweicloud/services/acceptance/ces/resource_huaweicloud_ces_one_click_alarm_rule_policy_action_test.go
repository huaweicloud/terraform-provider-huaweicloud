package ces

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCesOneClickAlarmRulePolicyAction_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCesOneClickAlarmRulePolicyAction_basic(name),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testCesOneClickAlarmRulePolicyAction_basic(name string) string {
	return fmt.Sprintf(`
%s

locals {
  alarm_id  = [for v in data.huaweicloud_ces_one_click_alarm_rules.test.alarms[*].alarm_id : v if v != ""][0]
  policy_id = [for v in data.huaweicloud_ces_one_click_alarm_rules.test.alarms[0].policies[*].alarm_policy_id :
    v if v != ""][0]
}

resource "huaweicloud_ces_one_click_alarm_rule_policy_action" "test" {
  one_click_alarm_id = huaweicloud_ces_one_click_alarm.test.id
  alarm_id           = local.alarm_id
  alarm_policy_ids   = [local.policy_id]
  enabled            = false
}
`, testDataSourceCesOneClickAlarmRules_basic(name))
}
