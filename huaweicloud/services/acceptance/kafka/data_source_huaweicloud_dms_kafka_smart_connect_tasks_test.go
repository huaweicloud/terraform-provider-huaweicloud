package kafka

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDmsKafkaSmartConnectTasks_basic(t *testing.T) {
	name := acceptance.RandomAccResourceNameWithDash()
	dataSourceName := "data.huaweicloud_dms_kafka_smart_connect_tasks.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDmsKafkaSmartConnectTasks_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "tasks.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tasks.0.task_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tasks.0.destination_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tasks.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tasks.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tasks.0.topics"),
					resource.TestCheckOutput("task_id_filter_is_useful", "true"),
					resource.TestCheckOutput("task_name_filter_is_useful", "true"),
					resource.TestCheckOutput("destination_type_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDmsKafkaSmartConnectTasks_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_kafka_smart_connect_tasks" "test" {
  depends_on   = [huaweicloud_dms_kafka_smart_connect_task.test]
  connector_id = huaweicloud_dms_kafka_smart_connect.test.id
}

data "huaweicloud_dms_kafka_smart_connect_tasks" "task_id_filter" {
  depends_on   = [huaweicloud_dms_kafka_smart_connect_task.test]
  connector_id = huaweicloud_dms_kafka_smart_connect.test.id
  task_id      = huaweicloud_dms_kafka_smart_connect_task.test.id
}

locals {
  id = huaweicloud_dms_kafka_smart_connect_task.test.id
}

output "task_id_filter_is_useful" {
  value = length(data.huaweicloud_dms_kafka_smart_connect_tasks.task_id_filter.tasks) > 0 && alltrue(
  [for v in data.huaweicloud_dms_kafka_smart_connect_tasks.task_id_filter.tasks[*].id : v == local.id]
  )  
}

data "huaweicloud_dms_kafka_smart_connect_tasks" "task_name_filter" {
  depends_on   = [huaweicloud_dms_kafka_smart_connect_task.test]
  connector_id = huaweicloud_dms_kafka_smart_connect.test.id
  task_name    = huaweicloud_dms_kafka_smart_connect_task.test.task_name
}
  
locals {
  task_name = huaweicloud_dms_kafka_smart_connect_task.test.task_name
}
  
output "task_name_filter_is_useful" {
  value = length(data.huaweicloud_dms_kafka_smart_connect_tasks.task_name_filter.tasks) > 0 && alltrue(
  [for v in data.huaweicloud_dms_kafka_smart_connect_tasks.task_name_filter.tasks[*].task_name : v == local.task_name]
  )
}

data "huaweicloud_dms_kafka_smart_connect_tasks" "destination_type_filter" {
  depends_on       = [huaweicloud_dms_kafka_smart_connect_task.test]
  connector_id     = huaweicloud_dms_kafka_smart_connect.test.id
  destination_type = huaweicloud_dms_kafka_smart_connect_task.test.destination_type
}
	
locals {
  destination_type = huaweicloud_dms_kafka_smart_connect_task.test.destination_type
}
	
output "destination_type_filter_is_useful" {
  value = length(data.huaweicloud_dms_kafka_smart_connect_tasks.destination_type_filter.tasks) > 0 && alltrue(
  [for v in data.huaweicloud_dms_kafka_smart_connect_tasks.destination_type_filter.tasks[*].destination_type : v == local.destination_type]
  ) 
}

data "huaweicloud_dms_kafka_smart_connect_tasks" "status_filter" {
  depends_on   = [huaweicloud_dms_kafka_smart_connect_task.test]
  connector_id = huaweicloud_dms_kafka_smart_connect.test.id
  status       = huaweicloud_dms_kafka_smart_connect_task.test.status
}
	  
locals {
  status = huaweicloud_dms_kafka_smart_connect_task.test.status
}
	  
output "status_filter_is_useful" {
  value = length(data.huaweicloud_dms_kafka_smart_connect_tasks.status_filter.tasks) > 0 && alltrue(
  [for v in data.huaweicloud_dms_kafka_smart_connect_tasks.status_filter.tasks[*].status : v == local.status]
  ) 
}

`, testDmsKafkaSmartConnectTask_basic(name))
}
