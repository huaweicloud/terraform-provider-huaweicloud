package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v1/eips"
)

func TestAccVpcV1EIP_basic(t *testing.T) {
	var eip eips.PublicIp

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcV1EIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1EIP_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1EIPExists("huaweicloud_vpc_eip_v1.eip_1", &eip),
				),
			},
		},
	})
}

func TestAccVpcV1EIP_share(t *testing.T) {
	var eip eips.PublicIp

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcV1EIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcV1EIP_share,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcV1EIPExists("huaweicloud_vpc_eip_v1.eip_1", &eip),
				),
			},
		},
	})
}

func testAccCheckVpcV1EIPDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating EIP: %s", err)
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

func testAccCheckVpcV1EIPExists(n string, kp *eips.PublicIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.networkingV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating networking client: %s", err)
		}

		found, err := eips.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("EIP not found")
		}

		kp = &found

		return nil
	}
}

const testAccVpcV1EIP_basic = `
resource "huaweicloud_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name = "test"
    size = 8
    share_type = "PER"
    charge_mode = "traffic"
  }
}
`

const testAccVpcV1EIP_share = `
resource "huaweicloud_vpc_bandwidth_v2" "bandwidth_1" {
	name = "bandwidth_1"
	size = 5
}

resource "huaweicloud_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    id = "${huaweicloud_vpc_bandwidth_v2.bandwidth_1.id}"
    share_type = "WHOLE"
  }
}
`
