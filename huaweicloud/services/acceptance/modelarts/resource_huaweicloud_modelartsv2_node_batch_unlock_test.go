package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceV2NodeBatchUnlock_basic(t *testing.T) {
	var (
		rName = "huaweicloud_modelartsv2_node_batch_unlock.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceV2NodeBatchUnlock_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(rName, "node_names.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
				),
			},
		},
	})
}

func testAccResourceV2NodeBatchUnlock_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_modelartsv2_resource_pool_nodes" "test" {
  resource_pool_name = "%[1]s"
}

locals {
  node_names = try(slice([for o in data.huaweicloud_modelartsv2_resource_pool_nodes.test.nodes:
    o.metadata[0].name], 0, 1), [])
}

resource "huaweicloud_modelartsv2_node_batch_lock" "test" {
  pool_id    = "%[1]s"
  node_names = local.node_names
}

resource "huaweicloud_modelartsv2_node_batch_unlock" "test" {
  pool_id    = "%[1]s"
  node_names = local.node_names

  depends_on = [
    huaweicloud_modelartsv2_node_batch_lock.test
  ]
}`, acceptance.HW_MODELARTS_RESOURCE_POOL_ID)
}
