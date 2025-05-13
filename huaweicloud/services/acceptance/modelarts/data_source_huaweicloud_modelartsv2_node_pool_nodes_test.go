package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this acceptance test, please support at least two resource nodes under the resource pool.
func TestAccDataSourceV2NodePoolNodes_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_modelartsv2_node_pool_nodes.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelArtsResourcePoolName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceV2NodePoolNodes_invalidResourcePoolName,
				ExpectError: regexp.MustCompile(`\\"invalid-resource-pool-name\\" not found`),
			},
			{
				Config:      testAccDataSourceV2NodePoolNodes_invalidNodePoolName(),
				ExpectError: regexp.MustCompile(`invalid nodepool name`),
			},
			{
				Config: testAccDataSourceV2NodePoolNodes_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "nodes.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_metadata_set_and_valid", "true"),
					resource.TestCheckOutput("is_spec_set_and_valid", "true"),
					resource.TestCheckOutput("is_status_set_and_valid", "true"),
				),
			},
		},
	})
}

const testAccDataSourceV2NodePoolNodes_invalidResourcePoolName string = `
data "huaweicloud_modelartsv2_node_pool_nodes" "test" {
  resource_pool_name = "invalid-resource-pool-name"
  node_pool_name     = "invalid-node-pool-name"
}
`

func testAccDataSourceV2NodePoolNodes_invalidNodePoolName() string {
	return fmt.Sprintf(`
data "huaweicloud_modelartsv2_node_pool_nodes" "test" {
  resource_pool_name = "%[1]s"
  node_pool_name     = "invalid-node-pool-name"
}
`, acceptance.HW_MODELARTS_RESOURCE_POOL_NAME)
}

func testAccDataSourceV2NodePoolNodes_basic_step1() string {
	return fmt.Sprintf(`
data "huaweicloud_modelartsv2_resource_pools" "test" {}

locals {
  node_pool_flavor_id = try([for o in data.huaweicloud_modelartsv2_resource_pools.test.resource_pools:
    o.resources[0].flavor_id if o.metadata[0].name == "%[1]s"][0], "")
}

data "huaweicloud_modelartsv2_node_pool_nodes" "test" {
  resource_pool_name = "%[1]s"
  node_pool_name     = format("%%s-default", local.node_pool_flavor_id)
}

locals {
  metadata = data.huaweicloud_modelartsv2_node_pool_nodes.test.nodes[0].metadata
  spec     = data.huaweicloud_modelartsv2_node_pool_nodes.test.nodes[0].spec
  status   = data.huaweicloud_modelartsv2_node_pool_nodes.test.nodes[0].status
}

output "is_metadata_set_and_valid" {
  value = try(length(local.metadata) > 0 && alltrue([
    local.metadata[0].annotations != "",
    local.metadata[0].creation_timestamp != "",
    local.metadata[0].labels != "",
    local.metadata[0].name != "",
  ]))
}

output "is_spec_set_and_valid" {
  value = try(length(local.spec) > 0 && alltrue([
    local.spec[0].flavor != "",
    length(local.spec[0].host_network) > 0 && alltrue([
      local.spec[0].host_network.0.vpc != "",
      local.spec[0].host_network.0.subnet != "",
    ]),
    length(local.spec[0].os) > 0 && alltrue([
      local.spec[0].os[0].image_id != "",
    ]),
  ]))
}

output "is_status_set_and_valid" {
  value = try(length(local.status) > 0 && alltrue([
    local.status[0].phase != "",
    local.status[0].az != "",
    length(local.status[0].os) > 0 && alltrue([
      local.status[0].os.0.name != "",
    ]),
  ]))
}
`, acceptance.HW_MODELARTS_RESOURCE_POOL_NAME)
}
