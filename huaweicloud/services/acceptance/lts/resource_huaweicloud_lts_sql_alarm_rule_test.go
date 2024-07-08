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
		obj      interface{}
		name     = acceptance.RandomAccResourceName()
		rName    = "huaweicloud_lts_sql_alarm_rule.test"
		password = acceptance.RandomPassword()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSQLAlarmRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSQLAlarmRule_basic(name, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(rName, "condition_expression", "t>0"),
					resource.TestCheckResourceAttr(rName, "alarm_level", "CRITICAL"),
					resource.TestCheckResourceAttr(rName, "status", "STOPPING"),
					resource.TestCheckResourceAttr(rName, "frequency.0.type", "HOURLY"),
					resource.TestCheckResourceAttrPair(rName, "sql_requests.0.log_group_id",
						"huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "sql_requests.0.log_stream_id",
						"huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(rName, "notification_rule.0.language", "en-us"),
					resource.TestCheckResourceAttrPair(rName, "notification_rule.0.template_name",
						"huaweicloud_lts_notification_template.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "notification_rule.0.user_name",
						"huaweicloud_identity_user.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "notification_rule.0.topics.0.name",
						"huaweicloud_smn_topic.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "notification_rule.0.topics.0.topic_urn",
						"huaweicloud_smn_topic.test", "topic_urn"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				Config: testSQLAlarmRule_basic_update(name, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "condition_expression", "t>2"),
					resource.TestCheckResourceAttr(rName, "alarm_level", "INFO"),
					resource.TestCheckResourceAttr(rName, "status", "RUNNING"),
					resource.TestCheckResourceAttr(rName, "frequency.0.type", "DAILY"),
					resource.TestCheckResourceAttr(rName, "frequency.0.hour_of_day", "6"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
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
		},
	})
}

func testSQLAlarmRule_basic(name, password string) string {
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
    log_group_id           = huaweicloud_lts_group.test.id
    log_stream_id          = huaweicloud_lts_stream.test.id
    search_time_range_unit = "minute"
    search_time_range      = 5
  }

  frequency {
    type = "HOURLY"
  }

  notification_rule {
    template_name = huaweicloud_lts_notification_template.test.name
    user_name     = huaweicloud_identity_user.test.name
    language      = "en-us"

    topics {
      name      = huaweicloud_smn_topic.test.name
      topic_urn = huaweicloud_smn_topic.test.topic_urn
    }
  }
}
`, testAlarmRule_base(name, password), name)
}

func testSQLAlarmRule_basic_update(name, password string) string {
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
    log_group_id           = huaweicloud_lts_group.test.id
    log_stream_id          = huaweicloud_lts_stream.test.id
    search_time_range_unit = "minute"
    search_time_range      = 5
  }

  frequency {
    type        = "DAILY"
    hour_of_day = 6
  }

  notification_rule {
    template_name = huaweicloud_lts_notification_template.test.name
    user_name     = huaweicloud_identity_user.test.name
    language      = "en-us"

    topics {
      name      = huaweicloud_smn_topic.test.name
      topic_urn = huaweicloud_smn_topic.test.topic_urn
    }
  }
}
`, testAlarmRule_base(name, password), name)
}
