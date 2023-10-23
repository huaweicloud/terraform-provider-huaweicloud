package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceLoadBalancers_basic(t *testing.T) {
	rName := "data.huaweicloud_elb_loadbalancers.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceLoadBalancers_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.#"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.name"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.ipv4_address"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.ipv4_port_id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.l4_flavor_id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.l7_flavor_id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.vpc_id"),
					resource.TestCheckOutput("filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceLoadBalancers_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_elb_loadbalancers" "test" {
  name = huaweicloud_elb_loadbalancer.test.name

  depends_on = [
    huaweicloud_elb_loadbalancer.test
  ]
}

data "huaweicloud_elb_loadbalancers" "name_filter" {
  name = huaweicloud_elb_loadbalancer.test.name
}

locals {
  name_filter_result = [for v in data.huaweicloud_elb_loadbalancers.name_filter.loadbalancers[*].name : v == data.huaweicloud_elb_loadbalancers.test.name]
}
 
output "filter_is_useful" {
  value = length(local.name_filter_result) > 0
}
`, testAccElbV3LoadBalancerConfig_basic(name))
}
