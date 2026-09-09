package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this acceptance test, please creating a resource pool.
func TestAccDataSourceV2NodePools_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_modelartsv2_node_pools.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelArtsResourcePoolName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceV2NodePools_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "node_pools.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "node_pools.0.metadata.#"),
					resource.TestCheckResourceAttrSet(dataSource, "node_pools.0.metadata.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "node_pools.0.spec.0.resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "node_pools.0.spec.0.resources.0.node_pool"),
					resource.TestCheckResourceAttrSet(dataSource, "node_pools.0.spec.0.resources.0.flavor"),
					resource.TestCheckResourceAttrSet(dataSource, "node_pools.0.spec.0.resources.0.count"),
					resource.TestCheckResourceAttrSet(dataSource, "node_pools.0.spec.0.resources.0.extend_params.#"),
					resource.TestCheckResourceAttrSet(dataSource, "node_pools.0.spec.0.resources.0.extend_params.0.docker_base_size"),
					resource.TestCheckResourceAttrSet(dataSource, "node_pools.0.spec.0.resources.0.extend_params.0.runtime"),
					resource.TestCheckResourceAttrSet(dataSource, "node_pools.0.spec.0.resources.0.extend_params.0.node_pool_name"),
					resource.TestCheckResourceAttrSet(dataSource, "node_pools.0.spec.0.resources.0.os.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "node_pools.0.status.0.resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "node_pools.0.status.0.resources.0.available.#"),
					resource.TestCheckResourceAttrSet(dataSource, "node_pools.0.status.0.resources.0.available.0.flavor"),
					resource.TestCheckResourceAttrSet(dataSource, "node_pools.0.status.0.resources.0.available.0.max_count"),
					resource.TestCheckResourceAttrSet(dataSource, "node_pools.0.status.0.resources.0.available.0.azs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "node_pools.0.status.0.resources.0.available.0.azs.0.az"),
					resource.TestCheckResourceAttrSet(dataSource, "node_pools.0.status.0.resources.0.available.0.azs.0.count"),
				),
			},
		},
	})
}

func testAccDataSourceV2NodePools_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_modelartsv2_node_pools" "test" {
  pool_name = "%s"
}
`, acceptance.HW_MODELARTS_RESOURCE_POOL_NAME)
}
