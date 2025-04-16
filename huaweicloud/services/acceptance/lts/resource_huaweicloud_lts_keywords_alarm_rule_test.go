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

func getKeywordsAlarmRule(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("lts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}

	return lts.GetKeywordsAlarmRuleById(client, state.Primary.ID)
}

func TestAccKeywordsAlarmRule_basic(t *testing.T) {
	var (
		name      = acceptance.RandomAccResourceName()
		aliasName = acceptance.RandomAccResourceName()
		password  = acceptance.RandomPassword()

		keywordAlarmRule interface{}
		rName            = "huaweicloud_lts_keywords_alarm_rule.test"
		rc               = acceptance.InitResourceCheck(rName, &keywordAlarmRule, getKeywordsAlarmRule)

		withNotiSaveRule   = "huaweicloud_lts_keywords_alarm_rule.with_notification_save_rule"
		rcWithNotiSaveRule = acceptance.InitResourceCheck(withNotiSaveRule, &keywordAlarmRule, getKeywordsAlarmRule)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDomainId(t)
			acceptance.TestAccPreCheckLtsAlarmActionRuleName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeywordsAlarmRule_step1(name, password, aliasName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "keywords_requests.0.keywords", name+"_key_words"),
					resource.TestCheckResourceAttr(rName, "keywords_requests.0.condition", ">"),
					resource.TestCheckResourceAttr(rName, "keywords_requests.0.number", "100"),
					resource.TestCheckResourceAttr(rName, "keywords_requests.0.search_time_range_unit", "minute"),
					resource.TestCheckResourceAttr(rName, "keywords_requests.0.search_time_range", "5"),
					resource.TestCheckResourceAttrPair(rName, "keywords_requests.0.log_group_id",
						"huaweicloud_lts_group.test.0", "id"),
					resource.TestCheckResourceAttrPair(rName, "keywords_requests.0.log_stream_id",
						"huaweicloud_lts_stream.test.0", "id"),
					resource.TestCheckResourceAttrPair(rName, "keywords_requests.0.log_group_name",
						"huaweicloud_lts_group.test.0", "group_name"),
					resource.TestCheckResourceAttrPair(rName, "keywords_requests.0.log_stream_name",
						"huaweicloud_lts_stream.test.0", "stream_name"),
					resource.TestCheckResourceAttr(rName, "frequency.0.type", "HOURLY"),
					resource.TestCheckResourceAttr(rName, "alarm_level", "CRITICAL"),
					// check optional parameters.
					resource.TestCheckResourceAttr(rName, "send_notifications", "false"),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(rName, "alarm_action_rule_name", ""),
					resource.TestCheckResourceAttr(rName, "notification_save_rule.#", "0"),
					resource.TestCheckResourceAttr(rName, "trigger_condition_count", "1"),
					resource.TestCheckResourceAttr(rName, "trigger_condition_frequency", "1"),
					resource.TestCheckResourceAttr(rName, "send_recovery_notifications", "false"),
					resource.TestCheckResourceAttr(rName, "recovery_frequency", "3"),
					resource.TestCheckResourceAttr(rName, "alarm_rule_alias", name),
					resource.TestCheckResourceAttr(rName, "notification_frequency", "0"),
					resource.TestCheckResourceAttr(rName, "status", "RUNNING"),
					// Check attributes.
					resource.TestCheckResourceAttrSet(rName, "domain_id"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "condition_expression"),
					// Mainly verify the modification logic of the `notification_save_rule` parameter.
					rcWithNotiSaveRule.CheckResourceExists(),
					resource.TestCheckResourceAttr(withNotiSaveRule, "status", "STOPPING"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "send_notifications", "true"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "trigger_condition_count", "2"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "trigger_condition_frequency", "3"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "send_recovery_notifications", "true"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "recovery_frequency", "4"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "alarm_rule_alias", aliasName),
					resource.TestCheckResourceAttr(withNotiSaveRule, "notification_frequency", "15"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "frequency.0.type", "FIXED_RATE"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "frequency.0.fixed_rate_unit", "minute"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "frequency.0.fixed_rate", "30"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "notification_save_rule.0.template_name", "keywords_template"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "notification_save_rule.0.topics.#", "1"),
					resource.TestCheckResourceAttrPair(withNotiSaveRule, "notification_save_rule.0.topics.0.name",
						"huaweicloud_smn_topic.test.0", "name"),
					resource.TestCheckResourceAttrPair(withNotiSaveRule, "notification_save_rule.0.topics.0.topic_urn",
						"huaweicloud_smn_topic.test.0", "topic_urn"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "notification_save_rule.0.topics.0.display_name", ""),
				),
			},
			{
				Config: testKeywordsAlarmRule_step2(name, password, aliasName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "keywords_requests.0.keywords", name+"_key_words_update"),
					resource.TestCheckResourceAttr(rName, "keywords_requests.0.condition", ">="),
					resource.TestCheckResourceAttr(rName, "keywords_requests.0.number", "50"),
					resource.TestCheckResourceAttr(rName, "keywords_requests.0.search_time_range_unit", "hour"),
					resource.TestCheckResourceAttr(rName, "keywords_requests.0.search_time_range", "1"),
					resource.TestCheckResourceAttrPair(rName, "keywords_requests.0.log_group_id",
						"huaweicloud_lts_group.test.1", "id"),
					resource.TestCheckResourceAttrPair(rName, "keywords_requests.0.log_stream_id",
						"huaweicloud_lts_stream.test.1", "id"),
					resource.TestCheckResourceAttrPair(rName, "keywords_requests.0.log_group_name",
						"huaweicloud_lts_group.test.1", "group_name"),
					resource.TestCheckResourceAttrPair(rName, "keywords_requests.0.log_stream_name",
						"huaweicloud_lts_stream.test.1", "stream_name"),
					resource.TestCheckResourceAttr(rName, "frequency.0.type", "DAILY"),
					resource.TestCheckResourceAttr(rName, "frequency.0.hour_of_day", "6"),
					resource.TestCheckResourceAttr(rName, "alarm_level", "INFO"),
					// check optional parameters.
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "send_notifications", "true"),
					resource.TestCheckResourceAttr(rName, "trigger_condition_count", "2"),
					resource.TestCheckResourceAttr(rName, "trigger_condition_frequency", "3"),
					resource.TestCheckResourceAttr(rName, "send_recovery_notifications", "true"),
					resource.TestCheckResourceAttr(rName, "recovery_frequency", "4"),
					resource.TestCheckResourceAttr(rName, "notification_frequency", "30"),
					resource.TestCheckResourceAttr(rName, "notification_save_rule.0.template_name", "keywords_template"),
					resource.TestCheckResourceAttrPair(rName, "notification_save_rule.0.user_name",
						"huaweicloud_identity_user.test", "name"),
					resource.TestCheckResourceAttr(rName, "notification_save_rule.0.topics.#", "1"),
					resource.TestCheckResourceAttrPair(rName, "notification_save_rule.0.topics.0.name",
						"huaweicloud_smn_topic.test.0", "name"),
					resource.TestCheckResourceAttrPair(rName, "notification_save_rule.0.topics.0.topic_urn",
						"huaweicloud_smn_topic.test.0", "topic_urn"),
					resource.TestCheckResourceAttrPair(rName, "notification_save_rule.0.topics.0.display_name",
						"huaweicloud_smn_topic.test.0", "display_name"),
					resource.TestCheckResourceAttrPair(rName, "notification_save_rule.0.topics.0.push_policy",
						"huaweicloud_smn_topic.test.0", "push_policy"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					rcWithNotiSaveRule.CheckResourceExists(),
					resource.TestCheckResourceAttr(withNotiSaveRule, "alarm_rule_alias", aliasName+"_update"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "status", "RUNNING"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "notification_save_rule.0.topics.#", "2"),
				),
			},
			{
				Config: testKeywordsAlarmRule_step3(name, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "frequency.0.type", "WEEKLY"),
					resource.TestCheckResourceAttr(rName, "frequency.0.day_of_week", "6"),
					resource.TestCheckResourceAttr(rName, "frequency.0.hour_of_day", "10"),
					resource.TestCheckResourceAttr(rName, "send_notifications", "true"),
					resource.TestCheckResourceAttr(rName, "alarm_action_rule_name", acceptance.HW_LTS_ALARM_ACTION_RULE_NAME),
					resource.TestCheckResourceAttr(rName, "notification_save_rule.#", "0"),
					rcWithNotiSaveRule.CheckResourceExists(),
					resource.TestCheckResourceAttr(withNotiSaveRule, "frequency.0.type", "CRON"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "frequency.0.cron_expression", "0 18 * * *"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "notification_save_rule.0.language", "en-us"),
				),
			},
			{
				Config: testKeywordsAlarmRule_step4(name, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "alarm_action_rule_name", acceptance.HW_LTS_ALARM_ACTION_RULE_NAME),
					rcWithNotiSaveRule.CheckResourceExists(),
					resource.TestCheckResourceAttr(withNotiSaveRule, "send_notifications", "false"),
					resource.TestCheckResourceAttr(rName, "notification_save_rule.#", "0"),
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

func testAlarmRule_base(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  count = 2

  group_name  = "%[1]s_${count.index}"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test" {
  count = 2

  group_id    = huaweicloud_lts_group.test[count.index].id
  stream_name = "%[1]s_${count.index}"
}

resource "huaweicloud_smn_topic" "test" {
  count = 2

  name         = "%[1]s_${count.index}"
  display_name = "The display name of topic"
}

data "huaweicloud_lts_notification_templates" "test" {
  domain_id = "%[3]s"
}

# 'keywords_template' is the default template provided by the system.
locals {
  template_name = [for v in data.huaweicloud_lts_notification_templates.test.templates[*].name : v if v == "keywords_template"][0]
}

resource "huaweicloud_identity_user" "test" {
  name     = "%[1]s"
  enabled  = true
  password = "%[2]s"
}
`, name, password, acceptance.HW_DOMAIN_ID)
}

func testKeywordsAlarmRule_step1(name, password, aliasName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_keywords_alarm_rule" "test" {
  name        = "%[2]s"
  alarm_level = "CRITICAL"
  description = "created by terraform"

  keywords_requests {
    keywords               = "%[2]s_key_words"
    condition              = ">"
    number                 = 100
    log_group_id           = huaweicloud_lts_group.test[0].id
    log_stream_id          = huaweicloud_lts_stream.test[0].id
    search_time_range_unit = "minute"
    search_time_range      = 5
  }

  frequency {
    type = "HOURLY"
  }
}

resource "huaweicloud_lts_keywords_alarm_rule" "with_notification_save_rule" {
  name                        = "%[2]s_notification"
  alarm_level                 = "MINOR"
  send_notifications          = true
  trigger_condition_count     = 2
  trigger_condition_frequency = 3
  send_recovery_notifications = true
  recovery_frequency          = 4
  alarm_rule_alias            = "%[3]s"
  notification_frequency      = 15
  status                      = "STOPPING"

  keywords_requests {
    keywords               = "%[2]s_key_words"
    condition              = "<"
    number                 = 50
    log_group_id           = huaweicloud_lts_group.test[0].id
    log_stream_id          = huaweicloud_lts_stream.test[0].id
    search_time_range_unit = "hour"
    search_time_range      = 1
    log_group_name         = huaweicloud_lts_group.test[0].group_name
    log_stream_name        = huaweicloud_lts_stream.test[0].stream_name
  }

  frequency {
    type            = "FIXED_RATE"
    fixed_rate_unit = "minute"
    fixed_rate      = 30
  }

  notification_save_rule {
    template_name = local.template_name
    user_name     = huaweicloud_identity_user.test.name

    topics {
      name      = huaweicloud_smn_topic.test[0].name
      topic_urn = huaweicloud_smn_topic.test[0].topic_urn
    }
  }
}
`, testAlarmRule_base(name, password), name, aliasName)
}

