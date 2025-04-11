package lb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
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
					resource.TestCheckResourceAttrSet(rName, "pools.#"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.name"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.id"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.description"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.protocol"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.lb_method"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.healthmonitor_id"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.protection_status"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.listeners.#"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.listeners.0.id"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.loadbalancers.#"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.loadbalancers.0.id"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.members.#"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.members.0.id"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.persistence.#"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.persistence.0.type"),
					resource.TestCheckResourceAttrSet(rName, "pools.0.persistence.0.cookie_name"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("pool_id_filter_is_useful", "true"),
					resource.TestCheckOutput("description_filter_is_useful", "true"),
					resource.TestCheckOutput("loadbalancer_id_filter_is_useful", "true"),
					resource.TestCheckOutput("member_address_filter_is_useful", "true"),
					resource.TestCheckOutput("member_device_id_filter_is_useful", "true"),
					resource.TestCheckOutput("healthmonitor_id_filter_is_useful", "true"),
					resource.TestCheckOutput("protocol_filter_is_useful", "true"),
					resource.TestCheckOutput("lb_method_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataLBPools_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 22.04 server 64bit"
  most_recent = true
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_lb_loadbalancer" "test" {
  name          = "%[2]s"
  vip_subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
}

resource "huaweicloud_lb_listener" "test" {
  name            = "%[2]s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_lb_loadbalancer.test.id
}

resource "huaweicloud_lb_pool" "test" {
  name        = "%[2]s"
  description = "test description"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = huaweicloud_lb_listener.test.id

  persistence {
    type        = "APP_COOKIE"
    cookie_name = "testCookie"
  }
}

resource "huaweicloud_lb_member" "test" {
  address       = huaweicloud_compute_instance.test.access_ip_v4
  protocol_port = 8080
  pool_id       = huaweicloud_lb_pool.test.id
  subnet_id     = huaweicloud_vpc_subnet.test.ipv4_subnet_id
}

resource "huaweicloud_lb_monitor" "test" {
  pool_id     = huaweicloud_lb_pool.test.id
  name        = "%[2]s"
  type        = "TCP"
  delay       = 20
  timeout     = 10
  max_retries = 5
  domain_name = "testdomain.com"
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccDataLBPools_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_lb_pools" "test" {
  depends_on = [
    huaweicloud_lb_member.test,
    huaweicloud_lb_monitor.test
  ]
}

data "huaweicloud_lb_pools" "name_filter" {
  depends_on = [
    huaweicloud_lb_member.test,
    huaweicloud_lb_monitor.test
  ]

  name = "%[2]s"
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_lb_pools.name_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_lb_pools.name_filter.pools[*].name :v == "%[2]s"]
  )
}

locals {
  pool_id = huaweicloud_lb_pool.test.id
}

data "huaweicloud_lb_pools" "pool_id_filter" {
  depends_on = [
    huaweicloud_lb_member.test,
    huaweicloud_lb_monitor.test
  ]

  pool_id = huaweicloud_lb_pool.test.id
}

output "pool_id_filter_is_useful" {
  value = length(data.huaweicloud_lb_pools.pool_id_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_lb_pools.pool_id_filter.pools[*].id : v == local.pool_id]
  )
}

locals {
  description = huaweicloud_lb_pool.test.description
}

data "huaweicloud_lb_pools" "description_filter" {
  depends_on = [
    huaweicloud_lb_member.test,
    huaweicloud_lb_monitor.test
  ]

  description = huaweicloud_lb_pool.test.description
}

output "description_filter_is_useful" {
  value = length(data.huaweicloud_lb_pools.description_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_lb_pools.description_filter.pools[*].description : v == local.description]
  )
}

locals {
  loadbalancer_id = huaweicloud_lb_pool.test.loadbalancer_id
}

data "huaweicloud_lb_pools" "loadbalancer_id_filter" {
  depends_on = [
    huaweicloud_lb_member.test,
    huaweicloud_lb_monitor.test
  ]

  loadbalancer_id = huaweicloud_lb_pool.test.loadbalancer_id
}

output "loadbalancer_id_filter_is_useful" {
  value = length(data.huaweicloud_lb_pools.loadbalancer_id_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_lb_pools.loadbalancer_id_filter.pools[*].loadbalancers.0.id : v == local.loadbalancer_id]
  )
}

locals {
  member_address = huaweicloud_compute_instance.test.access_ip_v4
}

data "huaweicloud_lb_pools" "member_address_filter" {
  depends_on = [
    huaweicloud_lb_member.test,
    huaweicloud_lb_monitor.test
  ]

  member_address = huaweicloud_compute_instance.test.access_ip_v4
}

output "member_address_filter_is_useful" {
  value = length(data.huaweicloud_lb_pools.member_address_filter.pools) > 0
}

locals {
  member_device_id = huaweicloud_compute_instance.test.id
}

data "huaweicloud_lb_pools" "member_device_id_filter" {
  depends_on = [
    huaweicloud_lb_member.test,
    huaweicloud_lb_monitor.test
  ]

  member_device_id = huaweicloud_compute_instance.test.id
}

output "member_device_id_filter_is_useful" {
  value = length(data.huaweicloud_lb_pools.member_device_id_filter.pools) > 0
}

locals {
  healthmonitor_id = huaweicloud_lb_monitor.test.id
}

data "huaweicloud_lb_pools" "healthmonitor_id_filter" {
  depends_on = [
    huaweicloud_lb_member.test,
    huaweicloud_lb_monitor.test
  ]

  healthmonitor_id = huaweicloud_lb_monitor.test.id
}

output "healthmonitor_id_filter_is_useful" {
  value = length(data.huaweicloud_lb_pools.healthmonitor_id_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_lb_pools.healthmonitor_id_filter.pools[*].healthmonitor_id : v == local.healthmonitor_id]
  )
}

locals {
  protocol = huaweicloud_lb_pool.test.protocol
}

data "huaweicloud_lb_pools" "protocol_filter" {
  depends_on = [
    huaweicloud_lb_member.test,
    huaweicloud_lb_monitor.test
  ]

  protocol = huaweicloud_lb_pool.test.protocol
}

output "protocol_filter_is_useful" {
  value = length(data.huaweicloud_lb_pools.protocol_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_lb_pools.protocol_filter.pools[*].protocol : v == local.protocol]
  )
}

locals {
  lb_method = huaweicloud_lb_pool.test.lb_method
}

data "huaweicloud_lb_pools" "lb_method_filter" {
  depends_on = [
    huaweicloud_lb_member.test,
    huaweicloud_lb_monitor.test
  ]

  lb_method = huaweicloud_lb_pool.test.lb_method
}

output "lb_method_filter_is_useful" {
  value = length(data.huaweicloud_lb_pools.lb_method_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_lb_pools.lb_method_filter.pools[*].lb_method : v == local.lb_method]
  )
}
`, testAccDataLBPools_base(name), name)
}
