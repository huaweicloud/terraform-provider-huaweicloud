package cts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	cts "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cts/v3/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getCTSNotificationResourceObj(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcCtsV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CTS client: %s", err)
	}

	notificationName := state.Primary.ID
	notificationType := cts.GetListNotificationsRequestNotificationTypeEnum().SMN
	listOpts := &cts.ListNotificationsRequest{
		NotificationType: notificationType,
		NotificationName: &notificationName,
	}

	response, err := client.ListNotifications(listOpts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CTS notification: %s", err)
	}

	if response.Notifications == nil || len(*response.Notifications) == 0 {
		return nil, fmt.Errorf("can not find the CTS notification %s", notificationName)
	}

	allNotifications := *response.Notifications
	ctsNotification := allNotifications[0]

	return ctsNotification, nil
}

func TestAccCTSNotification_basic(t *testing.T) {
	var notify cts.NotificationsResponseBody
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_cts_notification.notify"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&notify,
		getCTSNotificationResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCTSNotification_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "operation_type", "complete"),
					resource.TestCheckResourceAttr(resourceName, "status", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.condition", "AND"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.rule.#", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "smn_topic",
						"huaweicloud_smn_topic.topic_1", "id"),
				),
			},
			{
				Config: testAccCTSNotification_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "operation_type", "customized"),
					resource.TestCheckResourceAttr(resourceName, "status", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.condition", "OR"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.rule.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "operations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "operations.0.service", "ECS"),
					resource.TestCheckResourceAttr(resourceName, "operation_users.0.group", "devops"),
					resource.TestCheckResourceAttr(resourceName, "operation_users.0.users.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "operation_users.0.users.0", "tf-user10"),
					resource.TestCheckResourceAttrPair(resourceName, "smn_topic",
						"huaweicloud_smn_topic.topic_1", "id"),
				),
			},
			{
				Config: testAccCTSNotification_disable(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "operation_type", "customized"),
					resource.TestCheckResourceAttr(resourceName, "status", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "operations.#", "1"),
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

func testAccCTSNotification_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "topic_1" {
  name  = "%[1]s"
}

resource "huaweicloud_cts_notification" "notify" {
  name           = "%[1]s"
  operation_type = "complete"
  smn_topic      = huaweicloud_smn_topic.topic_1.id

  filter {
    condition = "AND"
    rule      = ["code = 200","resource_name = test"]
  }
}
`, rName)
}

func testAccCTSNotification_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "topic_1" {
  name  = "%[1]s"
}

resource "huaweicloud_cts_notification" "notify" {
  name           = "%[1]s"
  operation_type = "customized"
  smn_topic      = huaweicloud_smn_topic.topic_1.id

  filter {
    condition = "OR"
    rule      = ["code = 400","resource_name = name","api_version = 1.0"]
  }

  operations {
    service     = "ECS"
    resource    = "ecs"
    trace_names = ["createServer", "deleteServer"]
  }

  operation_users {
    group = "devops"
    users = ["tf-user10"]
  }
}
`, rName)
}

func testAccCTSNotification_disable(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "topic_1" {
  name  = "%[1]s"
}

resource "huaweicloud_cts_notification" "notify" {
  name           = "%[1]s"
  operation_type = "customized"
  smn_topic      = huaweicloud_smn_topic.topic_1.id
  enabled        = false

  operations {
    service     = "ECS"
    resource    = "ecs"
    trace_names = ["createServer", "deleteServer"]
  }
}
`, rName)
}