func testKeywordsAlarmRule_step2(name, password, aliasName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_keywords_alarm_rule" "test" {
  name                        = "%[2]s"
  alarm_level                 = "INFO"
  send_notifications          = true
  trigger_condition_count     = 2
  trigger_condition_frequency = 3
  send_recovery_notifications = true
  recovery_frequency          = 4
  notification_frequency      = 30

  keywords_requests {
    keywords               = "%[2]s_key_words_update"
    condition              = ">="
    number                 = 50
    log_group_id           = huaweicloud_lts_group.test[1].id
    log_stream_id          = huaweicloud_lts_stream.test[1].id
    search_time_range_unit = "hour"
    search_time_range      = 1
  }

  frequency {
    type        = "DAILY"
    hour_of_day = 6
  }

  notification_save_rule {
    template_name = local.template_name
    user_name     = huaweicloud_identity_user.test.name

    # Verify the 'keywords_alarm_send_code' is '1' in the modification logic.
    topics {
      name         = huaweicloud_smn_topic.test[0].name
      topic_urn    = huaweicloud_smn_topic.test[0].topic_urn
      display_name = huaweicloud_smn_topic.test[0].display_name
      push_policy  = huaweicloud_smn_topic.test[0].push_policy
    }
  }
}

resource "huaweicloud_lts_keywords_alarm_rule" "with_notification_save_rule" {
  name                        = "%[2]s_notification"
  alarm_level                 = "MINOR"
  send_notifications          = true
  trigger_condition_count     = 2
  trigger_condition_frequency = 3
  recovery_frequency          = 4
  alarm_rule_alias            = "%[3]s_update"
  notification_frequency      = 15
  status                      = "RUNNING"
  
  keywords_requests {
    keywords               = "%[2]s_key_words"
    condition              = "<"
    number                 = 50
    log_group_id           = huaweicloud_lts_group.test[0].id
    log_stream_id          = huaweicloud_lts_stream.test[0].id
    search_time_range_unit = "hour"
    search_time_range      = 1
  }

  frequency {
    type            = "FIXED_RATE"
    fixed_rate_unit = "minute"
    fixed_rate      = 30
  }

  notification_save_rule {
    template_name = local.template_name
    user_name     = huaweicloud_identity_user.test.name

    # Verify the 'keywords_alarm_send_code' is '2' in the modification logic.
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
`, testAlarmRule_base(name, password), name, aliasName)
}

func testKeywordsAlarmRule_step3(name, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_keywords_alarm_rule" "test" {
  name                        = "%[2]s"
  alarm_level                 = "INFO"
  send_notifications          = true
  trigger_condition_count     = 2
  trigger_condition_frequency = 3
  send_recovery_notifications = true
  recovery_frequency          = 4
  notification_frequency      = 15
  alarm_action_rule_name      = "%[3]s"

  keywords_requests {
    keywords               = "%[2]s_key_words_update"
    condition              = ">="
    number                 = 50
    log_group_id           = huaweicloud_lts_group.test[1].id
    log_stream_id          = huaweicloud_lts_stream.test[1].id
    search_time_range_unit = "hour"
    search_time_range      = 1
  }

  frequency {
    type        = "WEEKLY"
    day_of_week = 6
    hour_of_day = 10
  }
}

resource "huaweicloud_lts_keywords_alarm_rule" "with_notification_save_rule" {
  name                        = "%[2]s_notification"
  alarm_level                 = "CRITICAL"
  trigger_condition_count     = 2
  trigger_condition_frequency = 3
  send_recovery_notifications = true
  recovery_frequency          = 4
  notification_frequency      = 15
  status                      = "RUNNING"
  
  keywords_requests {
    keywords               = "%[2]s_key_words"
    condition              = "<"
    number                 = 50
    log_group_id           = huaweicloud_lts_group.test[0].id
    log_stream_id          = huaweicloud_lts_stream.test[0].id
    search_time_range_unit = "hour"
    search_time_range      = 1
  }

  frequency {
    type            = "CRON"
    cron_expression = "0 18 * * *"
  }

  notification_save_rule {
    template_name = local.template_name
    user_name     = huaweicloud_identity_user.test.name
    language      = "en-us"

    # Verify the 'keywords_alarm_send_code' is '0' in the modification logic.
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
`, testAlarmRule_base(name, password), name, acceptance.HW_LTS_ALARM_ACTION_RULE_NAME)
}

