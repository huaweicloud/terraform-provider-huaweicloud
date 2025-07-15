package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this acceptance test, make sure that the node you are operating is an A3 (hyperinstance) prepaid node.
func TestAccResourceResourcePoolNodeBatchResize_basic(t *testing.T) {
	name := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelArtsResourcePoolName(t)
			acceptance.TestAccPreCheckModelArtsResourcePoolBatchResize(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceResourcePoolNodeBatchResize_invalidResourcePoolName(name),
				ExpectError: regexp.MustCompile(`\\"invalid-resource-pool-name\\" not found`),
			},
			{
				Config: testAccResourceResourcePoolNodeBatchResize_basic(name),
			},
		},
	})
}

func testAccResourceResourcePoolNodeBatchResize_base() string {
	return fmt.Sprintf(`
data "huaweicloud_modelartsv2_resource_pool_nodes" "test" {
  resource_pool_name = "%[1]s"
}

data "huaweicloud_modelartsv2_resource_pools" "test" {
  status = "created"
}

locals {
  action_nodes = try([for v in data.huaweicloud_modelartsv2_resource_pool_nodes.test.nodes : jsondecode(v.metadata[0].labels)
  if contains(split(",","%[2]s"), lookup(jsondecode(v.metadata[0].labels), "os.modelarts.node/batch.name", ""))], [])
  source_node_pool_name = element(local.action_nodes[*]["os.modelarts.node/nodepool"], 0)
  source_node_pool_configs = try([for item in data.huaweicloud_modelartsv2_resource_pools.test.resource_pools :
    [for v in item.resources : v if "%%{if v.node_pool == ""}${v.flavor_id}-default%%{else}${v.node_pool}%%{endif}" == local.source_node_pool_name]
    if item.metadata[0].name == "%[1]s"
  ][0], [])
}
`, acceptance.HW_MODELARTS_RESOURCE_POOL_NAME, acceptance.HW_MODELARTS_RESOURCE_POOL_BATCH_RESIZE_NODE_NAME)
}

func testAccResourceResourcePoolNodeBatchResize_invalidResourcePoolName(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_resource_pool_node_batch_resize" "test" {
  resource_pool_name = "invalid-resource-pool-name"

  dynamic "nodes" {
    for_each = distinct(local.action_nodes[*]["os.modelarts.node/batch.uid"])

    content {
      batch_uid = nodes.value
    }
  }

  source {
    node_pool = local.source_node_pool_name
    flavor    = try(local.source_node_pool_configs[0].flavor_id, "")

    dynamic "creating_step" {
      for_each = try(local.source_node_pool_configs[0].creating_step, [])

      content {
        type = creating_step.value.type
        step = creating_step.value.step
      }
    }
  }

  target {
    node_pool = "%[2]s"
    flavor    = try(local.source_node_pool_configs[0].flavor_id, "")

    creating_step {
      type = "hyperinstance"
      step = 2
    }
  }

  billing = jsonencode({
    autoPay = "1"
  })

  lifecycle {
    # After the deletion is completed, the query result of the data source will change, so the reference to the source
    # need to be ignored.
    ignore_changes = [source]
  }
}
`, testAccResourceResourcePoolNodeBatchResize_base(), name)
}

func testAccResourceResourcePoolNodeBatchResize_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_resource_pool_node_batch_resize" "test" {
  resource_pool_name = "%[2]s"

  dynamic "nodes" {
    for_each = distinct(local.action_nodes[*]["os.modelarts.node/batch.uid"])

    content {
      batch_uid = nodes.value
    }
  }

  source {
    node_pool = local.source_node_pool_name
    flavor    = try(local.source_node_pool_configs[0].flavor_id, "")

    dynamic "creating_step" {
      for_each = try(local.source_node_pool_configs[0].creating_step, [])

      content {
        type = creating_step.value.type
        step = creating_step.value.step
      }
    }
  }

  target {
    node_pool = "%[3]s"
    flavor    = try(local.source_node_pool_configs[0].flavor_id, "")

    creating_step {
      type = "hyperinstance"
      step = 2
    }
  }

  billing = jsonencode({
    autoPay = "1"
  })

  lifecycle {
    # After the deletion is completed, the query result of the data source will change, so the reference to the source
    # need to be ignored.
    ignore_changes = [source]
  }
}
`, testAccResourceResourcePoolNodeBatchResize_base(), acceptance.HW_MODELARTS_RESOURCE_POOL_NAME, name)
}
