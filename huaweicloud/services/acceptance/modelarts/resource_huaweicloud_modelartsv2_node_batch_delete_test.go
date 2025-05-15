package modelarts

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceV2NodeBatchDelete_basic(t *testing.T) {
	var (
		allResourcePools      = "data.huaweicloud_modelartsv2_resource_pools.test"
		dcForAllResourcePools = acceptance.InitDataSourceCheck(allResourcePools)

		allNodes      = "data.huaweicloud_modelartsv2_node_pool_nodes.test"
		dcForAllNodes = acceptance.InitDataSourceCheck(allNodes)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceV2NodeBatchDelete_basic_step1,
				Check: resource.ComposeTestCheckFunc(
					dcForAllResourcePools.CheckResourceExists(),
					resource.TestMatchResourceAttr(allResourcePools, "resource_pools.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
			{
				Config: testAccResourceV2NodeBatchDelete_basic_step2,
				Check: resource.ComposeTestCheckFunc(
					dcForAllNodes.CheckResourceExists(),
					resource.TestMatchResourceAttr(allNodes, "nodes.#", regexp.MustCompile(`^([2-9]\d*|\d{2,})$`)),
				),
			},
			{
				Config: testAccResourceV2NodeBatchDelete_basic_step3,
			},
		},
	})
}

const testAccResourceV2NodeBatchDelete_basic_step1 string = `
data "huaweicloud_modelartsv2_resource_pools" "test" {}
`

const testAccResourceV2NodeBatchDelete_basic_step2 string = `
data "huaweicloud_modelartsv2_resource_pools" "test" {}

data "huaweicloud_modelartsv2_node_pool_nodes" "test" {
  resource_pool_name = data.huaweicloud_modelartsv2_resource_pools.test.resource_pools[0].metadata[0].name
  node_pool_name     = format("%s-default", data.huaweicloud_modelartsv2_resource_pools.test.resource_pools[0].resources[0].flavor_id)
}
`

const testAccResourceV2NodeBatchDelete_basic_step3 string = `
data "huaweicloud_modelartsv2_resource_pools" "test" {}

data "huaweicloud_modelartsv2_node_pool_nodes" "test" {
  resource_pool_name = data.huaweicloud_modelartsv2_resource_pools.test.resource_pools[0].metadata[0].name
  node_pool_name     = format("%s-default", data.huaweicloud_modelartsv2_resource_pools.test.resource_pools[0].resources[0].flavor_id)
}

resource "huaweicloud_modelartsv2_node_batch_delete" "test" {
  resource_pool_name = data.huaweicloud_modelartsv2_resource_pools.test.resource_pools[0].metadata[0].name
  node_names         = try(slice(data.huaweicloud_modelartsv2_node_pool_nodes.test.nodes[*].metadata[0].name, 0, 1), [])

  lifecycle {
    ignore_changes = [node_names]
  }
}
`
