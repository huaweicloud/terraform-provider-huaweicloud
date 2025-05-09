package lts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceKeywordAlarmRules_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	dataSource := "data.huaweicloud_lts_keyword_alarm_rules.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLtsAlarmActionRuleName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKeywordAlarmRules_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "keyword_alarm_rules.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckOutput("is_name_set", "true"),
					resource.TestCheckOutput("is_keywords_requests_set", "true"),
					resource.TestCheckOutput("is_frequency_set", "true"),
					resource.TestCheckOutput("is_alarm_level_set", "true"),
					resource.TestCheckOutput("is_send_notifications_set", "true"),
					resource.TestCheckOutput("is_alarm_action_rule_name_set", "true"),
					resource.TestCheckOutput("is_trigger_condition_count_set", "true"),
					resource.TestCheckOutput("is_trigger_condition_frequency_set", "true"),
					resource.TestCheckOutput("is_send_recovery_notifications_set", "true"),
					resource.TestCheckOutput("is_recovery_frequency_set", "true"),
					resource.TestCheckOutput("is_notification_frequency_set", "true"),
					resource.TestCheckOutput("is_status_set", "true"),
					resource.TestCheckOutput("is_topics_set", "true"),
					resource.TestCheckOutput("is_template_name_set", "true"),
					resource.TestMatchResourceAttr(dataSource, "keyword_alarm_rules.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource, "keyword_alarm_rules.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(dataSource, "keyword_alarm_rules.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "keyword_alarm_rules.0.condition_expression"),
				),
			},
		},
	})
}

func testAccDataSourceKeywordAlarmRules_base(name string) string {
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
  display_name = "%[1]s"
}

