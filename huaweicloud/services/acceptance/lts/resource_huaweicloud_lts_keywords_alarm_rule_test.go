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
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_lts_keywords_alarm_rule.test"

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
				Config: testKeywordsAlarmRule_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(rName, "alarm_level", "CRITICAL"),
					resource.TestCheckResourceAttr(rName, "status", "RUNNING"),
					resource.TestCheckResourceAttr(rName, "frequency.0.type", "HOURLY"),
					resource.TestCheckResourceAttrPair(rName, "keywords_requests.0.log_group_id",
						"huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "keywords_requests.0.log_stream_id",
						"huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				Config: testKeywordsAlarmRule_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "alarm_level", "INFO"),
					resource.TestCheckResourceAttr(rName, "status", "STOPPING"),
					resource.TestCheckResourceAttr(rName, "frequency.0.type", "DAILY"),
					resource.TestCheckResourceAttr(rName, "frequency.0.hour_of_day", "6"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testKeywordsAlarmRule_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_keywords_alarm_rule" "test" {
  name        = "%[2]s"
  description = "created by terraform"
  alarm_level = "CRITICAL"
  
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
}
`, testAccLtsStream_basic(name), name)
}

func testKeywordsAlarmRule_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_keywords_alarm_rule" "test" {
  name        = "%[2]s"
  description = ""
  alarm_level = "INFO"
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
    type        = "DAILY"
    hour_of_day = 6
  }
}
`, testAccLtsStream_basic(name), name)
}
