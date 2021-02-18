package huaweicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/ports"
)

func TestAccNetworkingV2VIP_basic(t *testing.T) {
	rand := acctest.RandString(5)
	resourceName := "huaweicloud_networking_vip.vip_1"
	var vip ports.Port

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2VIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2VIPConfig_basic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2VIPExists(resourceName, &vip),
					resource.TestCheckResourceAttr(resourceName, "name", "acc-test-vip"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkingV2VIPConfig_update(rand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "acc-test-vip-updated"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2VIPDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_networking_vip" {
			continue
		}

		_, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("VIP still exists")
		}
	}

	log.Printf("[DEBUG] testAccCheckNetworkingV2VIPDestroy success!")

	return nil
}

func testAccCheckNetworkingV2VIPExists(n string, vip *ports.Port) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.NetworkingV2Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
		}

		found, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("VIP not found")
		}
		log.Printf("[DEBUG] test found is: %#v", found)
		*vip = *found

		return nil
	}
}

func testAccNetworkingV2VIPConfig_basic(rand string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "vpc_1" {
  name = "acc-test-vpc-%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "subnet_1" {
  vpc_id      = huaweicloud_vpc.vpc_1.id
  name        = "acc-test-subnet-%s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
}

resource "huaweicloud_networking_vip" "vip_1" {
  name       = "acc-test-vip"
  network_id = huaweicloud_vpc_subnet.subnet_1.id
}
	`, rand, rand)
}

func testAccNetworkingV2VIPConfig_update(rand string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "vpc_1" {
  name = "acc-test-vpc-%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "subnet_1" {
  vpc_id      = huaweicloud_vpc.vpc_1.id
  name        = "acc-test-subnet-%s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
}

resource "huaweicloud_networking_vip" "vip_1" {
  name       = "acc-test-vip-updated"
  network_id = huaweicloud_vpc_subnet.subnet_1.id
}
	`, rand, rand)
}