resource "huaweicloud_lts_keywords_alarm_rule" "test" {
  count                       = 2
  name                        = "%[1]s${count.index}"
  alarm_level                 = "CRITICAL"
  description                 = "created by terraform"
  send_notifications          = true
  alarm_action_rule_name      = count.index == 0 ? "%[2]s" : null
  trigger_condition_count     = 2
  trigger_condition_frequency = 4
  send_recovery_notifications = true
  recovery_frequency          = 5
  notification_frequency      = 10

  keywords_requests {
    keywords               = "tf_test1_key_words"
    condition              = ">"
    number                 = 100
    log_group_id           = huaweicloud_lts_group.test.id
    log_stream_id          = huaweicloud_lts_stream.test.id
    search_time_range_unit = "minute"
    search_time_range      = 5
  }

  frequency {
    type        = "WEEKLY"
    day_of_week = 6
    hour_of_day = 10
  }

  dynamic "notification_save_rule" {
    for_each = count.index == 0 ? [] : [1]

    content {
      template_name = "keywords_template"
      user_name     = "tf-user13"

      topics {
        name         = huaweicloud_smn_topic.test.name
        topic_urn    = huaweicloud_smn_topic.test.topic_urn
        display_name = huaweicloud_smn_topic.test.display_name
        push_policy  = huaweicloud_smn_topic.test.push_policy
      }
    }
  }
}
`, name, acceptance.HW_LTS_ALARM_ACTION_RULE_NAME)
}

func testAccDataSourceKeywordAlarmRules_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_lts_keyword_alarm_rules" "test" {
  depends_on = [huaweicloud_lts_keywords_alarm_rule.test]
}

locals {
  keyword_alarm_rule = try([for v in data.huaweicloud_lts_keyword_alarm_rules.test.keyword_alarm_rules :
  v if v.id == huaweicloud_lts_keywords_alarm_rule.test[0].id][0], {})
  keywords_request          = try(local.keyword_alarm_rule.keywords_requests[0], {})
  resource_keywords_request = try(huaweicloud_lts_keywords_alarm_rule.test[0].keywords_requests[0], {})
  frequency                 = try(local.keyword_alarm_rule.frequency[0], {})
  resource_frequency        = try(huaweicloud_lts_keywords_alarm_rule.test[0].frequency[0], {})

  alarm_rule_with_notification_save_rule = try(
    [for v in data.huaweicloud_lts_keyword_alarm_rules.test.keyword_alarm_rules :
    v if v.id == huaweicloud_lts_keywords_alarm_rule.test[1].id][0], {}
  )
}

output "is_name_set" {
  value = try(local.keyword_alarm_rule.name == huaweicloud_lts_keywords_alarm_rule.test[0].name, false)
}

output "is_keywords_requests_set" {
  value = try(
    length(local.keyword_alarm_rule.keywords_requests) == length(huaweicloud_lts_keywords_alarm_rule.test[0].keywords_requests) &&
    local.keywords_request.keywords == local.resource_keywords_request.keywords &&
    local.keywords_request.condition == local.resource_keywords_request.condition &&
    local.keywords_request.number == local.resource_keywords_request.number &&
    local.keywords_request.log_stream_id == local.resource_keywords_request.log_stream_id &&
    local.keywords_request.search_time_range_unit == local.resource_keywords_request.search_time_range_unit &&
    local.keywords_request.search_time_range == local.resource_keywords_request.search_time_range &&
    local.keywords_request.log_stream_name == local.resource_keywords_request.log_stream_name &&
    local.keywords_request.log_group_name == local.resource_keywords_request.log_group_name
  , false)
}

output "is_frequency_set" {
  value = try(
    local.frequency.type == local.resource_frequency.type &&
    local.frequency.cron_expression == local.resource_frequency.cron_expression &&
    local.frequency.day_of_week == local.resource_frequency.day_of_week &&
    local.frequency.fixed_rate == local.resource_frequency.fixed_rate &&
    local.frequency.fixed_rate_unit == local.resource_frequency.fixed_rate_unit &&
    local.frequency.hour_of_day == local.resource_frequency.hour_of_day
  , false)
}

output "is_alarm_level_set" {
  value = try(local.keyword_alarm_rule.alarm_level == huaweicloud_lts_keywords_alarm_rule.test[0].alarm_level, false)
}

output "is_send_notifications_set" {
  value = try(local.keyword_alarm_rule.send_notifications == huaweicloud_lts_keywords_alarm_rule.test[0].send_notifications, false)
}

output "is_alarm_action_rule_name_set" {
  value = try(local.keyword_alarm_rule.alarm_action_rule_name == huaweicloud_lts_keywords_alarm_rule.test[0].alarm_action_rule_name, false)
}

output "is_trigger_condition_count_set" {
  value = try(local.keyword_alarm_rule.trigger_condition_count == huaweicloud_lts_keywords_alarm_rule.test[0].trigger_condition_count, false)
}

output "is_trigger_condition_frequency_set" {
  value = try(local.keyword_alarm_rule.trigger_condition_frequency == huaweicloud_lts_keywords_alarm_rule.test[0].trigger_condition_frequency, false)
}

output "is_send_recovery_notifications_set" {
  value = try(local.keyword_alarm_rule.send_recovery_notifications == huaweicloud_lts_keywords_alarm_rule.test[0].send_recovery_notifications, false)
}

output "is_recovery_frequency_set" {
  value = try(local.keyword_alarm_rule.recovery_frequency == huaweicloud_lts_keywords_alarm_rule.test[0].recovery_frequency, false)
}

output "is_notification_frequency_set" {
  value = try(local.keyword_alarm_rule.notification_frequency == huaweicloud_lts_keywords_alarm_rule.test[0].notification_frequency, false)
}

output "is_status_set" {
  value = try(local.keyword_alarm_rule.status == huaweicloud_lts_keywords_alarm_rule.test[0].status, false)
}

output "is_topics_set" {
  value = try(local.alarm_rule_with_notification_save_rule.topics ==
  huaweicloud_lts_keywords_alarm_rule.test[1].notification_save_rule[0].topics, false)
}

output "is_template_name_set" {
  value = try(local.alarm_rule_with_notification_save_rule.template_name ==
  huaweicloud_lts_keywords_alarm_rule.test[1].notification_save_rule[0].template_name, false)
}
`, testAccDataSourceKeywordAlarmRules_base(name))
}
