package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceV2NodeBatchReset_basic(t *testing.T) {
	var (
		rName = "huaweicloud_modelartsv2_node_batch_reset.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelArtsResourcePoolIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceV2NodeBatchReset_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "pool_id"),
					resource.TestMatchResourceAttr(rName, "rolling_config.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttr(rName, "rolling_config.0.strategy", "RollingByPercent"),
					resource.TestCheckResourceAttr(rName, "rolling_config.0.max_unavailable", "20"),
				),
			},
			{
				Config: testAccResourceV2NodeBatchReset_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "pool_id"),
					resource.TestMatchResourceAttr(rName, "rolling_config.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttr(rName, "rolling_config.0.strategy", "RollingByNumber"),
					resource.TestCheckResourceAttr(rName, "rolling_config.0.max_unavailable", "1"),
				),
			},
		},
	})
}

func testAccResourceV2NodeBatchReset_base() string {
	return fmt.Sprintf(`
locals {
  resource_pood_ids = split(",", "%[1]s")
}

data "huaweicloud_modelartsv2_resource_pool_nodes" "test" {
  resource_pool_name = local.resource_pood_ids[0]
}

locals {
  node_names = try(slice([for o in data.huaweicloud_modelartsv2_resource_pool_nodes.test.nodes:
    o.metadata[0].name], 0, 1), [])
}
`, acceptance.HW_MODELARTS_RESOURCE_POOL_IDS)
}

func testAccResourceV2NodeBatchReset_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelartsv2_node_batch_reset" "test" {
  pool_id    = local.resource_pood_ids[0]
  node_names = local.node_names

  rolling_config {
    strategy        = "RollingByPercent"
    max_unavailable = 20
  }
}
`, testAccResourceV2NodeBatchReset_base())
}

func testAccResourceV2NodeBatchReset_basic_update() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelartsv2_node_batch_reset" "test" {
  pool_id    = local.resource_pood_ids[0]
  node_names = local.node_names

  rolling_config {
    strategy        = "RollingByNumber"
    max_unavailable = 1
  }

  enable_force_new = "true"
}
`, testAccResourceV2NodeBatchReset_base())
}
