package lts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/lts"
)

func getSQLAlarmRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	getSQLAlarmRuleClient, err := cfg.NewServiceClient("lts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}

	return lts.GetSQLAlarmRuleById(getSQLAlarmRuleClient, state.Primary.ID)
}

func TestAccSQLAlarmRule_basic(t *testing.T) {
	var (
		name      = acceptance.RandomAccResourceName()
		aliasName = acceptance.RandomAccResourceName()
		password  = acceptance.RandomPassword()

		obj                interface{}
		rName              = "huaweicloud_lts_sql_alarm_rule.with_action_alarm_rule_name"
		withNotiSaveRule   = "huaweicloud_lts_sql_alarm_rule.with_notification_save_rule"
		rc                 = acceptance.InitResourceCheck(rName, &obj, getSQLAlarmRuleResourceFunc)
		rcWithNotiSaveRule = acceptance.InitResourceCheck(withNotiSaveRule, &obj, getSQLAlarmRuleResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDomainId(t)
			acceptance.TestAccPreCheckLtsAlarmActionRuleName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rc.CheckResourceDestroy(),
			rcWithNotiSaveRule.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccSQLAlarmRule_basic_step1(name, password, aliasName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "frequency.0.type", "HOURLY"),
					resource.TestCheckResourceAttr(rName, "condition_expression", "t>0"),
					resource.TestCheckResourceAttr(rName, "alarm_level", "CRITICAL"),
					resource.TestCheckResourceAttr(rName, "sql_requests.0.title", name),
					resource.TestCheckResourceAttr(rName, "sql_requests.0.sql", "select count(*) as t"),
					resource.TestCheckResourceAttrPair(rName, "sql_requests.0.log_group_id",
						"huaweicloud_lts_group.test.0", "id"),
					resource.TestCheckResourceAttrPair(rName, "sql_requests.0.log_stream_id",
						"huaweicloud_lts_stream.test.0", "id"),
					resource.TestCheckResourceAttr(rName, "sql_requests.0.search_time_range_unit", "minute"),
					resource.TestCheckResourceAttr(rName, "sql_requests.0.search_time_range", "5"),
					// Check optional parameters.
					resource.TestCheckResourceAttrPair(rName, "sql_requests.0.log_stream_name",
						"huaweicloud_lts_stream.test.0", "stream_name"),
					resource.TestCheckResourceAttrPair(rName, "sql_requests.0.log_group_name",
						"huaweicloud_lts_group.test.0", "group_name"),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(rName, "send_notifications", "true"),
					resource.TestCheckResourceAttr(rName, "trigger_condition_count", "1"),
					resource.TestCheckResourceAttr(rName, "trigger_condition_frequency", "1"),
					resource.TestCheckResourceAttr(rName, "send_recovery_notifications", "false"),
					resource.TestCheckResourceAttr(rName, "recovery_frequency", "3"),
					resource.TestCheckResourceAttr(rName, "notification_frequency", "0"),
					resource.TestCheckResourceAttr(rName, "alarm_rule_alias", name),
					resource.TestCheckResourceAttr(rName, "alarm_action_rule_name", acceptance.HW_LTS_ALARM_ACTION_RULE_NAME),
					resource.TestCheckResourceAttr(rName, "status", "STOPPING"),
					// Check Attributes.
					resource.TestCheckResourceAttrSet(rName, "domain_id"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),

					rcWithNotiSaveRule.CheckResourceExists(),
					resource.TestCheckResourceAttr(withNotiSaveRule, "send_notifications", "true"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "trigger_condition_count", "2"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "trigger_condition_frequency", "3"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "send_recovery_notifications", "true"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "recovery_frequency", "4"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "notification_frequency", "15"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "alarm_rule_alias", aliasName),
					resource.TestCheckResourceAttr(withNotiSaveRule, "notification_save_rule.0.template_name", "sql_template"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "notification_save_rule.0.topics.#", "1"),
					resource.TestCheckResourceAttrPair(withNotiSaveRule, "notification_save_rule.0.topics.0.name",
						"huaweicloud_smn_topic.test.0", "name"),
					resource.TestCheckResourceAttrPair(withNotiSaveRule, "notification_save_rule.0.topics.0.topic_urn",
						"huaweicloud_smn_topic.test.0", "topic_urn"),
					resource.TestCheckResourceAttrPair(withNotiSaveRule, "notification_save_rule.0.topics.0.display_name",
						"huaweicloud_smn_topic.test.0", "display_name"),
					resource.TestCheckResourceAttrPair(withNotiSaveRule, "notification_save_rule.0.topics.0.push_policy",
						"huaweicloud_smn_topic.test.0", "push_policy"),
				),
			},
			{
				Config: testAccSQLAlarmRule_basic_step2(name, password, aliasName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "condition_expression", "t>2"),
					resource.TestCheckResourceAttr(rName, "alarm_level", "INFO"),
					resource.TestCheckResourceAttr(rName, "sql_requests.0.title", name+"_sql"),
					resource.TestCheckResourceAttr(rName, "sql_requests.0.sql", "select *"),
					resource.TestCheckResourceAttrPair(rName, "sql_requests.0.log_group_id",
						"huaweicloud_lts_group.test.1", "id"),
					resource.TestCheckResourceAttrPair(rName, "sql_requests.0.log_stream_id",
						"huaweicloud_lts_stream.test.1", "id"),
					resource.TestCheckResourceAttr(rName, "sql_requests.0.search_time_range_unit", "hour"),
					resource.TestCheckResourceAttr(rName, "sql_requests.0.search_time_range", "10"),
					resource.TestCheckResourceAttr(rName, "sql_requests.0.is_time_range_relative", "true"),
					resource.TestCheckResourceAttr(rName, "frequency.0.type", "DAILY"),
					resource.TestCheckResourceAttr(rName, "frequency.0.hour_of_day", "6"),
					// Check Optional parameter.
					resource.TestCheckResourceAttr(rName, "trigger_condition_count", "6"),
					resource.TestCheckResourceAttr(rName, "trigger_condition_frequency", "6"),
					resource.TestCheckResourceAttr(rName, "send_recovery_notifications", "true"),
					resource.TestCheckResourceAttr(rName, "recovery_frequency", "5"),
					resource.TestCheckResourceAttr(rName, "notification_frequency", "30"),
					resource.TestCheckResourceAttrPair(rName, "sql_requests.0.log_stream_name",
						"huaweicloud_lts_stream.test.1", "stream_name"),
					resource.TestCheckResourceAttrPair(rName, "sql_requests.0.log_group_name",
						"huaweicloud_lts_group.test.1", "group_name"),
					resource.TestCheckResourceAttr(rName, "notification_save_rule.0.template_name", "sql_template"),
					resource.TestCheckResourceAttr(rName, "notification_save_rule.0.topics.#", "1"),
					resource.TestCheckResourceAttrPair(rName, "notification_save_rule.0.topics.0.name",
						"huaweicloud_smn_topic.test.0", "name"),
					resource.TestCheckResourceAttrPair(rName, "notification_save_rule.0.topics.0.topic_urn",
						"huaweicloud_smn_topic.test.0", "topic_urn"),
					resource.TestCheckResourceAttr(rName, "notification_save_rule.0.topics.0.display_name", ""),
					resource.TestCheckResourceAttr(rName, "notification_save_rule.0.topics.0.push_policy", "0"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "status", "RUNNING"),
					// Check attributes.
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					rcWithNotiSaveRule.CheckResourceExists(),
					resource.TestCheckResourceAttr(withNotiSaveRule, "frequency.0.type", "CRON"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "frequency.0.cron_expression", "0 18 * * *"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "send_recovery_notifications", "false"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "notification_save_rule.0.template_name", "sql_template"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "notification_save_rule.0.topics.#", "2"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "alarm_rule_alias", aliasName+"_update"),
				),
			},
			{
				Config: testAccSQLAlarmRule_basic_step3(name, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "frequency.0.type", "WEEKLY"),
					resource.TestCheckResourceAttr(rName, "frequency.0.hour_of_day", "6"),
					resource.TestCheckResourceAttr(rName, "frequency.0.day_of_week", "2"),
					resource.TestCheckResourceAttr(rName, "send_notifications", "true"),
					resource.TestCheckResourceAttr(rName, "alarm_action_rule_name", acceptance.HW_LTS_ALARM_ACTION_RULE_NAME),
					resource.TestCheckResourceAttr(rName, "notification_save_rule.#", "0"),
					rcWithNotiSaveRule.CheckResourceExists(),
					resource.TestCheckResourceAttr(withNotiSaveRule, "frequency.0.type", "FIXED_RATE"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "frequency.0.fixed_rate_unit", "minute"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "frequency.0.fixed_rate", "30"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "send_notifications", "false"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "notification_save_rule.#", "0"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:            withNotiSaveRule,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"notification_save_rule.0.user_name", "notification_save_rule.0.timezone"},
			},
		},
	})
}

