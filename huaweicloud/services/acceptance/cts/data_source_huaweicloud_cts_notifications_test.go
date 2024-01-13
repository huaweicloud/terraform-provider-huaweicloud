package cts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceCTSNotifications_basic(t *testing.T) {

	defaultDataSourceName := "data.huaweicloud_cts_notifications.test"
	nameFilterDataSourceName := "data.huaweicloud_cts_notifications.name_filter"
	operationTypeFilterDataSourceName := "data.huaweicloud_cts_notifications.operation_type_filter"
	statusFilterDataSourceName := "data.huaweicloud_cts_notifications.status_filter"
	topicIDFilterDataSourceName := "data.huaweicloud_cts_notifications.topic_id_filter"

	dc := acceptance.InitDataSourceCheck(defaultDataSourceName)
	dcNameFilter := acceptance.InitDataSourceCheck(nameFilterDataSourceName)
	dcOperationTypeFilter := acceptance.InitDataSourceCheck(operationTypeFilterDataSourceName)
	dcStatusFilter := acceptance.InitDataSourceCheck(statusFilterDataSourceName)
	dcTopicIDFilter := acceptance.InitDataSourceCheck(topicIDFilterDataSourceName)

	name := acceptance.RandomAccResourceName()
	name2 := acceptance.RandomAccResourceName()
	ctsNotificationBaseConfig := testAccCTSNotification_basic(name)
	baseConfig := testAccDatasourceCTSNotifications_base(ctsNotificationBaseConfig, name2)
	disabledConfig := testAccDatasourceCTSNotifications_disabled_base(ctsNotificationBaseConfig, name2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceCTSNotifications_basic(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(defaultDataSourceName, "notifications.0.name"),
					resource.TestCheckResourceAttrSet(defaultDataSourceName, "notifications.0.operation_type"),
					resource.TestCheckResourceAttrSet(defaultDataSourceName, "notifications.0.status"),
					resource.TestCheckResourceAttrSet(defaultDataSourceName, "notifications.0.id"),
					resource.TestCheckResourceAttrSet(defaultDataSourceName, "notifications.0.topic_id"),
					resource.TestCheckOutput("default_is_useful", "true"),
				),
			},
			{
				Config: testAccDatasourceCTSNotifications_nameFilter(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dcNameFilter.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(nameFilterDataSourceName, "notifications.0.id"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
				),
			},
			{
				Config: testAccDatasourceCTSNotifications_operationTypeFilter(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dcOperationTypeFilter.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(operationTypeFilterDataSourceName, "notifications.0.id"),
					resource.TestCheckOutput("operation_type_filter_is_useful", "true"),
				),
			},
			{
				Config: testAccDatasourceCTSNotifications_statusFilter(disabledConfig),
				Check: resource.ComposeTestCheckFunc(
					dcStatusFilter.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(statusFilterDataSourceName, "notifications.0.id"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
			{
				Config: testAccDatasourceCTSNotifications_topicIDFilter(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dcTopicIDFilter.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(topicIDFilterDataSourceName, "notifications.0.id"),
					resource.TestCheckOutput("topic_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceCTSNotifications_base(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_smn_topic" "topic_2" {
  name  = "%[2]s"
}
  
resource "huaweicloud_cts_notification" "test" {
  name           = "%[2]s"
  operation_type = "customized"
  smn_topic      = huaweicloud_smn_topic.topic_2.id
  
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
`, baseConfig, name)
}

func testAccDatasourceCTSNotifications_disabled_base(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_smn_topic" "topic_2" {
  name  = "%[2]s"
}
  
resource "huaweicloud_cts_notification" "test" {
  name           = "%[2]s"
  operation_type = "customized"
  smn_topic      = huaweicloud_smn_topic.topic_2.id
  enabled        = false
  
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
`, baseConfig, name)
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

output "default_is_useful" {
  value = length(data.huaweicloud_cts_notifications.test.notifications) >= 2
}
`, config)
}

func testAccDatasourceCTSNotifications_nameFilter(config string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  name = huaweicloud_cts_notification.test.name
}

data "huaweicloud_cts_notifications" "name_filter" {
  type = "smn"
  name = local.name
}
  
output "name_filter_is_useful" {
  value = length(
    data.huaweicloud_cts_notifications.name_filter.notifications
  ) == 1 && data.huaweicloud_cts_notifications.name_filter.notifications[0].name == local.name
}
`, config)
}

func testAccDatasourceCTSNotifications_operationTypeFilter(config string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  operation_type = huaweicloud_cts_notification.notify.operation_type
}

data "huaweicloud_cts_notifications" "operation_type_filter" {
  type = "smn"
  operation_type = local.operation_type
}
  
output "operation_type_filter_is_useful" {
  value = length(data.huaweicloud_cts_notifications.operation_type_filter.notifications) > 0 && alltrue(
    [for v in data.huaweicloud_cts_notifications.operation_type_filter.notifications[*].operation_type : v == local.operation_type]
  )  
}
`, config)
}

func testAccDatasourceCTSNotifications_statusFilter(config string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  status = huaweicloud_cts_notification.test.status
}

data "huaweicloud_cts_notifications" "status_filter" {
  type = "smn"
  status = local.status
}
  
output "status_filter_is_useful" {
  value = length(data.huaweicloud_cts_notifications.status_filter.notifications) > 0 && alltrue(
    [for v in data.huaweicloud_cts_notifications.status_filter.notifications[*].status : v == local.status]
  )  
}
`, config)
}

func testAccDatasourceCTSNotifications_topicIDFilter(config string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  topic_id = huaweicloud_cts_notification.notify.smn_topic
}

data "huaweicloud_cts_notifications" "topic_id_filter" {
  type = "smn"
  topic_id = local.topic_id
}
  
output "topic_id_filter_is_useful" {
  value = length(data.huaweicloud_cts_notifications.topic_id_filter.notifications) > 0 && alltrue(
	[for v in data.huaweicloud_cts_notifications.topic_id_filter.notifications[*].topic_id : v == local.topic_id]
  )
}
`, config)
}
