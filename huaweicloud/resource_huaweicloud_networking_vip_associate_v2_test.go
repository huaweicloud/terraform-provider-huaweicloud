package huaweicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/ports"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccNetworkingV2VIPAssociate_basic(t *testing.T) {
	rand := acctest.RandString(5)
	var vip ports.Port
	var port1 ports.Port
	var port2 ports.Port

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2VIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2VIPAssociateConfig_basic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("huaweicloud_networking_port.port_1", &port1),
					testAccCheckNetworkingV2PortExists("huaweicloud_networking_port.port_2", &port2),
					testAccCheckNetworkingV2VIPExists("huaweicloud_networking_vip.vip_1", &vip),
					testAccCheckNetworkingV2VIPAssociated(&port1, &vip),
					testAccCheckNetworkingV2VIPAssociated(&port2, &vip),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2VIPAssociateDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV2Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_networking_vip_associate" {
			continue
		}

		vipID := rs.Primary.Attributes["vip_id"]
		_, err = ports.Get(networkingClient, vipID).Extract()
		if err != nil {
			// If the error is a 404, then the vip port does not exist,
			// and therefore the floating IP cannot be associated to it.
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}
	}

	log.Printf("[DEBUG] testAccCheckNetworkingV2VIPAssociateDestroy success!")
	return nil
}

func testAccCheckNetworkingV2VIPAssociated(p *ports.Port, vip *ports.Port) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*config.Config)
		networkingClient, err := config.NetworkingV2Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
		}

		p, err := ports.Get(networkingClient, p.ID).Extract()
		if err != nil {
			// If the error is a 404, then the port does not exist,
			// and therefore the VIP cannot be associated to it.
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}

		vipport, err := ports.Get(networkingClient, vip.ID).Extract()
		if err != nil {
			// If the error is a 404, then the vip port does not exist,
			// and therefore the VIP cannot be associated to it.
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}

		for _, ip := range p.FixedIPs {
			for _, addresspair := range vipport.AllowedAddressPairs {
				if ip.IPAddress == addresspair.IPAddress {
					log.Printf("[DEBUG] testAccCheckNetworkingV2VIPAssociated success!")
					return nil
				}
			}
		}

		return fmt.Errorf("VIP %s was not attached to port %s", vipport.ID, p.ID)
	}
}

func testAccNetworkingV2VIPAssociateConfig_basic(rand string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "vpc_1" {
  name = "acc-test-vpc-%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "subnet_1" {
  vpc_id     = huaweicloud_vpc.vpc_1.id
  name       = "acc-test-subnet-%s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}

resource "huaweicloud_networking_port" "port_1" {
  name       = "port_1"
  network_id = huaweicloud_vpc_subnet.subnet_1.id

  fixed_ip {
    subnet_id = huaweicloud_vpc_subnet.subnet_1.subnet_id
  }
}

resource "huaweicloud_networking_port" "port_2" {
  name       = "port_2"
  network_id = huaweicloud_vpc_subnet.subnet_1.id

  fixed_ip {
    subnet_id = huaweicloud_vpc_subnet.subnet_1.subnet_id
  }
}

resource "huaweicloud_networking_vip" "vip_1" {
  network_id = huaweicloud_vpc_subnet.subnet_1.id
}

resource "huaweicloud_networking_vip_associate" "vip_associate_1" {
  vip_id   = huaweicloud_networking_vip.vip_1.id
  port_ids = [huaweicloud_networking_port.port_1.id, huaweicloud_networking_port.port_2.id]
}
	`, rand, rand)
}
