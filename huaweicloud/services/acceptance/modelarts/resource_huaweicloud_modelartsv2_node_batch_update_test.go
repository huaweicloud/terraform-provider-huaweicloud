package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceV2NodeBatchUpdate_basic(t *testing.T) {
	var (
		rName = "huaweicloud_modelartsv2_node_batch_update.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceV2NodeBatchUpdate_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(rName, "node_names.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttr(rName, "action", "closeHaRedundant"),
				),
			},
			{
				Config: testAccResourceV2NodeBatchUpdate_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(rName, "node_names.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttr(rName, "action", "createTags"),
					resource.TestCheckResourceAttr(rName, "tags.#", "2"),
				),
			},
		},
	})
}

func testAccResourceV2NodeBatchUpdate_base() string {
	return fmt.Sprintf(`
data "huaweicloud_modelartsv2_resource_pool_nodes" "test" {
  resource_pool_name = "%[1]s"
}

locals {
  node_names = try(slice([for o in data.huaweicloud_modelartsv2_resource_pool_nodes.test.nodes:
    o.metadata[0].name], 0, 1), [])
}
`, acceptance.HW_MODELARTS_RESOURCE_POOL_ID)
}

func testAccResourceV2NodeBatchUpdate_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelartsv2_node_batch_update" "test" {
  pool_id             = "%[2]s"
  node_names          = local.node_names
  action              = "closeHaRedundant"
  ha_redundant_effect = "NoExecute"
}
`, testAccResourceV2NodeBatchUpdate_base(), acceptance.HW_MODELARTS_RESOURCE_POOL_ID)
}

func testAccResourceV2NodeBatchUpdate_basic_update() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelartsv2_node_batch_update" "test" {
  pool_id    = "%[2]s"
  node_names = local.node_names
  action     = "createTags"

  tags {
    key   = "foo1"
    value = "bar1"
  }

  tags {
    key   = "foo2"
    value = "bar2"
  }

  enable_force_new = "true"
}
`, testAccResourceV2NodeBatchUpdate_base(), acceptance.HW_MODELARTS_RESOURCE_POOL_ID)
}
