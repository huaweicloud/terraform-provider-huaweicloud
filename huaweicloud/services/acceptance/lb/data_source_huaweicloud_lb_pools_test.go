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

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("pool_id_filter_is_useful", "true"),
					resource.TestCheckOutput("description_filter_is_useful", "true"),
					resource.TestCheckOutput("loadbalancer_id_filter_is_useful", "true"),
					resource.TestCheckOutput("healthmonitor_id_filter_is_useful", "true"),
					resource.TestCheckOutput("protocol_filter_is_useful", "true"),
					resource.TestCheckOutput("lb_method_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataLBPools_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_lb_pools" "test" {
  depends_on = [huaweicloud_lb_monitor.monitor_1]

  pool_id = huaweicloud_lb_pool.pool_1.id
}

data "huaweicloud_lb_pools" "name_filter" {
  depends_on = [huaweicloud_lb_monitor.monitor_1]

  name = "%[2]s"
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_lb_pools.name_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_lb_pools.name_filter.pools[*].name :v == "%[2]s"]
  )
}

locals {
  pool_id = huaweicloud_lb_pool.pool_1.id
}

data "huaweicloud_lb_pools" "pool_id_filter" {
  depends_on = [huaweicloud_lb_monitor.monitor_1]

  pool_id = huaweicloud_lb_pool.pool_1.id
}

output "pool_id_filter_is_useful" {
  value = length(data.huaweicloud_lb_pools.pool_id_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_lb_pools.pool_id_filter.pools[*].id : v == local.pool_id]
  )
}

locals {
  description = huaweicloud_lb_pool.pool_1.description
}

data "huaweicloud_lb_pools" "description_filter" {
  depends_on = [huaweicloud_lb_monitor.monitor_1]

  description = huaweicloud_lb_pool.pool_1.description
}

output "description_filter_is_useful" {
  value = length(data.huaweicloud_lb_pools.description_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_lb_pools.description_filter.pools[*].description : v == local.description]
  )
}

locals {
  loadbalancer_id = huaweicloud_lb_pool.pool_1.loadbalancer_id
}

data "huaweicloud_lb_pools" "loadbalancer_id_filter" {
  depends_on = [huaweicloud_lb_monitor.monitor_1]

  loadbalancer_id = huaweicloud_lb_pool.pool_1.loadbalancer_id
}

output "loadbalancer_id_filter_is_useful" {
  value = length(data.huaweicloud_lb_pools.loadbalancer_id_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_lb_pools.loadbalancer_id_filter.pools[*].loadbalancers.0.id : v == local.loadbalancer_id]
  )
}

locals {
  healthmonitor_id = huaweicloud_lb_monitor.monitor_1.id
}

data "huaweicloud_lb_pools" "healthmonitor_id_filter" {
  depends_on = [huaweicloud_lb_monitor.monitor_1]

  healthmonitor_id = huaweicloud_lb_monitor.monitor_1.id
}

output "healthmonitor_id_filter_is_useful" {
  value = length(data.huaweicloud_lb_pools.healthmonitor_id_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_lb_pools.healthmonitor_id_filter.pools[*].healthmonitor_id : v == local.healthmonitor_id]
  )
}

locals {
  protocol = huaweicloud_lb_pool.pool_1.protocol
}

data "huaweicloud_lb_pools" "protocol_filter" {
  depends_on = [huaweicloud_lb_monitor.monitor_1]

  protocol = huaweicloud_lb_pool.pool_1.protocol
}

output "protocol_filter_is_useful" {
  value = length(data.huaweicloud_lb_pools.protocol_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_lb_pools.protocol_filter.pools[*].protocol : v == local.protocol]
  )
}

locals {
  lb_method = huaweicloud_lb_pool.pool_1.lb_method
}

data "huaweicloud_lb_pools" "lb_method_filter" {
  depends_on = [huaweicloud_lb_monitor.monitor_1]

  lb_method = huaweicloud_lb_pool.pool_1.lb_method
}

output "lb_method_filter_is_useful" {
  value = length(data.huaweicloud_lb_pools.lb_method_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_lb_pools.lb_method_filter.pools[*].lb_method : v == local.lb_method]
  )
}
`, testAccLBV2MonitorConfig_basic(name), name)
}
