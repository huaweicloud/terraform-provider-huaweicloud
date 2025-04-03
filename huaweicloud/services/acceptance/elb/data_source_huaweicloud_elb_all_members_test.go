package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceElbAllMembers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_elb_all_members.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceElbAllMembers_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "members.#"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.member_type"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.address"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.subnet_cidr_id"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.protocol_port"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.weight"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.operating_status"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.ip_version"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.pool_id"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.loadbalancer_id"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.status.#"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.status.0.listener_id"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.status.0.operating_status"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.reason.#"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.updated_at"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("weight_filter_is_useful", "true"),
					resource.TestCheckOutput("subnet_cidr_id_filter_is_useful", "true"),
					resource.TestCheckOutput("address_filter_is_useful", "true"),
					resource.TestCheckOutput("protocol_port_filter_is_useful", "true"),
					resource.TestCheckOutput("member_id_filter_is_useful", "true"),
					resource.TestCheckOutput("operating_status_filter_is_useful", "true"),
					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
					resource.TestCheckOutput("ip_version_filter_is_useful", "true"),
					resource.TestCheckOutput("pool_id_filter_is_useful", "true"),
					resource.TestCheckOutput("loadbalancer_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceElbAllMembers_base(name string) string {
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
  name           = "%[2]s"
  vpc_id         = huaweicloud_vpc.test.id
  ipv4_subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  backend_subnets = [
    huaweicloud_vpc_subnet.test.id
 ]
}

resource "huaweicloud_elb_listener" "test" {
  name                        = "%[2]s"
  description                 = "test description"
  protocol                    = "HTTP"
  protocol_port               = 8083
  loadbalancer_id             = huaweicloud_elb_loadbalancer.test.id
  advanced_forwarding_enabled = false
}

resource "huaweicloud_elb_pool" "test" {
  name        = "%[2]s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = huaweicloud_elb_listener.test.id
}

resource "huaweicloud_elb_member" "test" {
  name          = "%[2]s"
  address       = huaweicloud_compute_instance.test.access_ip_v4
  weight        = 2
  protocol_port = 8080
  pool_id       = huaweicloud_elb_pool.test.id
  subnet_id     = huaweicloud_vpc_subnet.test.ipv4_subnet_id
}
`, common.TestBaseNetwork(name), name)
}

func testDataSourceElbAllMembers_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_elb_all_members" "test" {
  depends_on = [huaweicloud_elb_member.test]
}

locals {
  name = "%[2]s"
}

data "huaweicloud_elb_all_members" "name_filter" {
  depends_on = [huaweicloud_elb_member.test]

  name = ["%[2]s"]
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_elb_all_members.name_filter.members) > 0 && alltrue(
  [for v in data.huaweicloud_elb_all_members.name_filter.members[*].name : v == local.name]
  )
}

locals {
  weight = 2
}

data "huaweicloud_elb_all_members" "weight_filter" {
  depends_on = [huaweicloud_elb_member.test]

  weight = [2]
}

output "weight_filter_is_useful" {
  value = length(data.huaweicloud_elb_all_members.weight_filter.members) > 0 && alltrue(
  [for v in data.huaweicloud_elb_all_members.weight_filter.members[*].weight : v == local.weight]
  )
}

locals {
  subnet_cidr_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
}

data "huaweicloud_elb_all_members" "subnet_cidr_id_filter" {
  depends_on = [huaweicloud_elb_member.test]

  subnet_cidr_id = [huaweicloud_vpc_subnet.test.ipv4_subnet_id]
}

output "subnet_cidr_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_all_members.subnet_cidr_id_filter.members) > 0 && alltrue(
  [for v in data.huaweicloud_elb_all_members.subnet_cidr_id_filter.members[*].subnet_cidr_id : v == local.subnet_cidr_id]
  )
}

locals {
  address = huaweicloud_compute_instance.test.access_ip_v4
}

data "huaweicloud_elb_all_members" "address_filter" {
  depends_on = [huaweicloud_elb_member.test]

  address = [huaweicloud_compute_instance.test.access_ip_v4]
}

output "address_filter_is_useful" {
  value = length(data.huaweicloud_elb_all_members.address_filter.members) > 0 && alltrue(
  [for v in data.huaweicloud_elb_all_members.address_filter.members[*].address : v == local.address]
  )
}

locals {
  protocol_port = 8080
}

data "huaweicloud_elb_all_members" "protocol_port_filter" {
  depends_on = [huaweicloud_elb_member.test]

  protocol_port = [8080]
}

output "protocol_port_filter_is_useful" {
  value = length(data.huaweicloud_elb_all_members.protocol_port_filter.members) > 0 && alltrue(
  [for v in data.huaweicloud_elb_all_members.protocol_port_filter.members[*].protocol_port : v == local.protocol_port]
  )
}

locals {
  member_id = huaweicloud_elb_member.test.id
}

data "huaweicloud_elb_all_members" "member_id_filter" {
  depends_on = [huaweicloud_elb_member.test]

  member_id = [huaweicloud_elb_member.test.id]
}

output "member_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_all_members.member_id_filter.members) > 0 && alltrue(
  [for v in data.huaweicloud_elb_all_members.member_id_filter.members[*].id : v == local.member_id]
  )
}

locals {
  operating_status = data.huaweicloud_elb_all_members.test.members[0].operating_status
}

data "huaweicloud_elb_all_members" "operating_status_filter" {
  depends_on = [data.huaweicloud_elb_all_members.test]

  operating_status = [data.huaweicloud_elb_all_members.test.members[0].operating_status]
}

output "operating_status_filter_is_useful" {
  value = length(data.huaweicloud_elb_all_members.operating_status_filter.members) > 0 && alltrue(
  [for v in data.huaweicloud_elb_all_members.operating_status_filter.members[*].operating_status : v == local.operating_status]
  )
}

data "huaweicloud_elb_all_members" "enterprise_project_id_filter" {
  depends_on = [huaweicloud_elb_member.test]

  enterprise_project_id = [huaweicloud_elb_loadbalancer.test.enterprise_project_id]
}

output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_all_members.enterprise_project_id_filter.members) > 0
}

locals {
  ip_version = "v4"
}

data "huaweicloud_elb_all_members" "ip_version_filter" {
  depends_on = [huaweicloud_elb_member.test]

  ip_version = ["v4"]
}

output "ip_version_filter_is_useful" {
  value = length(data.huaweicloud_elb_all_members.ip_version_filter.members) > 0 && alltrue(
  [for v in data.huaweicloud_elb_all_members.ip_version_filter.members[*].ip_version : v == local.ip_version]
  )
}

locals {
  pool_id = huaweicloud_elb_pool.test.id
}

data "huaweicloud_elb_all_members" "pool_id_filter" {
  depends_on = [huaweicloud_elb_member.test]

  pool_id = [huaweicloud_elb_pool.test.id]
}

output "pool_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_all_members.pool_id_filter.members) > 0 && alltrue(
  [for v in data.huaweicloud_elb_all_members.pool_id_filter.members[*].pool_id : v == local.pool_id]
  )
}

locals {
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
}

data "huaweicloud_elb_all_members" "loadbalancer_id_filter" {
  depends_on = [huaweicloud_elb_member.test]

  loadbalancer_id = [huaweicloud_elb_loadbalancer.test.id]
}

output "loadbalancer_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_all_members.loadbalancer_id_filter.members) > 0 && alltrue(
  [for v in data.huaweicloud_elb_all_members.loadbalancer_id_filter.members[*].loadbalancer_id : v == local.loadbalancer_id]
  )
}
`, testDataSourceElbAllMembers_base(name), name)
}
