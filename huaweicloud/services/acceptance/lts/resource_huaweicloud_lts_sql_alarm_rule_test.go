package lts

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getSQLAlarmRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getSQLAlarmRule: Query the LTS SQLAlarmRule detail
	var (
		getSQLAlarmRuleHttpUrl = "v2/{project_id}/lts/alarms/sql-alarm-rule"
		getSQLAlarmRuleProduct = "lts"
	)
	getSQLAlarmRuleClient, err := cfg.NewServiceClient(getSQLAlarmRuleProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}

	getSQLAlarmRulePath := getSQLAlarmRuleClient.Endpoint + getSQLAlarmRuleHttpUrl
	getSQLAlarmRulePath = strings.ReplaceAll(getSQLAlarmRulePath, "{project_id}", getSQLAlarmRuleClient.ProjectID)

	getSQLAlarmRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getSQLAlarmRuleResp, err := getSQLAlarmRuleClient.Request("GET", getSQLAlarmRulePath, &getSQLAlarmRuleOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SQL alarm rule: %s", err)
	}

	getSQLAlarmRuleRespBody, err := utils.FlattenResponse(getSQLAlarmRuleResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SQL alarm rule: %s", err)
	}

	jsonPath := fmt.Sprintf("sql_alarm_rules[?sql_alarm_rule_id =='%s']|[0]", state.Primary.ID)
	getSQLAlarmRuleRespBody = utils.PathSearch(jsonPath, getSQLAlarmRuleRespBody, nil)
	if getSQLAlarmRuleRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return getSQLAlarmRuleRespBody, nil
}

func TestAccSQLAlarmRule_basic(t *testing.T) {
	var (
		name     = acceptance.RandomAccResourceName()
		password = acceptance.RandomPassword()

		obj                interface{}
		rName              = "huaweicloud_lts_sql_alarm_rule.test"
		withNotiSaveRule   = "huaweicloud_lts_sql_alarm_rule.with_notification_save_rule"
		rc                 = acceptance.InitResourceCheck(rName, &obj, getSQLAlarmRuleResourceFunc)
		rcWithNotiSaveRule = acceptance.InitResourceCheck(withNotiSaveRule, &obj, getSQLAlarmRuleResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rc.CheckResourceDestroy(),
			rcWithNotiSaveRule.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccSQLAlarmRule_basic_step1(name, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttrPair(rName, "sql_requests.0.log_group_id",
						"huaweicloud_lts_group.test.0", "id"),
					resource.TestCheckResourceAttrPair(rName, "sql_requests.0.log_stream_id",
						"huaweicloud_lts_stream.test.0", "id"),
					resource.TestCheckResourceAttr(rName, "frequency.0.type", "HOURLY"),
					resource.TestCheckResourceAttr(rName, "condition_expression", "t>0"),
					resource.TestCheckResourceAttr(rName, "alarm_level", "CRITICAL"),
					// Check optional parameters.
					resource.TestCheckResourceAttr(rName, "status", "STOPPING"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),

					rcWithNotiSaveRule.CheckResourceExists(),
					resource.TestCheckResourceAttr(withNotiSaveRule, "send_notifications", "true"),
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
				Config: testAccSQLAlarmRule_basic_step2(name, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "condition_expression", "t>2"),
					resource.TestCheckResourceAttr(rName, "alarm_level", "INFO"),
					resource.TestCheckResourceAttr(rName, "frequency.0.type", "DAILY"),
					resource.TestCheckResourceAttr(rName, "frequency.0.hour_of_day", "6"),
					// Check Optional parameter.
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
					resource.TestCheckResourceAttr(rName, "notification_save_rule.0.template_name", "sql_template"),
					resource.TestCheckResourceAttr(withNotiSaveRule, "notification_save_rule.0.topics.#", "2"),
				),
			},
			{
				Config: testAccSQLAlarmRule_basic_step3(name, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					rcWithNotiSaveRule.CheckResourceExists(),
					resource.TestCheckResourceAttr(withNotiSaveRule, "send_notifications", "false"),
					resource.TestCheckResourceAttr(rName, "notification_save_rule.#", "0"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"notification_rule",
				},
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

func testAccSQLAlarmRule_basic_step1(name, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_sql_alarm_rule" "test" {
  name                 = "%[2]s"
  description          = "created by terraform"
  condition_expression = "t>0"
  alarm_level          = "CRITICAL"
  status               = "STOPPING"

  sql_requests {
    title                  = "%[2]s_sql"
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
  name                 = "%[2]s_notification"
  condition_expression = "t>0"
  alarm_level          = "MINOR"
  send_notifications   = true
  
  sql_requests {
    title                  = "%[2]s_sql"
    sql                    = "select count(*) as t"
    log_group_id           = huaweicloud_lts_group.test[0].id
    log_stream_id          = huaweicloud_lts_stream.test[0].id
    search_time_range_unit = "minute"
    search_time_range      = 5
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
`, testAccSQLAlarmRule_base(name, password), name)
}

func testAccSQLAlarmRule_basic_step2(name, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_sql_alarm_rule" "test" {
  name                 = "%[2]s"
  condition_expression = "t>2"
  alarm_level          = "INFO"
  send_notifications   = true
  status               = "RUNNING"

  sql_requests {
    title                  = "%[2]s_sql"
    sql                    = "select count(*) as t"
    log_group_id           = huaweicloud_lts_group.test[0].id
    log_stream_id          = huaweicloud_lts_stream.test[0].id
    search_time_range_unit = "minute"
    search_time_range      = 5
  }

  frequency {
    type        = "DAILY"
    hour_of_day = 6
  }

  # Verify the 'keywords_alarm_send_code' is '2' in the modification logic.
  notification_save_rule {
    template_name = local.sql_template_name
    user_name     = huaweicloud_identity_user.test.name
    language      = "en-us"

	# Verify the 'keywords_alarm_send_code' is '2' in the modification logic.
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
  
  sql_requests {
    title                  = "%[2]s_sql"
    sql                    = "select count(*) as t"
    log_group_id           = huaweicloud_lts_group.test[0].id
    log_stream_id          = huaweicloud_lts_stream.test[0].id
    search_time_range_unit = "minute"
    search_time_range      = 5
  }

  frequency {
    type = "HOURLY"
  }

  notification_save_rule {
    template_name = local.sql_template_name
    user_name     = huaweicloud_identity_user.test.name
    language      = "en-us"

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
`, testAccSQLAlarmRule_base(name, password), name)
}

func testAccSQLAlarmRule_basic_step3(name, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_sql_alarm_rule" "test" {
  name                 = "%[2]s"
  condition_expression = "t>2"
  alarm_level          = "INFO"
  status               = "RUNNING"

  sql_requests {
    title                  = "%[2]s_sql"
    sql                    = "select count(*) as t"
    log_group_id           = huaweicloud_lts_group.test[0].id
    log_stream_id          = huaweicloud_lts_stream.test[0].id
    search_time_range_unit = "minute"
    search_time_range      = 5
  }

  frequency {
    type        = "DAILY"
    hour_of_day = 6
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
    type = "HOURLY"
  }

  # Remove 'notification_save_rule', verify the 'keywords_alarm_send_code' is '3' in the modification logic.
}
`, testAccSQLAlarmRule_base(name, password), name)
}
