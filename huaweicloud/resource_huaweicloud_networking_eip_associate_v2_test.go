package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v1/eips"
)

func TestAccNetworkingV2EIPAssociate_basic(t *testing.T) {
	var eip eips.PublicIp

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_networking_eip_associate_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2EIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2EIPAssociate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1EIPExists("huaweicloud_vpc_eip_v1.test", &eip),
					resource.TestCheckResourceAttrPtr(
						resourceName, "floating_ip", &eip.PublicAddress),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2EIPAssociateDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating EIP Client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_vpc_eip_v1" {
			continue
		}

		_, err := eips.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("EIP still exists")
		}
	}

	return nil
}

func testAccNetworkingV2EIPAssociate_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_v1" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet_v1" "test" {
  name          = "%s"
  cidr          = "192.168.0.0/16"
  gateway_ip    = "192.168.0.1"
  vpc_id        = huaweicloud_vpc_v1.test.id
}

resource "huaweicloud_networking_port_v2" "test" {
  network_id = huaweicloud_vpc_subnet_v1.test.id

  fixed_ip {
    subnet_id  = huaweicloud_vpc_subnet_v1.test.subnet_id
    ip_address = "192.168.0.20"
  }
}

resource "huaweicloud_vpc_eip_v1" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%s"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_networking_eip_associate_v2" "test" {
  floating_ip = huaweicloud_vpc_eip_v1.test.address
  port_id     = huaweicloud_networking_port_v2.test.id
}
`, rName, rName, rName)
}
