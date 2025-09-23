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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/aom"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getAlarmSilenceRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getAlarmSilenceRule: Query the Alarm Silence Rule
	var (
		getAlarmSilenceRuleHttpUrl = "v2/{project_id}/alert/mute-rules"
		getAlarmSilenceRuleProduct = "aom"
	)
	getAlarmSilenceRuleClient, err := cfg.NewServiceClient(getAlarmSilenceRuleProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM Client: %s", err)
	}

	getAlarmSilenceRulePath := getAlarmSilenceRuleClient.Endpoint + getAlarmSilenceRuleHttpUrl
	getAlarmSilenceRulePath = strings.ReplaceAll(getAlarmSilenceRulePath, "{project_id}", getAlarmSilenceRuleClient.ProjectID)

	getAlarmSilenceRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getAlarmSilenceRuleOpt.MoreHeaders = map[string]string{
		"Content-Type": "application/json",
	}
	getAlarmSilenceRuleResp, err := getAlarmSilenceRuleClient.Request("GET", getAlarmSilenceRulePath, &getAlarmSilenceRuleOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving AlarmSilenceRule: %s", err)
	}

	getAlarmSilenceRuleRespBody, err := utils.FlattenResponse(getAlarmSilenceRuleResp)
	if err != nil {
		return nil, err
	}

	rules := aom.FilterListAlarmSilenceRules(getAlarmSilenceRuleRespBody.([]interface{}), state.Primary.ID)
	if len(rules) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}

	return rules[0], nil
}

func TestAccAlarmSilenceRule_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_aom_alarm_silence_rule.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAlarmSilenceRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAlarmSilenceRule_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
					resource.TestCheckResourceAttr(rName, "time_zone", "Asia/Shanghai"),
					resource.TestCheckResourceAttr(rName, "silence_time.0.type", "DAILY"),
					resource.TestCheckResourceAttr(rName, "silence_time.0.starts_at", "0"),
					resource.TestCheckResourceAttr(rName, "silence_time.0.ends_at", "86399"),
					resource.TestCheckResourceAttr(rName, "silence_conditions.0.conditions.0.key", "event_severity"),
					resource.TestCheckResourceAttr(rName, "silence_conditions.0.conditions.0.operate", "EQUALS"),
				),
			},
			{
				Config: testAlarmSilenceRule_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "terraform test update"),
					resource.TestCheckResourceAttr(rName, "time_zone", "Asia/Shanghai"),
					resource.TestCheckResourceAttr(rName, "silence_time.0.type", "WEEKLY"),
					resource.TestCheckResourceAttr(rName, "silence_time.0.starts_at", "64800"),
					resource.TestCheckResourceAttr(rName, "silence_time.0.ends_at", "86399"),
					resource.TestCheckResourceAttr(rName, "silence_conditions.0.conditions.0.key", "event_severity"),
					resource.TestCheckResourceAttr(rName, "silence_conditions.0.conditions.0.operate", "EXIST"),
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

func testAlarmSilenceRule_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_alarm_silence_rule" "test" {
  name        = "%s"
  description = "terraform test"
  time_zone   = "Asia/Shanghai"

  silence_time {
    type      = "DAILY"
    starts_at = 0
    ends_at   = 86399
  }

  silence_conditions {
    conditions {
      key     = "event_severity"
      operate = "EQUALS"
      value   = ["Info"]
    }
  }
}
`, name)
}

func testAlarmSilenceRule_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_alarm_silence_rule" "test" {
  name        = "%s"
  description = "terraform test update"
  time_zone   = "Asia/Shanghai"

  silence_time {
    type      = "WEEKLY"
    starts_at = 64800
    ends_at   = 86399
    scope     = [1, 2, 3, 4, 5]   
  }

  silence_conditions {
    conditions {
      key     = "event_severity"
      operate = "EXIST"
    }
  }
}
`, name)
}
