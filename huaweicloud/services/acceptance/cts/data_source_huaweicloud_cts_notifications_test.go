package cts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceCTSNotifications_basic(t *testing.T) {
	defaultDataSourceName := "data.huaweicloud_cts_notifications.filter_by_name"
	dc := acceptance.InitDataSourceCheck(defaultDataSourceName)
	name := acceptance.RandomAccResourceName()
	baseConfig := testAccDatasourceCTSNotifications_base(name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceCTSNotifications_basic(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(defaultDataSourceName, "notifications.0.agency_name", "cts_admin_trust"),
					resource.TestCheckResourceAttrSet(defaultDataSourceName, "notifications.0.name"),
					resource.TestCheckResourceAttrSet(defaultDataSourceName, "notifications.0.operation_type"),
					resource.TestCheckResourceAttrSet(defaultDataSourceName, "notifications.0.status"),
					resource.TestCheckResourceAttrSet(defaultDataSourceName, "notifications.0.id"),
					resource.TestCheckResourceAttrSet(defaultDataSourceName, "notifications.0.topic_id"),
					resource.TestCheckOutput("is_default_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_topic_id_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_operation_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceCTSNotifications_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "topic_1" {
  name = "%[1]s"
}
  
resource "huaweicloud_cts_notification" "notify" {
  name           = "%[1]s_1"
  operation_type = "complete"
  smn_topic      = huaweicloud_smn_topic.topic_1.id
  agency_name    = "cts_admin_trust"
  
  filter {
    condition = "AND"
    rule      = ["code = 200","resource_name = test"]
  }
}

resource "huaweicloud_cts_notification" "test" {
  name           = "%[1]s_2"
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
}
`, name)
}

func testAccDatasourceCTSNotifications_basic(config string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cts_notifications" "test" {
  type = "smn"
  
  depends_on = [
    huaweicloud_cts_notification.test,
    huaweicloud_cts_notification.notify
  ]
}

output "is_default_useful" {
  value = length(data.huaweicloud_cts_notifications.test.notifications) >= 2
}

locals {
  name = huaweicloud_cts_notification.notify.name
}
  
data "huaweicloud_cts_notifications" "filter_by_name" {
  type = "smn"
  name = local.name

  depends_on = [
    huaweicloud_cts_notification.test,
    huaweicloud_cts_notification.notify
  ]
}
	
output "is_name_filter_useful" {
  value = length(data.huaweicloud_cts_notifications.filter_by_name.notifications) == 1 && alltrue(
    [for v in data.huaweicloud_cts_notifications.filter_by_name.notifications[*].name : v == local.name]
  )
}

locals {
  topic_id = huaweicloud_cts_notification.notify.smn_topic
}
  
data "huaweicloud_cts_notifications" "filter_by_topic_id" {
  type     = "smn"
  topic_id = local.topic_id

  depends_on = [
    huaweicloud_cts_notification.test,
    huaweicloud_cts_notification.notify
  ]
}
	
output "is_topic_id_filter_useful" {
  value = length(data.huaweicloud_cts_notifications.filter_by_topic_id.notifications) > 0 && alltrue(
    [for v in data.huaweicloud_cts_notifications.filter_by_topic_id.notifications[*].topic_id : v == local.topic_id]
  )
}

locals {
  status = huaweicloud_cts_notification.notify.status
}

data "huaweicloud_cts_notifications" "filter_by_status" {
  type   = "smn"
  status = local.status

  depends_on = [
    huaweicloud_cts_notification.test,
    huaweicloud_cts_notification.notify
  ]
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_cts_notifications.filter_by_status.notifications) >= 1 && alltrue(
    [for v in data.huaweicloud_cts_notifications.filter_by_status.notifications[*].status : v == local.status]
  )
}

locals {
  operation_type = huaweicloud_cts_notification.notify.operation_type
}

data "huaweicloud_cts_notifications" "filter_by_operation_type" {
  type           = "smn"
  operation_type = local.operation_type

  depends_on = [
    huaweicloud_cts_notification.test,
    huaweicloud_cts_notification.notify,
  ]
}

output "is_operation_type_filter_useful" {
  value = length(data.huaweicloud_cts_notifications.filter_by_operation_type.notifications) >= 1 && alltrue(
    [for v in data.huaweicloud_cts_notifications.filter_by_operation_type.notifications[*].operation_type : v == local.operation_type]
  )
}
`, config)
}
