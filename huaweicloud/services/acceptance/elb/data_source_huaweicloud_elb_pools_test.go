package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDatasourcePools_basic(t *testing.T) {
	dataSource := "data.huaweicloud_elb_pools.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePools_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "pools.#"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.protocol"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.ip_version"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.lb_method"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.vpc_id"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.any_port_enable"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.member_deletion_protection_enable"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.public_border_group"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.protection_status"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.slow_start_enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.slow_start_duration"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.connection_drain_enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.connection_drain_timeout"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.minimum_healthy_member_count"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.healthmonitor_id"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.listeners.#"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.listeners.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.loadbalancers.#"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.loadbalancers.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.members.#"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.members.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.persistence.#"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.persistence.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.persistence.0.cookie_name"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.persistence.0.timeout"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.updated_at"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("pool_id_filter_is_useful", "true"),
					resource.TestCheckOutput("description_filter_is_useful", "true"),
					resource.TestCheckOutput("any_port_enable_filter_is_useful", "true"),
					resource.TestCheckOutput("connection_drain_filter_is_useful", "true"),
					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
					resource.TestCheckOutput("ip_version_filter_is_useful", "true"),
					resource.TestCheckOutput("member_address_filter_is_useful", "true"),
					resource.TestCheckOutput("member_device_id_filter_is_useful", "true"),
					resource.TestCheckOutput("member_instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("member_deletion_protection_enable_filter_is_useful", "true"),
					resource.TestCheckOutput("pool_health_filter_is_useful", "true"),
					resource.TestCheckOutput("public_border_group_filter_is_useful", "true"),
					resource.TestCheckOutput("healthmonitor_id_filter_is_useful", "true"),
					resource.TestCheckOutput("loadbalancer_id_filter_is_useful", "true"),
					resource.TestCheckOutput("listener_id_filter_is_useful", "true"),
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

func TestAccDatasourcePools_quic(t *testing.T) {
	dataSource := "data.huaweicloud_elb_pools.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePools_quic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.quic_cid_hash_strategy.#"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.quic_cid_hash_strategy.0.len"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.quic_cid_hash_strategy.0.offset"),

					resource.TestCheckOutput("quic_cid_len_filter_is_useful", "true"),
					resource.TestCheckOutput("quic_cid_offset_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourcePools_base(name string) string {
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
  name                         = "%[2]s"
  protocol                     = "HTTP"
  lb_method                    = "ROUND_ROBIN"
  type                         = "instance"
  vpc_id                       = huaweicloud_vpc.test.id
  description                  = "test pool description"
  loadbalancer_id              = huaweicloud_elb_loadbalancer.test.id
  listener_id                  = huaweicloud_elb_listener.test.id
  minimum_healthy_member_count = 1

  persistence {
    type        = "APP_COOKIE"
    cookie_name = "testCookie"
    timeout     = 10
  }
}

resource "huaweicloud_elb_member" "test" {
  address       = huaweicloud_compute_instance.test.access_ip_v4
  protocol_port = 8000
  name          = "%[2]s"
  weight        = 20
  pool_id       = huaweicloud_elb_pool.test.id
  subnet_id     = huaweicloud_vpc_subnet.test.ipv4_subnet_id
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
  depends_on = [
    huaweicloud_elb_monitor.test,
    huaweicloud_elb_member.test
  ]

  pool_id = huaweicloud_elb_pool.test.id
}

data "huaweicloud_elb_pools" "name_filter" {
  depends_on = [
    huaweicloud_elb_monitor.test,
    huaweicloud_elb_member.test
  ]

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
  depends_on = [
    huaweicloud_elb_monitor.test,
    huaweicloud_elb_member.test
  ]

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
  depends_on = [
    huaweicloud_elb_monitor.test,
    huaweicloud_elb_member.test
  ]

  description = huaweicloud_elb_pool.test.description
}

output "description_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.description_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_pools.description_filter.pools[*].description : v == local.description]
  )
}

locals {
  any_port_enable = huaweicloud_elb_pool.test.any_port_enable
}

