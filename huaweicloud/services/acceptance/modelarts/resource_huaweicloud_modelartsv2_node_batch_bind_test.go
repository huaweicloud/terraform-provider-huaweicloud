package modelarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceV2NodeBatchBind_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelArtsResourcePoolIds(t, 1)
			acceptance.TestAccPreCheckModelArtsResourcePoolQuotaName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceV2NodeBatchBind_basic(),
			},
		},
	})
}

func testAccResourceV2NodeBatchBind_base() string {
	return fmt.Sprintf(`
data "huaweicloud_modelartsv2_resource_pool_nodes" "test" {
  resource_pool_name = "%[1]s"
}

locals {
  node_names = [for o in data.huaweicloud_modelartsv2_resource_pool_nodes.test.nodes: o.metadata[0].name]
}
`, acceptance.HW_MODELARTS_RESOURCE_POOL_IDS)
}

func testAccResourceV2NodeBatchBind_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelartsv2_node_batch_bind" "test" {
  pool_id = "%[2]s"

  nodes {
    name       = local.node_names[0]
    quota_name = "%[3]s"
  }

  drain = true
}
`, testAccResourceV2NodeBatchBind_base(), acceptance.HW_MODELARTS_RESOURCE_POOL_IDS,
		acceptance.HW_MODELARTS_RESOURCE_POOL_QUOTA_NAME)
}
