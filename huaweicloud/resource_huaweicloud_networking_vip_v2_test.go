package huaweicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/ports"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccNetworkingV2VIP_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	updateName := rName + "update"
	resourceName := "huaweicloud_networking_vip.vip_1"
	var vip ports.Port

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2VIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2VIPConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2VIPExists(resourceName, &vip),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkingV2VIPConfig_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2VIPDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
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

		config := testAccProvider.Meta().(*config.Config)
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

func testAccNetworkingV2VIPConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "vpc_1" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "subnet_1" {
  vpc_id      = huaweicloud_vpc.vpc_1.id
  name        = "%s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
}

resource "huaweicloud_networking_vip" "vip_1" {
  name       = "%s"
  network_id = huaweicloud_vpc_subnet.subnet_1.id
}
	`, rName, rName, rName)
}

func testAccNetworkingV2VIPConfig_update(updateName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "vpc_1" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "subnet_1" {
  vpc_id      = huaweicloud_vpc.vpc_1.id
  name        = "%s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
}

resource "huaweicloud_networking_vip" "vip_1" {
  name       = "%s"
  network_id = huaweicloud_vpc_subnet.subnet_1.id
}
	`, updateName, updateName, updateName)
}