data "huaweicloud_elb_pools" "any_port_enable_filter" {
  depends_on = [
    huaweicloud_elb_monitor.test,
    huaweicloud_elb_member.test
  ]

  any_port_enable = huaweicloud_elb_pool.test.any_port_enable
}

output "any_port_enable_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.any_port_enable_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_pools.any_port_enable_filter.pools[*].any_port_enable : v == local.any_port_enable]
  )
}

locals {
  connection_drain = huaweicloud_elb_pool.test.connection_drain_enabled
}

data "huaweicloud_elb_pools" "connection_drain_filter" {
  depends_on = [
    huaweicloud_elb_monitor.test,
    huaweicloud_elb_member.test
  ]

  connection_drain = huaweicloud_elb_pool.test.connection_drain_enabled
}

output "connection_drain_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.connection_drain_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_pools.connection_drain_filter.pools[*].connection_drain_enabled : v == local.connection_drain]
  )
}

locals {
  enterprise_project_id = huaweicloud_elb_pool.test.enterprise_project_id
}

data "huaweicloud_elb_pools" "enterprise_project_id_filter" {
  depends_on = [
    huaweicloud_elb_monitor.test,
    huaweicloud_elb_member.test
  ]

  enterprise_project_id = huaweicloud_elb_pool.test.enterprise_project_id
}

output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.enterprise_project_id_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_pools.enterprise_project_id_filter.pools[*].enterprise_project_id :
  v == local.enterprise_project_id]
  )
}

locals {
  ip_version = huaweicloud_elb_pool.test.ip_version
}

data "huaweicloud_elb_pools" "ip_version_filter" {
  depends_on = [
    huaweicloud_elb_monitor.test,
    huaweicloud_elb_member.test
  ]

  ip_version = huaweicloud_elb_pool.test.ip_version
}

output "ip_version_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.ip_version_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_pools.ip_version_filter.pools[*].ip_version : v == local.ip_version]
  )
}

locals {
  member_address = huaweicloud_elb_member.test.address
}

data "huaweicloud_elb_pools" "member_address_filter" {
  depends_on = [
    huaweicloud_elb_monitor.test,
    huaweicloud_elb_member.test
  ]

  member_address = huaweicloud_elb_member.test.address
}

output "member_address_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.member_address_filter.pools) > 0
}

locals {
  member_device_id = huaweicloud_elb_member.test.instance_id
}

data "huaweicloud_elb_pools" "member_device_id_filter" {
  depends_on = [
    huaweicloud_elb_monitor.test,
    huaweicloud_elb_member.test
  ]

  member_device_id = huaweicloud_elb_member.test.instance_id
}

output "member_device_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.member_device_id_filter.pools) > 0
}

locals {
  member_instance_id = huaweicloud_elb_member.test.instance_id
}

data "huaweicloud_elb_pools" "member_instance_id_filter" {
  depends_on = [
    huaweicloud_elb_monitor.test,
    huaweicloud_elb_member.test
  ]

  member_instance_id = huaweicloud_elb_member.test.instance_id
}

output "member_instance_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.member_instance_id_filter.pools) > 0
}

locals {
  member_deletion_protection_enable = huaweicloud_elb_pool.test.deletion_protection_enable
}

data "huaweicloud_elb_pools" "member_deletion_protection_enable_filter" {
  depends_on = [
    huaweicloud_elb_monitor.test,
    huaweicloud_elb_member.test
  ]

  member_deletion_protection_enable = huaweicloud_elb_pool.test.deletion_protection_enable
}

output "member_deletion_protection_enable_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.member_deletion_protection_enable_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_pools.member_deletion_protection_enable_filter.pools[*].member_deletion_protection_enable :
  v == local.member_deletion_protection_enable]
  )
}

locals {
  minimum_healthy_member_count = huaweicloud_elb_pool.test.minimum_healthy_member_count
}

data "huaweicloud_elb_pools" "pool_health_filter" {
  depends_on = [
    huaweicloud_elb_monitor.test,
    huaweicloud_elb_member.test
  ]

  pool_health = "minimum_healthy_member_count=1"
}

