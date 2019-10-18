package huaweicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/ports"
)

func TestAccNetworkingV2VIPAssociate_basic(t *testing.T) {
	var vip ports.Port
	var port1 ports.Port
	var port2 ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2VIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccNetworkingV2VIPAssociateConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("huaweicloud_networking_port_v2.port_1", &port1),
					testAccCheckNetworkingV2PortExists("huaweicloud_networking_port_v2.port_2", &port2),
					testAccCheckNetworkingV2VIPExists("huaweicloud_networking_vip_v2.vip_1", &vip),
					testAccCheckNetworkingV2VIPAssociateAssociated(&port1, &vip),
					testAccCheckNetworkingV2VIPAssociateAssociated(&port2, &vip),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2VIPAssociateDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_networking_vip_associate_v2" {
			continue
		}

		vipid, portids, err := parseNetworkingVIPAssociateID(rs.Primary.ID)
		if err != nil {
			return err
		}

		vipport, err := ports.Get(networkingClient, vipid).Extract()
		if err != nil {
			// If the error is a 404, then the vip port does not exist,
			// and therefore the floating IP cannot be associated to it.
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}

		// port by port
		for _, portid := range portids {
			p, err := ports.Get(networkingClient, portid).Extract()
			if err != nil {
				// If the error is a 404, then the port does not exist,
				// and therefore the floating IP cannot be associated to it.
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return nil
				}
				return err
			}

			// But if the port and vip still exists
			for _, ip := range p.FixedIPs {
				for _, addresspair := range vipport.AllowedAddressPairs {
					if ip.IPAddress == addresspair.IPAddress {
						return fmt.Errorf("VIP %s is still associated to port %s", vipid, portid)
					}
				}
			}
		}
	}

	log.Printf("[DEBUG] testAccCheckNetworkingV2VIPAssociateDestroy success!")
	return nil
}

func testAccCheckNetworkingV2VIPAssociateAssociated(p *ports.Port, vip *ports.Port) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
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
					log.Printf("[DEBUG] testAccCheckNetworkingV2VIPAssociateAssociated success!")
					return nil
				}
			}
		}

		return fmt.Errorf("VIP %s was not attached to port %s", vipport.ID, p.ID)
	}
}

var TestAccNetworkingV2VIPAssociateConfig_basic = fmt.Sprintf(`
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_router_interface_v2" "router_interface_1" {
  router_id = "${huaweicloud_networking_router_v2.router_1.id}"
  subnet_id = "${huaweicloud_networking_subnet_v2.subnet_1.id}"
}

resource "huaweicloud_networking_router_v2" "router_1" {
  name = "router_1"
  external_network_id = "%s"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  
  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
  }
}

resource "huaweicloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
	  
  network {
    port = "${huaweicloud_networking_port_v2.port_1.id}"
  }
}

resource "huaweicloud_networking_port_v2" "port_2" {
  name = "port_2"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  
  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
  }
}

resource "huaweicloud_compute_instance_v2" "instance_2" {
  name = "instance_2"
  security_groups = ["default"]
  availability_zone = "%s"
	  
  network {
    port = "${huaweicloud_networking_port_v2.port_2.id}"
  }
}

resource "huaweicloud_networking_vip_v2" "vip_1" {
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  subnet_id = "${huaweicloud_networking_subnet_v2.subnet_1.id}"
}

resource "huaweicloud_networking_vip_associate_v2" "vip_associate_1" {
  vip_id = "${huaweicloud_networking_vip_v2.vip_1.id}"
  port_ids = ["${huaweicloud_networking_port_v2.port_1.id}", "${huaweicloud_networking_port_v2.port_2.id}"]
}
`, OS_EXTGW_ID, OS_AVAILABILITY_ZONE, OS_AVAILABILITY_ZONE)