func testAccSQLAlarmRule_base(name, password string) string {
	return fmt.Sprintf(`
%[1]s

# 'sql_template' is the default template provided by the system.
locals {
  sql_template_name = [for v in data.huaweicloud_lts_notification_templates.test.templates[*].name :v if v == "sql_template"][0]
}
`, testAlarmRule_base(name, password), name)
}

func testAccSQLAlarmRule_basic_step1(name, password, aliasName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_sql_alarm_rule" "with_action_alarm_rule_name" {
  name                   = "%[2]s"
  description            = "created by terraform"
  condition_expression   = "t>0"
  alarm_level            = "CRITICAL"
  send_notifications     = true
  alarm_action_rule_name = "%[4]s"
  status                 = "STOPPING"

  sql_requests {
    title                  = "%[2]s"
    sql                    = "select count(*) as t"
    log_group_id           = huaweicloud_lts_group.test[0].id
    log_stream_id          = huaweicloud_lts_stream.test[0].id
    search_time_range_unit = "minute"
    search_time_range      = 5
  }

  frequency {
    type = "HOURLY"
  }
}

resource "huaweicloud_lts_sql_alarm_rule" "with_notification_save_rule" {
  name                        = "%[2]s_sql"
  condition_expression        = "t>0"
  alarm_level                 = "MINOR"
  send_notifications          = true
  trigger_condition_count     = 2
  trigger_condition_frequency = 3
  send_recovery_notifications = true
  recovery_frequency          = 4
  notification_frequency      = 15
  alarm_rule_alias            = "%[3]s"

  sql_requests {
    title                  = "%[2]s"
    sql                    = "select count(*) as t"
    log_group_id           = huaweicloud_lts_group.test[0].id
    log_stream_id          = huaweicloud_lts_stream.test[0].id
    search_time_range_unit = "minute"
    search_time_range      = 5
    log_group_name         = huaweicloud_lts_group.test[0].group_name
    log_stream_name        = huaweicloud_lts_stream.test[0].stream_name
  }

  frequency {
    type = "HOURLY"
  }

  notification_save_rule {
    template_name = local.sql_template_name
    user_name     = huaweicloud_identity_user.test.name
    language      = "en-us"

    topics {
      name         = huaweicloud_smn_topic.test[0].name
      topic_urn    = huaweicloud_smn_topic.test[0].topic_urn
      display_name = huaweicloud_smn_topic.test[0].display_name
      push_policy  = huaweicloud_smn_topic.test[0].push_policy
    }
  }
}
`, testAccSQLAlarmRule_base(name, password), name, aliasName, acceptance.HW_LTS_ALARM_ACTION_RULE_NAME)
}

func testAccSQLAlarmRule_basic_step2(name, password, aliasName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_sql_alarm_rule" "with_action_alarm_rule_name" {
  name                        = "%[2]s"
  condition_expression        = "t>2"
  alarm_level                 = "INFO"
  trigger_condition_count     = 6
  trigger_condition_frequency = 6
  send_recovery_notifications = true
  recovery_frequency          = 5
  send_notifications          = true
  notification_frequency      = 30
  status                      = "RUNNING"

  sql_requests {
    title                  = "%[2]s_sql"
    sql                    = "select *"
    log_group_id           = huaweicloud_lts_group.test[1].id
    log_stream_id          = huaweicloud_lts_stream.test[1].id
    search_time_range_unit = "hour"
    search_time_range      = 10
    is_time_range_relative = true
    log_group_name         = huaweicloud_lts_group.test[1].group_name
    log_stream_name        = huaweicloud_lts_stream.test[1].stream_name
  }

  frequency {
    type        = "DAILY"
    hour_of_day = 6
  }

  # Verify the 'sql_alarm_send_code' is '1' in the modification logic.
  notification_save_rule {
    template_name = local.sql_template_name
    user_name     = huaweicloud_identity_user.test.name
    language      = "en-us"

    topics {
      name      = huaweicloud_smn_topic.test[0].name
      topic_urn = huaweicloud_smn_topic.test[0].topic_urn
    }
  }
}

resource "huaweicloud_lts_sql_alarm_rule" "with_notification_save_rule" {
  name                 = "%[2]s_notification"
  condition_expression = "t>0"
  alarm_level          = "MINOR"
  send_notifications   = true
  alarm_rule_alias     = "%[3]s_update"

  sql_requests {
    title                  = "%[2]s_sql"
    sql                    = "select count(*) as t"
    log_group_id           = huaweicloud_lts_group.test[0].id
    log_stream_id          = huaweicloud_lts_stream.test[0].id
    search_time_range_unit = "minute"
    search_time_range      = 5
    log_group_name         = huaweicloud_lts_group.test[0].group_name
    log_stream_name        = huaweicloud_lts_stream.test[0].stream_name
  }

  frequency {
    type            = "CRON"
    cron_expression = "0 18 * * *"
  }

  notification_save_rule {
    template_name = local.sql_template_name
    user_name     = huaweicloud_identity_user.test.name
    language      = "en-us"

    # Verify the 'sql_alarm_send_code' is '2' in the modification logic.
    topics {
      name         = huaweicloud_smn_topic.test[0].name
      topic_urn    = huaweicloud_smn_topic.test[0].topic_urn
      display_name = huaweicloud_smn_topic.test[0].display_name
      push_policy  = huaweicloud_smn_topic.test[0].push_policy
    }
    topics {
      name      = huaweicloud_smn_topic.test[1].name
      topic_urn = huaweicloud_smn_topic.test[1].topic_urn
    }
  }
}
`, testAccSQLAlarmRule_base(name, password), name, aliasName)
}

