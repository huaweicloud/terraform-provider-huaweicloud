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
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("member_id_filter_is_useful", "true"),
					resource.TestCheckOutput("address_filter_is_useful", "true"),
					resource.TestCheckOutput("protocol_port_filter_is_useful", "true"),
					resource.TestCheckOutput("weight_filter_is_useful", "true"),
					resource.TestCheckOutput("subnet_id_filter_is_useful", "true"),
					resource.TestCheckOutput("member_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccElbMemberConfig_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

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
  address       = "192.168.0.10"
  protocol_port = 8080
  pool_id       = huaweicloud_elb_pool.test.id
  subnet_id     = huaweicloud_vpc_subnet.test.ipv4_subnet_id
}
`, common.TestVpc(name), name)
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

data "huaweicloud_elb_members" "member_type_filter" {
  pool_id     = huaweicloud_elb_pool.test.id
  member_type = "instance"
  depends_on  = [huaweicloud_elb_member.test]
}

output "member_type_filter_is_useful" {
  value = length(data.huaweicloud_elb_members.member_type_filter.members) > 0 && alltrue(
  [for v in data.huaweicloud_elb_members.member_type_filter.members[*].member_type : v == "instance"]
  )  
}
`, testAccElbMemberConfig_basic(name), name)
}
