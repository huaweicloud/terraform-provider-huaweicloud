package lb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataLBPools_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_lb_pools.test"

	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataLBPools_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "pools.0.name", name),
					resource.TestCheckResourceAttrPair(rName, "pools.0.id",
						"huaweicloud_lb_pool.pool_1", "id"),
					resource.TestCheckResourceAttrPair(rName, "pools.0.description",
						"huaweicloud_lb_pool.pool_1", "description"),
					resource.TestCheckResourceAttrPair(rName, "pools.0.protocol",
						"huaweicloud_lb_pool.pool_1", "protocol"),
					resource.TestCheckResourceAttrPair(rName, "pools.0.lb_method",
						"huaweicloud_lb_pool.pool_1", "lb_method"),
				),
			},
		},
	})
}

func testAccDataLBPools_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_lb_pools" "test" {
  name = "%s"

  depends_on = [
    huaweicloud_lb_pool.pool_1
  ]
}
`, testAccLBV2PoolConfig_basic(name), name)
}
