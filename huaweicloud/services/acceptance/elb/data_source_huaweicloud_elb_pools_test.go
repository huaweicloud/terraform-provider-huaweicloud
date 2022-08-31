package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourcePools_basic(t *testing.T) {
	rName := "data.huaweicloud_elb_pools.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePools_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "pools.0.name", name),
					resource.TestCheckResourceAttrPair(rName, "pools.0.id",
						"huaweicloud_elb_pool.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "pools.0.description",
						"huaweicloud_elb_pool.test", "description"),
					resource.TestCheckResourceAttrPair(rName, "pools.0.protocol",
						"huaweicloud_elb_pool.test", "protocol"),
					resource.TestCheckResourceAttrPair(rName, "pools.0.lb_method",
						"huaweicloud_elb_pool.test", "lb_method"),
				),
			},
		},
	})
}

func testAccDatasourcePools_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_elb_pools" "test" {
  name = "%s"

  depends_on = [
    huaweicloud_elb_pool.test
  ]
}
`, testAccElbV3PoolConfig_basic(name), name)
}
