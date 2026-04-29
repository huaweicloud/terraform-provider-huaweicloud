package dws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccWorkloadQueueUpdateAction_basic(t *testing.T) {
	var (
		queueName    = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_dws_workload_queue_update_action.test"
	)

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWorkloadQueueUpdateAction_basic(queueName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "cluster_id", acceptance.HW_DWS_CLUSTER_ID),
					resource.TestCheckResourceAttr(resourceName, "name", queueName),
					resource.TestCheckResourceAttr(resourceName, "configuration.#", "5"),
				),
			},
		},
	})
}

func testAccWorkloadQueueUpdateAction_basic(name string) string {
	return fmt.Sprintf(`
variable "queue_configuration_list" {
  type        = list(object({
    resource_name  = string
    resource_value = number
  }))
  description = "The configuration list to create the workload queue"

  default = [
    {
      resource_name  = "cpu_limit"
      resource_value = 10
    },
    {
      resource_name  = "cpu_share"
      resource_value = 10
    },
    {
      resource_name  = "memory"
      resource_value = 10
    },
    {
      resource_name  = "tablespace"
      resource_value = -1
    },
    {
      resource_name  = "activestatements"
      resource_value = -1
    },
  ]
}

variable "update_configuration_list" {
  type        = list(object({
    resource_name  = string
    resource_value = number
  }))
  description = "The configuration list to update the workload queue"

  default = [
    {
      resource_name  = "cpu_limit"
      resource_value = 12
    },
    {
      resource_name  = "cpu_share"
      resource_value = 12
    },
    {
      resource_name  = "memory"
      resource_value = 15
    },
    {
      resource_name  = "tablespace"
      resource_value = -1
    },
    {
      resource_name  = "activestatements"
      resource_value = -1
    },
  ]
}

resource "huaweicloud_dws_workload_queue" "test" {
  cluster_id           = "%[1]s"
  name                 = "%[2]s"
  logical_cluster_name = "%[3]s"

  dynamic "configuration" {
    for_each = var.queue_configuration_list

    content {
      resource_name  = configuration.value.resource_name
      resource_value = configuration.value.resource_value
    }
  }
}

resource "huaweicloud_dws_workload_queue_update_action" "test" {
  cluster_id           = "%[1]s"
  name                 = "%[2]s"
  logical_cluster_name = "%[3]s"

  dynamic "configuration" {
    for_each = var.update_configuration_list

    content {
      resource_name  = configuration.value.resource_name
      resource_value = configuration.value.resource_value
    }
  }

  depends_on = [huaweicloud_dws_workload_queue.test]
}
  `, acceptance.HW_DWS_CLUSTER_ID, name, acceptance.HW_DWS_LOGICAL_CLUSTER_NAME)
}
