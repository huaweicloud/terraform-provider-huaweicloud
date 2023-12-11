package dms

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

func getDmsKafkaSmartConnectTasksResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getKafkaSmartConnectTasks: query DMS kafka smart connect tasks
	var (
		getKafkaSmartConnectTasksHttpUrl = "v2/{project_id}/connectors/{connector_id}/sink-tasks"
		getKafkaSmartConnectTasksProduct = "dms"
	)
	getKafkaSmartConnectTasksClient, err := cfg.NewServiceClient(getKafkaSmartConnectTasksProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS Client: %s", err)
	}

	connnectorID := state.Primary.Attributes["connector_id"]
	getKafkaSmartConnectTasksPath := getKafkaSmartConnectTasksClient.Endpoint + getKafkaSmartConnectTasksHttpUrl
	getKafkaSmartConnectTasksPath = strings.ReplaceAll(getKafkaSmartConnectTasksPath, "{project_id}",
		getKafkaSmartConnectTasksClient.ProjectID)
	getKafkaSmartConnectTasksPath = strings.ReplaceAll(getKafkaSmartConnectTasksPath, "{connector_id}", connnectorID)

	getKafkaSmartConnectTasksOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getKafkaSmartConnectTasksResp, err := getKafkaSmartConnectTasksClient.Request("GET", getKafkaSmartConnectTasksPath,
		&getKafkaSmartConnectTasksOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DMS smart connect tasks: %s", err)
	}

	getKafkaSmartConnectTasksBody, err := utils.FlattenResponse(getKafkaSmartConnectTasksResp)
	if err != nil {
		return nil, err
	}
	return getKafkaSmartConnectTasksBody, nil
}

func TestAccDmsKafkaSmartConnectTasks_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_dms_kafka_smart_connect_tasks.test"

	rc := acceptance.InitResourceCheck(
		dataSourceName,
		&obj,
		getDmsKafkaSmartConnectTasksResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDmsKafkaSmartConnectTasks_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "tasks.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tasks.0.task_id"),
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
  task_id = huaweicloud_dms_kafka_smart_connect_task.test.id
}

output "task_id_filter_is_useful" {
  value = length(data.huaweicloud_dms_kafka_smart_connect_tasks.task_id_filter.tasks) > 0 && alltrue(
  [for v in data.huaweicloud_dms_kafka_smart_connect_tasks.task_id_filter.tasks[*].task_id : v == local.task_id]
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