output "pool_health_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.pool_health_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_pools.pool_health_filter.pools[*].minimum_healthy_member_count :
  v == local.minimum_healthy_member_count]
  )
}

locals {
  public_border_group = huaweicloud_elb_pool.test.public_border_group
}

data "huaweicloud_elb_pools" "public_border_group_filter" {
  depends_on = [
    huaweicloud_elb_monitor.test,
    huaweicloud_elb_member.test
  ]

  public_border_group = huaweicloud_elb_pool.test.public_border_group
}

output "public_border_group_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.public_border_group_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_pools.public_border_group_filter.pools[*].public_border_group : v == local.public_border_group]
  )
}

locals {
  healthmonitor_id = huaweicloud_elb_monitor.test.id
}

data "huaweicloud_elb_pools" "healthmonitor_id_filter" {
  depends_on = [
    huaweicloud_elb_monitor.test,
    huaweicloud_elb_member.test
  ]

  healthmonitor_id = huaweicloud_elb_monitor.test.id
}

output "healthmonitor_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.healthmonitor_id_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_pools.healthmonitor_id_filter.pools[*].healthmonitor_id : v == local.healthmonitor_id]
  )
}

locals {
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
}

data "huaweicloud_elb_pools" "loadbalancer_id_filter" {
  depends_on = [
    huaweicloud_elb_monitor.test,
    huaweicloud_elb_member.test
  ]

  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
}

output "loadbalancer_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.loadbalancer_id_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_pools.loadbalancer_id_filter.pools[*].loadbalancers[0].id : v == local.loadbalancer_id]
  )
}

locals {
  listener_id = huaweicloud_elb_listener.test.id
}

data "huaweicloud_elb_pools" "listener_id_filter" {
  depends_on = [
    huaweicloud_elb_monitor.test,
    huaweicloud_elb_member.test
  ]

  listener_id = huaweicloud_elb_listener.test.id
}

output "listener_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.listener_id_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_pools.listener_id_filter.pools[*].listeners[0].id : v == local.listener_id]
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

func testAccDatasourcePools_quic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_pool" "test" {
  name      = "%[2]s"
  protocol  = "QUIC"
  lb_method = "QUIC_CID"
  type      = "instance"
  vpc_id    = huaweicloud_vpc.test.id
}

data "huaweicloud_elb_pools" "test" {
  depends_on = [huaweicloud_elb_pool.test]
}

locals {
  quic_cid_len = data.huaweicloud_elb_pools.test.pools[0].quic_cid_hash_strategy[0].len
}

data "huaweicloud_elb_pools" "quic_cid_len_filter" {
  depends_on = [
    huaweicloud_elb_pool.test,
    data.huaweicloud_elb_pools.test
  ]

  quic_cid_len = data.huaweicloud_elb_pools.test.pools[0].quic_cid_hash_strategy[0].len
}

output "quic_cid_len_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.quic_cid_len_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_pools.quic_cid_len_filter.pools[*].quic_cid_hash_strategy.0.len : v == local.quic_cid_len]
  )
}

locals {
  quic_cid_offset = data.huaweicloud_elb_pools.test.pools[0].quic_cid_hash_strategy[0].offset
}

data "huaweicloud_elb_pools" "quic_cid_offset_filter" {
  depends_on = [
    huaweicloud_elb_pool.test,
    data.huaweicloud_elb_pools.test
  ]

  quic_cid_offset = data.huaweicloud_elb_pools.test.pools[0].quic_cid_hash_strategy[0].offset
}

output "quic_cid_offset_filter_is_useful" {
  value = length(data.huaweicloud_elb_pools.quic_cid_offset_filter.pools) > 0 && alltrue(
  [for v in data.huaweicloud_elb_pools.quic_cid_offset_filter.pools[*].quic_cid_hash_strategy.0.offset : v == local.quic_cid_offset]
  )
}
`, common.TestVpc(name), name)
}
