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

func getKeywordsAlarmRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getKeywordsAlarmRule: Query the LTS KeywordsAlarmRule detail
	var (
		getKeywordsAlarmRuleHttpUrl = "v2/{project_id}/lts/alarms/keywords-alarm-rule"
		getKeywordsAlarmRuleProduct = "lts"
	)
	getKeywordsAlarmRuleClient, err := cfg.NewServiceClient(getKeywordsAlarmRuleProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}

	getKeywordsAlarmRulePath := getKeywordsAlarmRuleClient.Endpoint + getKeywordsAlarmRuleHttpUrl
	getKeywordsAlarmRulePath = strings.ReplaceAll(getKeywordsAlarmRulePath, "{project_id}", getKeywordsAlarmRuleClient.ProjectID)

	getKeywordsAlarmRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getKeywordsAlarmRuleResp, err := getKeywordsAlarmRuleClient.Request("GET", getKeywordsAlarmRulePath, &getKeywordsAlarmRuleOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Keywords alarm rule: %s", err)
	}

	getKeywordsAlarmRuleRespBody, err := utils.FlattenResponse(getKeywordsAlarmRuleResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Keywords alarm rule: %s", err)
	}

	jsonPath := fmt.Sprintf("keywords_alarm_rules[?keywords_alarm_rule_id =='%s']|[0]", state.Primary.ID)
	getKeywordsAlarmRuleRespBody = utils.PathSearch(jsonPath, getKeywordsAlarmRuleRespBody, nil)
	if getKeywordsAlarmRuleRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return getKeywordsAlarmRuleRespBody, nil
}

func TestAccKeywordsAlarmRule_basic(t *testing.T) {
	var (
		obj      interface{}
		name     = acceptance.RandomAccResourceName()
		rName    = "huaweicloud_lts_keywords_alarm_rule.test"
		password = acceptance.RandomPassword()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getKeywordsAlarmRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeywordsAlarmRule_step1(name, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(rName, "alarm_level", "CRITICAL"),
					resource.TestCheckResourceAttr(rName, "status", "STOPPING"),
					resource.TestCheckResourceAttr(rName, "frequency.0.type", "HOURLY"),
					resource.TestCheckResourceAttrPair(rName, "keywords_requests.0.log_group_id",
						"huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "keywords_requests.0.log_stream_id",
						"huaweicloud_lts_stream.test", "id"),
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
				Config: testKeywordsAlarmRule_step2(name, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", ""),
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

func testAlarmRule_base(name, password string) string {
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

resource "huaweicloud_lts_notification_template" "test" {
  name   = "%[1]s"
  source = "LTS"
  locale = "en-us"
  
  templates {
    sub_type = "sms"
    content  = "This body content of template."
  }
}

resource "huaweicloud_identity_user" "test" {
  name     = "%[1]s"
  enabled  = true
  password = "%[2]s"
}
`, name, password)
}

func testKeywordsAlarmRule_step1(name, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_keywords_alarm_rule" "test" {
  name        = "%[2]s"
  description = "created by terraform"
  alarm_level = "CRITICAL"
  status      = "STOPPING"

  keywords_requests {
    keywords               = "%[2]s_key_words"
    condition              = ">"
    number                 = 100
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

    topics {
      name      = huaweicloud_smn_topic.test.name
      topic_urn = huaweicloud_smn_topic.test.topic_urn
    }
  }
}
`, testAlarmRule_base(name, password), name)
}

func testKeywordsAlarmRule_step2(name, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_keywords_alarm_rule" "test" {
  name        = "%[2]s"
  description = ""
  alarm_level = "INFO"
  status      = "RUNNING"
  
  keywords_requests {
    keywords               = "%[2]s_key_words"
    condition              = ">"
    number                 = 100
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

    topics {
      name      = huaweicloud_smn_topic.test.name
      topic_urn = huaweicloud_smn_topic.test.topic_urn
    }
  }
}
`, testAlarmRule_base(name, password), name)
}
