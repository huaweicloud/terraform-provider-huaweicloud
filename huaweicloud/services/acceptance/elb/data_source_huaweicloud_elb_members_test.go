package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDatasourceMembers_basic(t *testing.T) {
	rName := "data.huaweicloud_elb_members.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceMembers_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "members.#"),
					resource.TestCheckResourceAttrSet(rName, "members.0.name"),
					resource.TestCheckResourceAttrSet(rName, "members.0.id"),
					resource.TestCheckResourceAttrSet(rName, "members.0.address"),
					resource.TestCheckResourceAttrSet(rName, "members.0.protocol_port"),
					resource.TestCheckResourceAttrSet(rName, "members.0.subnet_id"),
					resource.TestCheckResourceAttrSet(rName, "members.0.weight"),
					resource.TestCheckResourceAttrSet(rName, "members.0.member_type"),
					resource.TestCheckResourceAttrSet(rName, "members.0.instance_id"),
					resource.TestCheckResourceAttrSet(rName, "members.0.ip_version"),
					resource.TestCheckResourceAttrSet(rName, "members.0.operating_status"),
					resource.TestCheckResourceAttrSet(rName, "members.0.reason.#"),
					resource.TestCheckResourceAttrSet(rName, "members.0.status.#"),
					resource.TestCheckResourceAttrSet(rName, "members.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "members.0.updated_at"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("member_id_filter_is_useful", "true"),
					resource.TestCheckOutput("address_filter_is_useful", "true"),
					resource.TestCheckOutput("protocol_port_filter_is_useful", "true"),
					resource.TestCheckOutput("weight_filter_is_useful", "true"),
					resource.TestCheckOutput("subnet_id_filter_is_useful", "true"),
					resource.TestCheckOutput("instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("ip_version_filter_is_useful", "true"),
					resource.TestCheckOutput("operating_status_filter_is_useful", "true"),
					resource.TestCheckOutput("member_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceMembers_base(name string) string {
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
  protocol_port = 8080
  pool_id       = huaweicloud_elb_pool.test.id
  subnet_id     = huaweicloud_vpc_subnet.test.ipv4_subnet_id
}
`, common.TestBaseNetwork(name), name)
}

func testAccDatasourceMembers_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_elb_members" "test" {
  pool_id    = huaweicloud_elb_pool.test.id
  depends_on = [huaweicloud_elb_member.test]
}

data "huaweicloud_elb_members" "name_filter" {
  pool_id    = huaweicloud_elb_pool.test.id
  name       = "%[2]s"
  depends_on = [huaweicloud_elb_member.test]
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_elb_members.name_filter.members) > 0 && alltrue(
  [for v in data.huaweicloud_elb_members.name_filter.members[*].name :v == "%[2]s"]
  )  
}

data "huaweicloud_elb_members" "member_id_filter" {
  pool_id    = huaweicloud_elb_pool.test.id
  member_id  = huaweicloud_elb_member.test.id
  depends_on = [huaweicloud_elb_member.test]
}

locals {
  member_id = huaweicloud_elb_member.test.id
}

output "member_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_members.member_id_filter.members) > 0 && alltrue(
  [for v in data.huaweicloud_elb_members.member_id_filter.members[*].id : v == local.member_id]
  )  
}

data "huaweicloud_elb_members" "address_filter" {
  pool_id    = huaweicloud_elb_pool.test.id
  address    = huaweicloud_elb_member.test.address
  depends_on = [huaweicloud_elb_member.test]
}

locals {
  address = huaweicloud_elb_member.test.address
}

output "address_filter_is_useful" {
  value = length(data.huaweicloud_elb_members.address_filter.members) > 0 && alltrue(
  [for v in data.huaweicloud_elb_members.address_filter.members[*].address : v == local.address]
  )  
}

data "huaweicloud_elb_members" "protocol_port_filter" {
  pool_id       = huaweicloud_elb_pool.test.id
  protocol_port = huaweicloud_elb_member.test.protocol_port
  depends_on    = [huaweicloud_elb_member.test]
}

locals {
  protocol_port = huaweicloud_elb_member.test.protocol_port
}

output "protocol_port_filter_is_useful" {
  value = length(data.huaweicloud_elb_members.protocol_port_filter.members) > 0 && alltrue(
  [for v in data.huaweicloud_elb_members.protocol_port_filter.members[*].protocol_port : v == local.protocol_port]
  )  
}

data "huaweicloud_elb_members" "weight_filter" {
  pool_id    = huaweicloud_elb_pool.test.id
  weight     = huaweicloud_elb_member.test.weight
  depends_on = [huaweicloud_elb_member.test]
}

locals {
  weight = huaweicloud_elb_member.test.weight
}

output "weight_filter_is_useful" {
  value = length(data.huaweicloud_elb_members.weight_filter.members) > 0 && alltrue(
  [for v in data.huaweicloud_elb_members.weight_filter.members[*].weight : v == local.weight]
  )  
}

data "huaweicloud_elb_members" "subnet_id_filter" {
  pool_id    = huaweicloud_elb_pool.test.id
  subnet_id  = huaweicloud_elb_member.test.subnet_id
  depends_on = [huaweicloud_elb_member.test]
}

locals {
  subnet_id = huaweicloud_elb_member.test.subnet_id
}

output "subnet_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_members.subnet_id_filter.members) > 0 && alltrue(
  [for v in data.huaweicloud_elb_members.subnet_id_filter.members[*].subnet_id : v == local.subnet_id]
  )  
}

data "huaweicloud_elb_members" "instance_id_filter" {
  pool_id     = huaweicloud_elb_pool.test.id
  instance_id = huaweicloud_elb_member.test.instance_id
  depends_on  = [huaweicloud_elb_member.test]
}

locals {
  instance_id = huaweicloud_elb_member.test.instance_id
}

output "instance_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_members.instance_id_filter.members) > 0 && alltrue(
  [for v in data.huaweicloud_elb_members.instance_id_filter.members[*].instance_id : v == local.instance_id]
  )  
}

data "huaweicloud_elb_members" "ip_version_filter" {
  pool_id    = huaweicloud_elb_pool.test.id
  ip_version = huaweicloud_elb_member.test.ip_version
  depends_on = [huaweicloud_elb_member.test]
}

locals {
  ip_version = huaweicloud_elb_member.test.ip_version
}

output "ip_version_filter_is_useful" {
  value = length(data.huaweicloud_elb_members.ip_version_filter.members) > 0 && alltrue(
  [for v in data.huaweicloud_elb_members.ip_version_filter.members[*].ip_version : v == local.ip_version]
  )  
}

data "huaweicloud_elb_members" "operating_status_filter" {
  pool_id          = huaweicloud_elb_pool.test.id
  operating_status = huaweicloud_elb_member.test.operating_status
  depends_on       = [huaweicloud_elb_member.test]
}

locals {
  operating_status = huaweicloud_elb_member.test.operating_status
}

output "operating_status_filter_is_useful" {
  value = length(data.huaweicloud_elb_members.operating_status_filter.members) > 0 && alltrue(
  [for v in data.huaweicloud_elb_members.operating_status_filter.members[*].operating_status : v == local.operating_status]
  )  
}

data "huaweicloud_elb_members" "member_type_filter" {
  pool_id     = huaweicloud_elb_pool.test.id
  member_type = huaweicloud_elb_member.test.member_type
  depends_on  = [huaweicloud_elb_member.test]
}

locals {
  member_type = huaweicloud_elb_member.test.member_type
}

output "member_type_filter_is_useful" {
  value = length(data.huaweicloud_elb_members.member_type_filter.members) > 0 && alltrue(
  [for v in data.huaweicloud_elb_members.member_type_filter.members[*].member_type : v == local.member_type]
  )  
}
`, testAccDatasourceMembers_base(name), name)
}
