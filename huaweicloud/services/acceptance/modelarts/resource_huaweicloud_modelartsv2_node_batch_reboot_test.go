package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceV2NodeBatchReboot_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelArtsResourcePoolName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceV2NodeBatchReboot_invalidResourcePoolName,
				ExpectError: regexp.MustCompile(`\\"invalid-resource-pool-name\\" not found`),
			},
			{
				Config: testAccResourceV2NodeBatchReboot_basic_step1(),
			},
		},
	})
}

const testAccResourceV2NodeBatchReboot_invalidResourcePoolName string = `
resource "huaweicloud_modelartsv2_node_batch_reboot" "invalid" {
  resource_pool_name = "invalid-resource-pool-name"
  node_names         = ["invalid-node-name"]
}
`

func testAccResourceV2NodeBatchReboot_basic_step1() string {
	return fmt.Sprintf(`
data "huaweicloud_modelartsv2_resource_pools" "test" {}

locals {
  resourcePool = [for pool in data.huaweicloud_modelartsv2_resource_pools.test.resource_pools : pool if pool.name == "%[1]s"][0]
}

data "huaweicloud_modelartsv2_resource_pool_nodes" "test" {
  resource_pool_name = local.resourcePool.metadata[0].name
}

locals {
  node_names = try(data.huaweicloud_modelartsv2_resource_pool_nodes.test.nodes[*].metadata[0].name, [])
}

resource "huaweicloud_modelartsv2_node_batch_reboot" "test" {
  resource_pool_name = local.resourcePool.metadata[0].name
  node_names         = local.node_names

  lifecycle {
    ignore_changes = [
      node_names,
    ]
  }
}
`, acceptance.HW_MODELARTS_RESOURCE_POOL_NAME)
}
