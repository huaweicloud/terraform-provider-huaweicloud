package ces

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCesAlarmRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ces_alarmrules.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCesAlarmRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.alarm_id"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.alarm_name"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.condition.0.metric_name"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.condition.0.comparison_operator"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.condition.0.period"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.condition.0.filter"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.condition.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.condition.0.suppress_duration"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.condition.0.alarm_level"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.alarm_type"),

					resource.TestCheckOutput("is_default_filter_useful", "true"),
					resource.TestCheckOutput("is_filter_by_id_useful", "true"),
					resource.TestCheckOutput("is_filter_by_name_useful", "true"),
					resource.TestCheckOutput("is_filter_by_namespace_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCesAlarmRules_basic() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
%[1]s

locals {
  id        = huaweicloud_ces_alarmrule.test.id
  name      = huaweicloud_ces_alarmrule.test.alarm_name
  namespace = huaweicloud_ces_alarmrule.test.metric.0.namespace
}

data "huaweicloud_ces_alarmrules" "test" {
  depends_on = [huaweicloud_ces_alarmrule.test]
}

output "is_default_filter_useful" {
  value = length(data.huaweicloud_ces_alarmrules.test.alarms) > 0
}

data "huaweicloud_ces_alarmrules" "filter_by_id" {
  alarm_id = local.id

  depends_on = [huaweicloud_ces_alarmrule.test]
}

output "is_filter_by_id_useful" {
  value = length(data.huaweicloud_ces_alarmrules.filter_by_id.alarms) > 0 && alltrue(
    [for alarm in data.huaweicloud_ces_alarmrules.filter_by_id.alarms[*] : alarm.alarm_id == local.id]
  )
}

data "huaweicloud_ces_alarmrules" "filter_by_name" {
  name = local.name

  depends_on = [huaweicloud_ces_alarmrule.test]
}

output "is_filter_by_name_useful" {
  value = length(data.huaweicloud_ces_alarmrules.filter_by_name.alarms) > 0 && alltrue(
    [for alarm in data.huaweicloud_ces_alarmrules.filter_by_name.alarms[*] : alarm.alarm_name == local.name]
  )
}

data "huaweicloud_ces_alarmrules" "filter_by_namespace" {
  namespace = local.namespace

  depends_on = [huaweicloud_ces_alarmrule.test]
}

output "is_filter_by_namespace_useful" {
  value = length(data.huaweicloud_ces_alarmrules.filter_by_namespace.alarms) > 0 && alltrue(
    [for alarm in data.huaweicloud_ces_alarmrules.filter_by_namespace.alarms[*] : alarm.namespace == local.namespace]
  )
}
`, testDataSourceCesAlarmRule_base(name))
}

func testDataSourceCesAlarmRule_base(name string) string {
	return testCESAlarmRule_basic(name)
}
