package aom

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

func getEventAlarmRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getEventAlarmRule: Query the Event Alarm Rule
	var (
		getEventAlarmRuleHttpUrl = "v2/{project_id}/event2alarm-rule"
		getEventAlarmRuleProduct = "aom"
	)
	getEventAlarmRuleClient, err := cfg.NewServiceClient(getEventAlarmRuleProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM Client: %s", err)
	}

	getEventAlarmRulePath := getEventAlarmRuleClient.Endpoint + getEventAlarmRuleHttpUrl
	getEventAlarmRulePath = strings.ReplaceAll(getEventAlarmRulePath, "{project_id}", getEventAlarmRuleClient.ProjectID)

	getEventAlarmRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getEventAlarmRuleOpt.MoreHeaders = map[string]string{
		"Content-Type": "application/json",
	}

	getEventAlarmRuleResp, err := getEventAlarmRuleClient.Request("GET", getEventAlarmRulePath, &getEventAlarmRuleOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving EventAlarmRule: %s", err)
	}

	getEventAlarmRuleRespBody, err := utils.FlattenResponse(getEventAlarmRuleResp)
	if err != nil {
		return nil, err
	}

	jsonPath := fmt.Sprintf("[?name=='%s']|[0]", state.Primary.ID)
	rule := utils.PathSearch(jsonPath, getEventAlarmRuleRespBody, nil)
	if rule == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return rule, nil
}

func TestAccEventAlarmRule_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_aom_event_alarm_rule.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getEventAlarmRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testEventAlarmRule_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
					resource.TestCheckResourceAttr(rName, "alarm_type", "notification"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttr(rName, "trigger_type", "accumulative"),
					resource.TestCheckResourceAttr(rName, "comparison_operator", ">="),
					resource.TestCheckResourceAttr(rName, "trigger_count", "2"),
					resource.TestCheckResourceAttr(rName, "period", "300"),
					resource.TestCheckResourceAttr(rName, "alarm_source", "AOM"),
					resource.TestCheckResourceAttr(rName, "select_object.event_type", "alarm"),
					resource.TestCheckResourceAttr(rName, "select_object.event_severity", "Critical"),
					resource.TestCheckResourceAttrPair(rName, "action_rule",
						"huaweicloud_aom_alarm_action_rule.test", "id"),
				),
			},
			{
				Config: testEventAlarmRule_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "terraform test update"),
					resource.TestCheckResourceAttr(rName, "alarm_type", "notification"),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
					resource.TestCheckResourceAttr(rName, "trigger_type", "immediately"),
					resource.TestCheckResourceAttr(rName, "alarm_source", "AOM"),
					resource.TestCheckResourceAttr(rName, "select_object.event_type", "SELECT_ALL"),
					resource.TestCheckResourceAttrPair(rName, "action_rule",
						"huaweicloud_aom_alarm_action_rule.test", "id"),
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

func testEventAlarmRule_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_aom_event_alarm_rule" "test" {
  name                = "%s"
  description         = "terraform test"
  alarm_type          = "notification"
  action_rule         = huaweicloud_aom_alarm_action_rule.test.id
  enabled             = false
  trigger_type        = "accumulative"
  period              = "300"
  comparison_operator = ">="
  trigger_count       = 2
  alarm_source        = "AOM"

  select_object = {
    "event_type"     ="alarm",
    "event_severity" = "Critical"
  }
}
`, testAlarmActionRule_basic(name), name)
}

func testEventAlarmRule_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_aom_event_alarm_rule" "test" {
  name         = "%s"
  description  = "terraform test update"
  alarm_type   = "notification"
  action_rule  = huaweicloud_aom_alarm_action_rule.test.id
  enabled      = true
  trigger_type = "immediately"
  alarm_source = "AOM"

  select_object = {
    "event_type" = "SELECT_ALL"
  }
}
`, testAlarmActionRule_basic(name), name)
}
