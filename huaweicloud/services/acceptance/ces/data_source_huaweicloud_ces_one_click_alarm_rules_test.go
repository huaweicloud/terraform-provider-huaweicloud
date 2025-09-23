package ces

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCesOneClickAlarmRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ces_one_click_alarm_rules.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCesOneClickAlarmRules_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.alarm_id"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.policies.#"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.policies.0.alarm_policy_id"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.policies.0.metric_name"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.notification_enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.alarm_notifications.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.alarm_notifications.0.notification_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.notification_begin_time"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.notification_end_time"),
				),
			},
		},
	})
}

func testDataSourceCesOneClickAlarmRules_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_ces_one_click_alarm_rules" "test" {
  one_click_alarm_id = huaweicloud_ces_one_click_alarm.test.id
}
`, testDataSourceCesOneClickAlarmRules_base(name))
}

func testDataSourceCesOneClickAlarmRules_base(name string) string {
	return testOneClickAlarm_basic(name)
}