func testKeywordsAlarmRule_step4(name, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_keywords_alarm_rule" "test" {
  name                        = "%[2]s"
  alarm_level                 = "INFO"
  send_notifications          = true
  trigger_condition_count     = 2
  trigger_condition_frequency = 3
  send_recovery_notifications = true
  recovery_frequency          = 4
  notification_frequency      = 15
  alarm_action_rule_name      = "%[3]s"

  keywords_requests {
    keywords               = "%[2]s_key_words_update"
    condition              = ">="
    number                 = 50
    log_group_id           = huaweicloud_lts_group.test[1].id
    log_stream_id          = huaweicloud_lts_stream.test[1].id
    search_time_range_unit = "hour"
    search_time_range      = 1
  }

  frequency {
    type        = "WEEKLY"
    day_of_week = 6
    hour_of_day = 10
  }
}

resource "huaweicloud_lts_keywords_alarm_rule" "with_notification_save_rule" {
  name                        = "%[2]s_notification"
  alarm_level                 = "CRITICAL"
  trigger_condition_count     = 2
  trigger_condition_frequency = 3
  recovery_frequency          = 4
  notification_frequency      = 15
  status                      = "RUNNING"
  
  keywords_requests {
    keywords               = "%[2]s_key_words"
    condition              = "<"
    number                 = 50
    log_group_id           = huaweicloud_lts_group.test[0].id
    log_stream_id          = huaweicloud_lts_stream.test[0].id
    search_time_range_unit = "hour"
    search_time_range      = 1
  }

  frequency {
    type            = "CRON"
    cron_expression = "0 18 * * *"
  }

  # Delete 'notification_save_rule' and 'send_recovery_notifications' parameter.
  # Verify the 'keywords_alarm_send_code' is '3' in the modification logic.
}
`, testAlarmRule_base(name, password), name, acceptance.HW_LTS_ALARM_ACTION_RULE_NAME)
}
