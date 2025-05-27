package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceV2NodeBatchDelete_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceV2NodeBatchDelete_invalidResourcePoolName,
				ExpectError: regexp.MustCompile(`\\"invalid-resource-pool-name\\" not found`),
			},
			{
				Config:      testAccResourceV2NodeBatchDelete_invalidResourceNodeNames(),
				ExpectError: regexp.MustCompile(`\\"invalid-node-name\\" not found`),
			},
			{
				Config: testAccResourceV2NodeBatchDelete_basic_step1(),
			},
		},
	})
}

const testAccResourceV2NodeBatchDelete_invalidResourcePoolName string = `
resource "huaweicloud_modelartsv2_node_batch_delete" "test" {
  resource_pool_name = "invalid-resource-pool-name"
  node_names         = ["invalid-node-name"]
}
`

func testAccResourceV2NodeBatchDelete_invalidResourceNodeNames() string {
	return fmt.Sprintf(`
resource "huaweicloud_modelartsv2_node_batch_delete" "test" {
  resource_pool_name = "%[1]s"
  node_names         = ["invalid-node-name"]
}
`, acceptance.HW_MODELARTS_RESOURCE_POOL_NAME)
}

func testAccResourceV2NodeBatchDelete_basic_step1() string {
	return fmt.Sprintf(`
data "huaweicloud_modelartsv2_resource_pool_nodes" "test" {
  resource_pool_name = "%[1]s"
}

resource "huaweicloud_modelartsv2_node_batch_delete" "test" {
  resource_pool_name = "%[1]s"
  node_names         = try(slice([for o in data.huaweicloud_modelartsv2_resource_pool_nodes.test.nodes:
    o.metadata[0].name if lookup(jsondecode(o.metadata[0].annotations), "os.modelarts/billing.mode",
    "0") == "0" && try(o.metadata[0].labels, "") != ""], 0, 1), [])

  lifecycle {
    # After the deletion is completed, the query result of the data source will change, so the reference to the names
    # need to be ignored.
    ignore_changes = [node_names]
  }
}
`, acceptance.HW_MODELARTS_RESOURCE_POOL_NAME)
}
