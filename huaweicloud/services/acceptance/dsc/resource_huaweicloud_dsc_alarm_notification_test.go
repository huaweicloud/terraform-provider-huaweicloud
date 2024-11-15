package dsc

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

func getResourceDscAlarmNotificationFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dsc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + "v1/{project_id}/sdg/smn/topics"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DSC alarm notification: %s", err)
	}

	return utils.FlattenResponse(resp)
}

func TestAccResourceDscAlarmNotification_basic(t *testing.T) {
	var (
		obj          interface{}
		randName     = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_dsc_alarm_notification.test"
		alarmTopicID = acceptance.HW_DSC_ALARM_TOPIC_ID
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceDscAlarmNotificationFunc,
	)

	// Avoid CheckDestroy
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please create a DSC instance in advance.
			acceptance.TestAccPrecheckDscInstance(t)
			// Please configure the alarm topic ID into the environment variable in advance.
			// Currently, it can only be obtained through F12 on the console.
			acceptance.TestAccPrecheckDscAlarmTopicID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceDscAlarmNotification_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "alarm_topic_id", alarmTopicID),
					resource.TestCheckResourceAttr(resourceName, "status", "0"),
					resource.TestCheckResourceAttrPair(resourceName, "topic_urn",
						"huaweicloud_smn_topic.test1", "topic_urn"),
				),
			},
			{
				Config: testResourceDscAlarmNotification_update(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "alarm_topic_id", alarmTopicID),
					resource.TestCheckResourceAttr(resourceName, "status", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "topic_urn",
						"huaweicloud_smn_topic.test2", "topic_urn"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status"},
			},
		},
	})
}

func testResourceDscAlarmNotification_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test1" {
  name         = "smn-%[1]s-1"
  display_name = "The display name of smn topic 1"
}

resource "huaweicloud_smn_topic" "test2" {
  name         = "smn-%[1]s-2"
  display_name = "The display name of smn topic 2"
}
`, name)
}

func testResourceDscAlarmNotification_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dsc_alarm_notification" "test" {
  alarm_topic_id = "%[2]s"
  status         = 0
  topic_urn      = huaweicloud_smn_topic.test1.topic_urn
}
`, testResourceDscAlarmNotification_base(name), acceptance.HW_DSC_ALARM_TOPIC_ID)
}

func testResourceDscAlarmNotification_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dsc_alarm_notification" "test" {
  alarm_topic_id = "%[2]s"
  status         = 1
  topic_urn      = huaweicloud_smn_topic.test2.topic_urn
}
`, testResourceDscAlarmNotification_base(name), acceptance.HW_DSC_ALARM_TOPIC_ID)
}
