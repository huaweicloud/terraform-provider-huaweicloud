package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
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
					resource.TestCheckResourceAttr(rName, "pools.0.type", "instance"),
					resource.TestCheckResourceAttrPair(rName, "pools.0.vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(rName, "pools.0.protection_status", "nonProtection"),
					resource.TestCheckResourceAttr(rName, "pools.0.slow_start_enabled", "false"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("pool_id_filter_is_useful", "true"),
					resource.TestCheckOutput("description_filter_is_useful", "true"),
					resource.TestCheckOutput("healthmonitor_id_filter_is_useful", "true"),
					resource.TestCheckOutput("protocol_filter_is_useful", "true"),
					resource.TestCheckOutput("lb_method_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("vpc_id_filter_is_useful", "true"),
					resource.TestCheckOutput("protection_status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourcePools_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name            = "%[2]s"
  ipv4_subnet_id  = huaweicloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
}

resource "huaweicloud_elb_listener" "test" {
  name            = "%[2]s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
}

resource "huaweicloud_elb_pool" "test" {
  name            = "%[2]s"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  type            = "instance"
  vpc_id          = huaweicloud_vpc.test.id
  description     = "test pool description"
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
  listener_id     = huaweicloud_elb_listener.test.id

  minimum_healthy_member_count = 1
}

resource "huaweicloud_elb_monitor" "test" {
  pool_id          = huaweicloud_elb_pool.test.id
  name             = "%[2]s"
  protocol         = "HTTP"
  interval         = 20
  timeout          = 10
  max_retries      = 5
  max_retries_down = 5
  url_path         = "/aa"
  domain_name      = "www.aa.com"
  port             = "8000"
  status_code      = "200,401-500,502"
  enabled          = false
}
`, common.TestBaseNetwork(name), name)
}

func testAccDatasourcePools_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_elb_pools" "test" {
  depends_on = [huaweicloud_elb_monitor.test]

  pool_id = huaweicloud_elb_pool.test.id
}

data "huaweicloud_elb_pools" "name_filter" {
  depends_on = [huaweicloud_elb_monitor.test]

  name = "%[2]s"
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.name_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_pools.name_filter.pools[*].name :v == "%[2]s"]
  )
}

locals {
  pool_id = huaweicloud_elb_pool.test.id
}

data "huaweicloud_elb_pools" "pool_id_filter" {
  depends_on = [huaweicloud_elb_monitor.test]

  pool_id = huaweicloud_elb_pool.test.id
}

output "pool_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.pool_id_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_pools.pool_id_filter.pools[*].id : v == local.pool_id]
  )
}

locals {
  description = huaweicloud_elb_pool.test.description
}

data "huaweicloud_elb_pools" "description_filter" {
  depends_on = [huaweicloud_elb_monitor.test]

  description = huaweicloud_elb_pool.test.description
}

output "description_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.description_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_pools.description_filter.pools[*].description : v == local.description]
  )
}

locals {
  healthmonitor_id = huaweicloud_elb_monitor.test.id
}

data "huaweicloud_elb_pools" "healthmonitor_id_filter" {
  depends_on = [huaweicloud_elb_monitor.test]

  healthmonitor_id = huaweicloud_elb_monitor.test.id
}

output "healthmonitor_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.healthmonitor_id_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_pools.healthmonitor_id_filter.pools[*].healthmonitor_id : v == local.healthmonitor_id]
  )
}

locals {
  protocol = huaweicloud_elb_pool.test.protocol
}

data "huaweicloud_elb_pools" "protocol_filter" {
  depends_on = [huaweicloud_elb_monitor.test]

  protocol = huaweicloud_elb_pool.test.protocol
}

output "protocol_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.description_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_pools.description_filter.pools[*].protocol : v == local.protocol]
  )
}

locals {
  lb_method = huaweicloud_elb_pool.test.lb_method
}

data "huaweicloud_elb_pools" "lb_method_filter" {
  depends_on = [huaweicloud_elb_monitor.test]

  lb_method = huaweicloud_elb_pool.test.lb_method
}

output "lb_method_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.lb_method_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_pools.lb_method_filter.pools[*].lb_method : v == local.lb_method]
  )
}

locals {
  type = huaweicloud_elb_pool.test.type
}

data "huaweicloud_elb_pools" "type_filter" {
  depends_on = [huaweicloud_elb_monitor.test]

  type = huaweicloud_elb_pool.test.type
}

output "type_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.type_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_pools.type_filter.pools[*].type : v == local.type]
  )
}

locals {
  vpc_id = huaweicloud_elb_pool.test.vpc_id
}

data "huaweicloud_elb_pools" "vpc_id_filter" {
  depends_on = [huaweicloud_elb_monitor.test]

  vpc_id = huaweicloud_elb_pool.test.vpc_id
}

output "vpc_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.vpc_id_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_pools.vpc_id_filter.pools[*].vpc_id : v == local.vpc_id]
  )
}

locals {
  protection_status = huaweicloud_elb_pool.test.protection_status
}

data "huaweicloud_elb_pools" "protection_status_filter" {
  depends_on = [huaweicloud_elb_monitor.test]

  protection_status = huaweicloud_elb_pool.test.protection_status
}

output "protection_status_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.protection_status_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_pools.protection_status_filter.pools[*].protection_status : v == local.protection_status]
  )
}
`, testAccDatasourcePools_base(name), name)
}
