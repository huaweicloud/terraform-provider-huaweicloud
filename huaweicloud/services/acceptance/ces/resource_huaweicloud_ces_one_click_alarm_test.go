package ces

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ces"
)

func getOneClickAlarmFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var getOneClickAlarmProduct = "ces"
	client, err := cfg.NewServiceClient(getOneClickAlarmProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CES client: %s", err)
	}

	return ces.GetOneClickAlarm(client, state.Primary.ID)
}

func TestAccOneClickAlarm_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_ces_one_click_alarm.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getOneClickAlarmFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testOneClickAlarm_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "one_click_alarm_id", "OBSSystemOneClickAlarm"),
					resource.TestCheckResourceAttr(rName, "dimension_names.0.metric.0", "bucket_name"),
					resource.TestCheckResourceAttr(rName, "dimension_names.0.event", "true"),
					resource.TestCheckResourceAttr(rName, "notification_enabled", "true"),
					resource.TestCheckResourceAttrSet(rName, "namespace"),
					resource.TestCheckResourceAttrSet(rName, "description"),
					resource.TestCheckResourceAttrSet(rName, "enabled"),
				),
			},
			{
				Config: testOneClickAlarm_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "one_click_alarm_id", "OBSSystemOneClickAlarm"),
					resource.TestCheckResourceAttr(rName, "notification_enabled", "false"),
					resource.TestCheckResourceAttrSet(rName, "namespace"),
					resource.TestCheckResourceAttrSet(rName, "description"),
					resource.TestCheckResourceAttrSet(rName, "enabled"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"one_click_alarm_id", "dimension_names", "notification_enabled",
					"alarm_notifications", "ok_notifications", "notification_begin_time", "notification_end_time",
				},
			},
		},
	})
}

func testOneClickAlarm_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ces_one_click_alarm" test {
  one_click_alarm_id = "OBSSystemOneClickAlarm"

  dimension_names {
    metric = ["bucket_name"]
    event  = true
  }

  notification_enabled = true

  alarm_notifications {
    type = "contact"

    notification_list = [
      huaweicloud_smn_topic.test.topic_urn
    ]
  }
  
  ok_notifications {
    type = "contact"

    notification_list = [
      huaweicloud_smn_topic.test.topic_urn
    ]
  }
  
  notification_begin_time = "00:00"
  notification_end_time   = "23:59"
}
`, testCESAlarmRule_topicBase(name))
}

func testOneClickAlarm_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ces_one_click_alarm" test {
  one_click_alarm_id = "OBSSystemOneClickAlarm"
  
  dimension_names {
    metric = ["bucket_name"]
    event  = true
  }
  
  notification_enabled    = false
  notification_begin_time = "00:00"
  notification_end_time   = "23:59"
}
`, testCESAlarmRule_topicBase(name))
}
