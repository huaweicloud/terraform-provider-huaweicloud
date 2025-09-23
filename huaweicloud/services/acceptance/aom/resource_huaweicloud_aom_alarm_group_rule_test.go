package aom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/aom"
)

func getAlarmGroupRuleResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("aom", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM client: %s", err)
	}

	rule, err := aom.GetAlarmGroupRule(client, state.Primary.Attributes["name"])
	if err != nil {
		return nil, err
	}

	return rule, nil
}

func TestAccAlarmGroupRule_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_aom_alarm_group_rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getAlarmGroupRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAlarmGroupRule_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "group_by.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "detail.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "group_interval", "5"),
					resource.TestCheckResourceAttr(resourceName, "group_repeat_waiting", "0"),
					resource.TestCheckResourceAttr(resourceName, "group_wait", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				Config: testAlarmGroupRule_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "group_by.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "detail.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "group_interval", "120"),
					resource.TestCheckResourceAttr(resourceName, "group_repeat_waiting", "180"),
					resource.TestCheckResourceAttr(resourceName, "group_wait", "60"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAlarmGroupRule_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_aom_alarm_group_rule" "test" {
  name                  = "%[2]s"
  group_by              = ["event_severity", "resource_provider", "key"]
  group_interval        = 5
  group_repeat_waiting  = 0
  group_wait            = 0
  description           = "test"
  enterprise_project_id = "%[3]s"

  detail {
    bind_notification_rule_ids = [huaweicloud_aom_alarm_action_rule.test.name]

    match {
      key     = "event_severity"
      operate = "EQUALS"
      value   = ["Critical", "Major"]
    }

    match {
      key     = "resource_type"
      operate = "EXIST"
    }
  }

  detail {
    bind_notification_rule_ids = [huaweicloud_aom_alarm_action_rule.test.name]

    match {
      key     = "notification_scene"
      operate = "EQUALS"
      value   = ["notify_triggered"]
    }

    match {
      key     = "key"
      operate = "EQUALS"
      value   = ["value"]
    }
  }
}`, testAlarmActionRule_basic(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAlarmGroupRule_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_aom_alarm_group_rule" "test" {
  name                  = "%[2]s"
  group_by              = ["resource_provider"]
  group_wait            = 60
  group_interval        = 120
  group_repeat_waiting  = 180
  enterprise_project_id = "%[3]s"

  detail {
    bind_notification_rule_ids = [huaweicloud_aom_alarm_action_rule.test.name]

    match {
      key     = "event_severity"
      operate = "EQUALS"
      value   = ["Critical", "Major"]
    }
  }
}`, testAlarmActionRule_basic(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
