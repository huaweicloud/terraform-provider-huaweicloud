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

func getAlarmActionRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getAlarmActionRule: Query the Alarm Action Rule
	var (
		getAlarmActionRuleHttpUrl = "v2/{project_id}/alert/action-rules/{id}"
		getAlarmActionRuleProduct = "aom"
	)
	getAlarmActionRuleClient, err := cfg.NewServiceClient(getAlarmActionRuleProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM Client: %s", err)
	}

	getAlarmActionRulePath := getAlarmActionRuleClient.Endpoint + getAlarmActionRuleHttpUrl
	getAlarmActionRulePath = strings.ReplaceAll(getAlarmActionRulePath, "{project_id}", getAlarmActionRuleClient.ProjectID)
	getAlarmActionRulePath = strings.ReplaceAll(getAlarmActionRulePath, "{id}", state.Primary.ID)

	getAlarmActionRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getAlarmActionRuleOpt.MoreHeaders = map[string]string{
		"Content-Type": "application/json",
	}
	getAlarmActionRuleResp, err := getAlarmActionRuleClient.Request("GET", getAlarmActionRulePath, &getAlarmActionRuleOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving AlarmActionRule: %s", err)
	}
	return utils.FlattenResponse(getAlarmActionRuleResp)
}

func TestAccAlarmActionRule_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_aom_alarm_action_rule.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAlarmActionRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAlarmActionRule_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
					resource.TestCheckResourceAttr(rName, "type", "1"),
					resource.TestCheckResourceAttr(rName, "notification_template", "aom.built-in.template.zh"),
					resource.TestCheckResourceAttrPair(rName, "smn_topics.0.topic_urn",
						"huaweicloud_smn_topic.topic_1", "topic_urn"),
				),
			},
			{
				Config: testAlarmActionRule_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "terraform test update"),
					resource.TestCheckResourceAttr(rName, "type", "1"),
					resource.TestCheckResourceAttr(rName, "notification_template", "aom.built-in.template.en"),
					resource.TestCheckResourceAttrPair(rName, "smn_topics.0.topic_urn",
						"huaweicloud_smn_topic.topic_1", "topic_urn"),
					resource.TestCheckResourceAttrPair(rName, "smn_topics.1.topic_urn",
						"huaweicloud_smn_topic.topic_2", "topic_urn"),
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

func testAlarmActionRule_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "topic_1" {
  name = "%[1]s_1"
}

resource "huaweicloud_aom_alarm_action_rule" "test" {
  name                  = "%[1]s"
  description           = "terraform test"
  type                  = "1"
  notification_template = "aom.built-in.template.zh"

  smn_topics {
    topic_urn = huaweicloud_smn_topic.topic_1.topic_urn
  }
}
`, name)
}

func testAlarmActionRule_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "topic_1" {
  name = "%[1]s_1"
}

resource "huaweicloud_smn_topic" "topic_2" {
  name = "%[1]s_2"
}

resource "huaweicloud_aom_alarm_action_rule" "test" {
  name                  = "%[1]s"
  description           = "terraform test update"
  type                  = "1"
  notification_template = "aom.built-in.template.en"

  smn_topics {
    topic_urn = huaweicloud_smn_topic.topic_1.topic_urn
  }

  smn_topics {
    topic_urn = huaweicloud_smn_topic.topic_2.topic_urn
  }
}
`, name)
}
