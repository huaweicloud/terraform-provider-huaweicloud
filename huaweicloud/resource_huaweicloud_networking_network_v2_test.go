package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/layer3/routers"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/networks"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/subnets"
)

func TestAccNetworkingV2Network_basic(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Network_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Network_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"huaweicloud_networking_network_v2.network_1", "name", "network_2"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_netstack(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Network_netstack,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2RouterExists("huaweicloud_networking_router_v2.router_1", &router),
					testAccCheckNetworkingV2RouterInterfaceExists(
						"huaweicloud_networking_router_interface_v2.ri_1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_timeout(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Network_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_multipleSegmentMappings(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Network_multipleSegmentMappings,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2NetworkDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_networking_network_v2" {
			continue
		}

		_, err := networks.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Network still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2NetworkExists(n string, network *networks.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
		}

		found, err := networks.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Network not found")
		}

		*network = *found

		return nil
	}
}

const testAccNetworkingV2Network_basic = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}
`

const testAccNetworkingV2Network_update = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_2"
  admin_state_up = "true"
}
`

const testAccNetworkingV2Network_netstack = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.10.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_router_v2" "router_1" {
  name = "router_1"
}

resource "huaweicloud_networking_router_interface_v2" "ri_1" {
  router_id = "${huaweicloud_networking_router_v2.router_1.id}"
  subnet_id = "${huaweicloud_networking_subnet_v2.subnet_1.id}"
}
`

const testAccNetworkingV2Network_timeout = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

const testAccNetworkingV2Network_multipleSegmentMappings = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  segments =[
    {
      segmentation_id = 2,
      network_type = "vxlan"
    }
  ],
  admin_state_up = "true"
}
`
