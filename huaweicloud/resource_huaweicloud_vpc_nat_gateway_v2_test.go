package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/huawei-clouds/golangsdk/openstack/vpc/v2/natgateways"
)

func TestAccVpcNatGateway_basic(t *testing.T) {
	var network networks.Network
	var router routers.Router
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcV2NatGatewayDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpcV2NatGateway_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2RouterExists("huaweicloud_networking_router_v2.router_1", &router),
					testAccCheckNetworkingV2RouterInterfaceExists("huaweicloud_networking_router_interface_v2.int_1"),
					testAccCheckVpcV2NatGatewayExists("huaweicloud_vpc_nat_gateway_v2.nat_1"),
				),
			},
			resource.TestStep{
				Config: testAccVpcV2NatGateway_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_vpc_nat_gateway_v2.nat_1", "name", "nat_1_updated"),
					resource.TestCheckResourceAttr("huaweicloud_vpc_nat_gateway_v2.nat_1", "description", "nat_1 updated"),
					resource.TestCheckResourceAttr("huaweicloud_vpc_nat_gateway_v2.nat_1", "spec", "2"),
				),
			},
		},
	})
}

func testAccCheckVpcV2NatGatewayDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	vpcClient, err := config.vpcV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud vpc client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_vpc_nat_gateway_v2" {
			continue
		}

		_, err := natgateways.Get(vpcClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Nat gateway still exists")
		}
	}

	return nil
}

func testAccCheckVpcV2NatGatewayExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		vpcClient, err := config.vpcV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud vpc client: %s", err)
		}

		found, err := natgateways.Get(vpcClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Nat gateway not found")
		}

		return nil
	}
}

const testAccVpcV2NatGateway_basic = `
resource "huaweicloud_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_router_interface_v2" "int_1" {
  subnet_id = "${huaweicloud_networking_subnet_v2.subnet_1.id}"
  router_id = "${huaweicloud_networking_router_v2.router_1.id}"
}

resource "huaweicloud_vpc_nat_gateway_v2" "nat_1" {
  name   = "nat_1"
  description = "test for terraform"
  spec = "1"
  internal_network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  router_id = "${huaweicloud_networking_router_v2.router_1.id}"
  depends_on = ["huaweicloud_networking_router_interface_v2.int_1"]
}
`

const testAccVpcV2NatGateway_update = `
resource "huaweicloud_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_router_interface_v2" "int_1" {
  subnet_id = "${huaweicloud_networking_subnet_v2.subnet_1.id}"
  router_id = "${huaweicloud_networking_router_v2.router_1.id}"
}

resource "huaweicloud_vpc_nat_gateway_v2" "nat_1" {
  name   = "nat_1_updated"
  description = "nat_1 updated"
  spec = "2"
  internal_network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  router_id = "${huaweicloud_networking_router_v2.router_1.id}"
  depends_on = ["huaweicloud_networking_router_interface_v2.int_1"]
}
`
