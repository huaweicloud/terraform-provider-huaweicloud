package aom

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAlarmRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aom_alarm_rules.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAlarmRules_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_rules.#"),
				),
			},
		},
	})
}

func testDataSourceAlarmRules_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_aom_alarm_rules" "test" {
  depends_on = [huaweicloud_aomv4_alarm_rule.test]
}
`, testAlarmRuleV4_basic(name))
}

func TestAccDataSourceAlarmRules_event(t *testing.T) {
	dataSource := "data.huaweicloud_aom_alarm_rules.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAlarmRules_event(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_rules.#"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_rules.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_rules.0.alarm_rule_id"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_rules.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_rules.0.enable"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_rules.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_rules.0.alarm_notifications.#"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_rules.0.event_alarm_spec.#"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_rules.0.status"),
					resource.TestMatchResourceAttr(dataSource, "alarm_rules.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),

					resource.TestCheckOutput("name_validation", "true"),
					resource.TestCheckOutput("eps_id_validation", "true"),
					resource.TestCheckOutput("event_source_validation", "true"),
					resource.TestCheckOutput("type_validation", "true"),
					resource.TestCheckOutput("status_validation", "true"),
					resource.TestCheckOutput("bind_notification_rule_id_validation", "true"),
				),
			},
		},
	})
}

func testDataSourceAlarmRules_event(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_aom_alarm_rules" "test" {
  depends_on = [huaweicloud_aomv4_alarm_rule.test]

  alarm_rule_name           = huaweicloud_aomv4_alarm_rule.test.name
  enterprise_project_id     = huaweicloud_aomv4_alarm_rule.test.enterprise_project_id
  event_source              = huaweicloud_aomv4_alarm_rule.test.event_alarm_spec.0.event_source
  alarm_rule_type           = huaweicloud_aomv4_alarm_rule.test.type
  alarm_rule_status         = huaweicloud_aomv4_alarm_rule.test.status
  bind_notification_rule_id = huaweicloud_aomv4_alarm_rule.test.alarm_notifications.0.bind_notification_rule_id
}

locals {
  test_results = data.huaweicloud_aom_alarm_rules.test
}

output "name_validation" {
  value = alltrue([for v in local.test_results.alarm_rules[*].name : v == huaweicloud_aomv4_alarm_rule.test.name])
}

output "eps_id_validation" {
  value = alltrue([for v in local.test_results.alarm_rules[*].enterprise_project_id : 
	v == huaweicloud_aomv4_alarm_rule.test.enterprise_project_id])
}

output "event_source_validation" {
  value = alltrue([for v in local.test_results.alarm_rules[*].event_alarm_spec.0.event_source : 
	v == huaweicloud_aomv4_alarm_rule.test.event_alarm_spec.0.event_source])
}

output "type_validation" {
  value = alltrue([for v in local.test_results.alarm_rules[*].type : v == huaweicloud_aomv4_alarm_rule.test.type])
}

output "status_validation" {
  value = alltrue([for v in local.test_results.alarm_rules[*].status : v == huaweicloud_aomv4_alarm_rule.test.status])
}

output "bind_notification_rule_id_validation" {
  value = alltrue([for v in local.test_results.alarm_rules[*].alarm_notifications.0.bind_notification_rule_id : 
	v == huaweicloud_aomv4_alarm_rule.test.alarm_notifications.0.bind_notification_rule_id])
}
`, testAlarmRuleV4_basic(name))
}

func TestAccDataSourceAlarmRules_metric(t *testing.T) {
	dataSource := "data.huaweicloud_aom_alarm_rules.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAlarmRules_metric(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_rules.#"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_rules.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_rules.0.alarm_rule_id"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_rules.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_rules.0.enable"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_rules.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_rules.0.alarm_notifications.#"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_rules.0.metric_alarm_spec.#"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_rules.0.status"),
					resource.TestMatchResourceAttr(dataSource, "alarm_rules.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),

					resource.TestCheckOutput("name_validation", "true"),
					resource.TestCheckOutput("eps_id_validation", "true"),
					resource.TestCheckOutput("type_validation", "true"),
					resource.TestCheckOutput("status_validation", "true"),
				),
			},
		},
	})
}

func testDataSourceAlarmRules_metric(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_aom_alarm_rules" "test" {
  depends_on = [huaweicloud_aomv4_alarm_rule.test]

  alarm_rule_name       = huaweicloud_aomv4_alarm_rule.test.name
  enterprise_project_id = huaweicloud_aomv4_alarm_rule.test.enterprise_project_id
  alarm_rule_type       = huaweicloud_aomv4_alarm_rule.test.type
  alarm_rule_status     = huaweicloud_aomv4_alarm_rule.test.status
}

locals {
  test_results = data.huaweicloud_aom_alarm_rules.test
}

output "name_validation" {
  value = alltrue([for v in local.test_results.alarm_rules[*].name : v == huaweicloud_aomv4_alarm_rule.test.name])
}

output "eps_id_validation" {
  value = alltrue([for v in local.test_results.alarm_rules[*].enterprise_project_id : 
	v == huaweicloud_aomv4_alarm_rule.test.enterprise_project_id])
}

output "type_validation" {
  value = alltrue([for v in local.test_results.alarm_rules[*].type : v == huaweicloud_aomv4_alarm_rule.test.type])
}

output "status_validation" {
  value = alltrue([for v in local.test_results.alarm_rules[*].status : v == huaweicloud_aomv4_alarm_rule.test.status])
}
`, testAlarmRuleV4_metric(name))
}
