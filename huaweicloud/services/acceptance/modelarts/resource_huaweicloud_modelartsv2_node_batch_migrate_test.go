package modelarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceV2NodeBatchMigrate_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelArtsResourcePoolIds(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceV2NodeBatchMigrate_basic(),
			},
		},
	})
}

func testAccResourceV2NodeBatchMigrate_base() string {
	return fmt.Sprintf(`
locals {
  resource_pood_ids = split(",", "%[1]s")
}

data "huaweicloud_modelartsv2_resource_pool_nodes" "test" {
  resource_pool_name = local.resource_pood_ids[0]
}

# The source resource pool must contain 2 or more nodes
locals {
  node_names = try(slice([for o in data.huaweicloud_modelartsv2_resource_pool_nodes.test.nodes:
    o.metadata[0].name], 0, 1), [])
  flavor = try(slice([for o in data.huaweicloud_modelartsv2_resource_pool_nodes.test.nodes:
    o.spec[0].flavor], 0, 1), [])
}
`, acceptance.HW_MODELARTS_RESOURCE_POOL_IDS)
}

func testAccResourceV2NodeBatchMigrate_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelartsv2_node_batch_migrate" "test" {
  source_pool_id    = local.resource_pood_ids[0]
  source_cluster_id = local.resource_pood_ids[0]
  target_pool_id    = local.resource_pood_ids[1]
  target_cluster_id = local.resource_pood_ids[1]
  node_names        = local.node_names

  resource_spec {
    flavor = local.flavor[0]
  }

  lifecycle {
    ignore_changes = [
      node_names,
    ]
  }
}`, testAccResourceV2NodeBatchMigrate_base())
}