func testAccSQLAlarmRule_basic_step3(name, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_sql_alarm_rule" "with_action_alarm_rule_name" {
  name                        = "%[2]s"
  condition_expression        = "t>2"
  alarm_level                 = "INFO"
  trigger_condition_count     = 6
  trigger_condition_frequency = 6
  recovery_frequency          = 5
  send_notifications          = true
  notification_frequency      = 30
  alarm_action_rule_name      = "%[3]s"

  sql_requests {
    title                  = "%[2]s_sql"
    sql                    = "select *"
    log_group_id           = huaweicloud_lts_group.test[1].id
    log_stream_id          = huaweicloud_lts_stream.test[1].id
    search_time_range_unit = "hour"
    search_time_range      = 10
    is_time_range_relative = true
    log_group_name         = huaweicloud_lts_group.test[1].group_name
    log_stream_name        = huaweicloud_lts_stream.test[1].stream_name
  }

  frequency {
    type        = "WEEKLY"
    hour_of_day = 6
    day_of_week = 2
  }
}

resource "huaweicloud_lts_sql_alarm_rule" "with_notification_save_rule" {
  name                 = "%[2]s_notification"
  condition_expression = "t>0"
  alarm_level          = "MINOR"

  sql_requests {
    title                  = "%[2]s_sql"
    sql                    = "select count(*) as t"
    log_group_id           = huaweicloud_lts_group.test[0].id
    log_stream_id          = huaweicloud_lts_stream.test[0].id
    search_time_range_unit = "minute"
    search_time_range      = 5
  }

  frequency {
    type            = "FIXED_RATE"
    fixed_rate_unit = "minute"
    fixed_rate      = "30"
  }

  # Remove 'notification_save_rule', verify the 'sql_alarm_send_code' is '3' in the modification logic.
}
`, testAccSQLAlarmRule_base(name, password), name, acceptance.HW_LTS_ALARM_ACTION_RULE_NAME)
}
