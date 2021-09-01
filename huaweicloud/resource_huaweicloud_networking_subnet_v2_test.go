package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v2/subnets"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccNetworkingV2Subnet_basic(t *testing.T) {
	var subnet subnets.Subnet

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"huaweicloud_networking_subnet_v2.subnet_1", "allocation_pools.0.start", "192.168.199.100"),
				),
			},
			{
				ResourceName:      "huaweicloud_networking_subnet_v2.subnet_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkingV2Subnet_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"huaweicloud_networking_subnet_v2.subnet_1", "name", "subnet_1"),
					resource.TestCheckResourceAttr(
						"huaweicloud_networking_subnet_v2.subnet_1", "gateway_ip", "192.168.199.1"),
					resource.TestCheckResourceAttr(
						"huaweicloud_networking_subnet_v2.subnet_1", "enable_dhcp", "true"),
					resource.TestCheckResourceAttr(
						"huaweicloud_networking_subnet_v2.subnet_1", "allocation_pools.0.start", "192.168.199.150"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_enableDHCP(t *testing.T) {
	var subnet subnets.Subnet

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_enableDHCP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"huaweicloud_networking_subnet_v2.subnet_1", "enable_dhcp", "true"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_noGateway(t *testing.T) {
	var subnet subnets.Subnet

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_noGateway,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"huaweicloud_networking_subnet_v2.subnet_1", "gateway_ip", ""),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_impliedGateway(t *testing.T) {
	var subnet subnets.Subnet

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_impliedGateway,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"huaweicloud_networking_subnet_v2.subnet_1", "gateway_ip", "192.168.199.1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Subnet_timeout(t *testing.T) {
	var subnet subnets.Subnet

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2SubnetDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV2Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_networking_subnet_v2" {
			continue
		}

		_, err := subnets.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Subnet still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2SubnetExists(n string, subnet *subnets.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		networkingClient, err := config.NetworkingV2Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
		}

		found, err := subnets.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Subnet not found")
		}

		*subnet = *found

		return nil
	}
}

const testAccNetworkingV2Subnet_basic = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"

  allocation_pools {
    start = "192.168.199.100"
    end = "192.168.199.200"
  }
}
`

const testAccNetworkingV2Subnet_update = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  gateway_ip = "192.168.199.1"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"

  allocation_pools {
    start = "192.168.199.150"
    end = "192.168.199.200"
  }
}
`

const testAccNetworkingV2Subnet_enableDHCP = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  gateway_ip = "192.168.199.1"
  enable_dhcp = true
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}
`

const testAccNetworkingV2Subnet_noGateway = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  no_gateway = true
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}
`

const testAccNetworkingV2Subnet_impliedGateway = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}
resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}
`

const testAccNetworkingV2Subnet_timeout = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"

  allocation_pools {
    start = "192.168.199.100"
    end = "192.168.199.200"
  }

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
