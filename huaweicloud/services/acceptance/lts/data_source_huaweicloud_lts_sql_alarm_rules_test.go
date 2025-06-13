package lts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSqlAlarmRules_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSource := "data.huaweicloud_lts_sql_alarm_rules.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSqlAlarmRules_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "sql_alarm_rules.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckResourceAttrSet(dataSource, "sql_alarm_rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_alarm_rules.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_alarm_rules.0.condition_expression"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_alarm_rules.0.alarm_level"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_alarm_rules.0.status"),
					resource.TestCheckOutput("is_sql_requests_set", "true"),
					resource.TestCheckOutput("is_frequency_set", "true"),
					resource.TestCheckOutput("is_send_notifications_set", "true"),
					resource.TestCheckOutput("is_alarm_action_rule_name_set", "true"),
					resource.TestCheckOutput("is_trigger_condition_count_set", "true"),
					resource.TestCheckOutput("is_trigger_condition_frequency_set", "true"),
					resource.TestCheckOutput("is_send_recovery_notifications_set", "true"),
					resource.TestCheckOutput("is_recovery_frequency_set", "true"),
					resource.TestCheckOutput("is_notification_frequency_set", "true"),
					resource.TestCheckOutput("is_description_set", "true"),
					resource.TestCheckOutput("is_topics_set", "true"),
					resource.TestCheckOutput("is_template_name_set", "true"),
					resource.TestMatchResourceAttr(dataSource, "sql_alarm_rules.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource, "sql_alarm_rules.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataSourceSqlAlarmRules_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[1]s"
}

resource "huaweicloud_smn_topic" "test" {
  name         = "%[1]s"
  display_name = "The display name of topic"
}

resource "huaweicloud_identity_user" "test" {
  name     = "%[1]s"
  enabled  = true
  password = "%[2]s"
}

resource "huaweicloud_lts_sql_alarm_rule" "test" {
  name                        = "%[1]s"
  description                 = "created by terraform"
  condition_expression        = "t>0"
  alarm_level                 = "CRITICAL"
  send_notifications          = true
  alarm_action_rule_name      = "rule"
  trigger_condition_count     = 2
  trigger_condition_frequency = 3
  send_recovery_notifications = true
  recovery_frequency          = 4
  notification_frequency      = 15
  alarm_rule_alias            = "%[1]s-ailas"

  sql_requests {
    title                  = "%[1]s"
    sql                    = "select count(*) as t"
    log_group_id           = huaweicloud_lts_group.test.id
    log_stream_id          = huaweicloud_lts_stream.test.id
    search_time_range_unit = "minute"
    search_time_range      = 5
    log_group_name         = huaweicloud_lts_group.test.group_name
    log_stream_name        = huaweicloud_lts_stream.test.stream_name
  }

  frequency {
    type        = "WEEKLY"
    hour_of_day = 6
    day_of_week = 2
  }
}

resource "huaweicloud_lts_sql_alarm_rule" "with_notification_save_rule" {
  name                 = "%[1]s_notification"
  condition_expression = "t>0"
  alarm_level          = "MINOR"
  send_notifications   = true

  sql_requests {
    title                  = "%[1]s"
    sql                    = "select count(*) as t"
    log_group_id           = huaweicloud_lts_group.test.id
    log_stream_id          = huaweicloud_lts_stream.test.id
    search_time_range_unit = "minute"
    search_time_range      = 5
  }

  frequency {
    type        = "DAILY"
    hour_of_day = 6
  }

  notification_save_rule {
    template_name = "sql_template"
    user_name     = huaweicloud_identity_user.test.name
    language      = "en-us"

    topics {
      name         = huaweicloud_smn_topic.test.name
      topic_urn    = huaweicloud_smn_topic.test.topic_urn
      display_name = huaweicloud_smn_topic.test.display_name
      push_policy  = huaweicloud_smn_topic.test.push_policy
    }
  }
}
`, name, acceptance.RandomPassword())
}

func testAccDataSourceSqlAlarmRules_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_lts_sql_alarm_rules" "test" {
  depends_on = [
    huaweicloud_lts_sql_alarm_rule.test,
    huaweicloud_lts_sql_alarm_rule.with_notification_save_rule
  ]
}

locals {
  sql_alarm_rules = data.huaweicloud_lts_sql_alarm_rules.test.sql_alarm_rules
  sql_alarm_rule  = try([for v in local.sql_alarm_rules : v if v.id == huaweicloud_lts_sql_alarm_rule.test.id][0], {})
  sql_request     = try(local.sql_alarm_rule.sql_requests[0], {})
  frequency       = try(local.sql_alarm_rule.frequency[0], {})
  alarm_rule_with_notification_save_rule = try([for v in local.sql_alarm_rules :
  v if v.id == huaweicloud_lts_sql_alarm_rule.with_notification_save_rule.id][0], {})
}

output "is_sql_requests_set" {
  value = try(
    length(local.sql_alarm_rule.sql_requests) == length(huaweicloud_lts_sql_alarm_rule.test.sql_requests) &&
    local.sql_request.title != "" &&
    local.sql_request.sql != "" &&
    local.sql_request.log_stream_id != "" &&
    local.sql_request.log_stream_name != "" &&
    local.sql_request.log_group_id != "" &&
    local.sql_request.log_group_name != "" &&
    local.sql_request.search_time_range > 0 &&
    local.sql_request.search_time_range_unit != ""
  , false)
}

output "is_frequency_set" {
  value = try(local.frequency.type == "WEEKLY" && local.frequency.hour_of_day == 6 && local.frequency.day_of_week == 2, false)
}

output "is_send_notifications_set" {
  value = try(local.sql_alarm_rule.send_notifications == huaweicloud_lts_sql_alarm_rule.test.send_notifications, false)
}

output "is_alarm_action_rule_name_set" {
  value = try(local.sql_alarm_rule.alarm_action_rule_name == huaweicloud_lts_sql_alarm_rule.test.alarm_action_rule_name, false)
}

output "is_trigger_condition_count_set" {
  value = try(local.sql_alarm_rule.trigger_condition_count == huaweicloud_lts_sql_alarm_rule.test.trigger_condition_count, false)
}

output "is_trigger_condition_frequency_set" {
  value = try(local.sql_alarm_rule.trigger_condition_frequency == huaweicloud_lts_sql_alarm_rule.test.trigger_condition_frequency, false)
}

output "is_send_recovery_notifications_set" {
  value = try(local.sql_alarm_rule.send_recovery_notifications == huaweicloud_lts_sql_alarm_rule.test.send_recovery_notifications, false)
}

output "is_recovery_frequency_set" {
  value = try(local.sql_alarm_rule.recovery_frequency == huaweicloud_lts_sql_alarm_rule.test.recovery_frequency, false)
}

output "is_notification_frequency_set" {
  value = try(local.sql_alarm_rule.notification_frequency == huaweicloud_lts_sql_alarm_rule.test.notification_frequency, false)
}

output "is_description_set" {
  value = try(local.sql_alarm_rule.description == huaweicloud_lts_sql_alarm_rule.test.description, false)
}

output "is_topics_set" {
  value = try(local.alarm_rule_with_notification_save_rule.topics ==
  huaweicloud_lts_sql_alarm_rule.with_notification_save_rule.notification_save_rule[0].topics, false)
}

output "is_template_name_set" {
  value = try(local.alarm_rule_with_notification_save_rule.template_name ==
  huaweicloud_lts_sql_alarm_rule.with_notification_save_rule.notification_save_rule[0].template_name, false)
}
`, testAccDataSourceSqlAlarmRules_base(name))
}
