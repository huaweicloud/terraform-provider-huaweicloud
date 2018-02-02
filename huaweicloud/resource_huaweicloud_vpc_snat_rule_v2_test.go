package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/huawei-clouds/golangsdk/openstack/vpc/v2/snatrules"
)

func TestAccVpcSnatRule_basic(t *testing.T) {
	var fip floatingips.FloatingIP
	var network networks.Network
	var router routers.Router
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcV2SnatRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpcV2SnatRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2RouterExists("huaweicloud_networking_router_v2.router_1", &router),
					testAccCheckNetworkingV2FloatingIPExists("huaweicloud_networking_floatingip_v2.fip_1", &fip),
					testAccCheckNetworkingV2RouterInterfaceExists("huaweicloud_networking_router_interface_v2.int_1"),
					testAccCheckVpcV2NatGatewayExists("huaweicloud_vpc_nat_gateway_v2.nat_1"),
					testAccCheckVpcV2SnatRuleExists("huaweicloud_vpc_snat_rule_v2.snat_1"),
				),
			},
		},
	})
}

func testAccCheckVpcV2SnatRuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	vpcClient, err := config.vpcV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud vpc client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_vpc_snat_rule_v2" {
			continue
		}

		_, err := snatrules.Get(vpcClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Snat rule still exists")
		}
	}

	return nil
}

func testAccCheckVpcV2SnatRuleExists(n string) resource.TestCheckFunc {
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

		found, err := snatrules.Get(vpcClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Snat rule not found")
		}

		return nil
	}
}

const testAccVpcV2SnatRule_basic = `
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

resource "huaweicloud_networking_floatingip_v2" "fip_1" {
}

resource "huaweicloud_vpc_nat_gateway_v2" "nat_1" {
  name   = "nat_1"
  description = "test for terraform"
  spec = "1"
  internal_network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  router_id = "${huaweicloud_networking_router_v2.router_1.id}"
  depends_on = ["huaweicloud_networking_router_interface_v2.int_1"]
}

resource "huaweicloud_vpc_snat_rule_v2" "snat_1" {
  nat_gateway_id = "${huaweicloud_vpc_nat_gateway_v2.nat_1.id}"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  floating_ip_id = "${huaweicloud_networking_floatingip_v2.fip_1.id}"
}
`
