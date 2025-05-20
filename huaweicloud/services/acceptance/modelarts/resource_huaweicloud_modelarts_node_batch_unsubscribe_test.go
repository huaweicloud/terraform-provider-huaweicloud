package modelarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceV2NodeBatchUnsubscribe_basic(t *testing.T) {

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
				Config: testAccResourceV2NodeBatchUnsubscribe_basic_step1(),
			},
		},
	})
}

func testAccResourceV2NodeBatchUnsubscribe_basic_step1() string {
	return fmt.Sprintf(`
data "huaweicloud_modelartsv2_resource_pool_nodes" "test" {
  resource_pool_name = "%[1]s"
}

resource "huaweicloud_modelartsv2_node_batch_unsubscribe" "test" {
  resource_pool_name = "%[1]s"
  node_ids           = try(slice([for o in data.huaweicloud_modelartsv2_resource_pool_nodes.test.nodes:
    lookup(jsondecode(o.metadata[0].labels), "os.modelarts/resource.id",
    "") if try(o.metadata[0].labels, "") != ""], 0, 1), [])

  lifecycle {
    ignore_changes = [node_ids]
  }
}
`, acceptance.HW_MODELARTS_RESOURCE_POOL_NAME)
}
