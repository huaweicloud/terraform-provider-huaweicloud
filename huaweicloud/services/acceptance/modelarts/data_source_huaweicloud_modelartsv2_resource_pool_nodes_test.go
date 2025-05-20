package modelarts

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this acceptance test, please support at least one resource pool.
func TestAccDataV2ResourcePoolNodes_basic(t *testing.T) {
	allResourcePools := "data.huaweicloud_modelartsv2_resource_pools.test"
	dcForAllResourcePools := acceptance.InitDataSourceCheck(allResourcePools)

	allResourcePoolNodes := "data.huaweicloud_modelartsv2_resource_pool_nodes.test"
	dcForAllResourcePoolNodes := acceptance.InitDataSourceCheck(allResourcePoolNodes)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataV2ResourcePoolNodes_invalidResourcePoolName,
				ExpectError: regexp.MustCompile(`\\"invalid-resource-pool-name\\" not found`),
			},
			{
				Config: testAccDataV2ResourcePoolNodes_basic_step1,
				Check: resource.ComposeTestCheckFunc(
					dcForAllResourcePools.CheckResourceExists(),
					resource.TestMatchResourceAttr(allResourcePools, "resource_pools.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
			{
				Config: testAccDataV2ResourcePoolNodes_basic_step2,
				Check: resource.ComposeTestCheckFunc(
					dcForAllResourcePoolNodes.CheckResourceExists(),
					resource.TestMatchResourceAttr(allResourcePoolNodes, "nodes.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckOutput("is_metadata_set_and_valid", "true"),
					resource.TestCheckOutput("is_spec_set_and_valid", "true"),
					resource.TestCheckOutput("is_status_set_and_valid", "true"),
				),
			},
		},
	})
}

const testAccDataV2ResourcePoolNodes_invalidResourcePoolName string = `
data "huaweicloud_modelartsv2_resource_pool_nodes" "invalid_resource_pool_name" {
  resource_pool_name = "invalid-resource-pool-name"
}
`

const testAccDataV2ResourcePoolNodes_basic_step1 string = `
data "huaweicloud_modelartsv2_resource_pools" "test" {}
`

const testAccDataV2ResourcePoolNodes_basic_step2 string = `
data "huaweicloud_modelartsv2_resource_pools" "test" {}

data "huaweicloud_modelartsv2_resource_pool_nodes" "test" {
  resource_pool_name = try(data.huaweicloud_modelartsv2_resource_pools.test.resource_pools[0].metadata[0].name, "")
}

locals {
  metadata = data.huaweicloud_modelartsv2_resource_pool_nodes.test.nodes[0].metadata
  spec     = data.huaweicloud_modelartsv2_resource_pool_nodes.test.nodes[0].spec
  status   = data.huaweicloud_modelartsv2_resource_pool_nodes.test.nodes[0].status
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
    local.spec[0].extend_params != "",
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
`
