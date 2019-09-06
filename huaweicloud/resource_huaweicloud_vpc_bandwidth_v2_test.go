package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v1/bandwidths"
)

func TestAccVpcBandWidthV2_basic(t *testing.T) {
	var bandwidth bandwidths.BandWidth

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcBandWidthV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcBandWidthV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcBandWidthV2Exists("huaweicloud_vpc_bandwidth_v2.bandwidth_1", &bandwidth),
					resource.TestCheckResourceAttr(
						"huaweicloud_vpc_bandwidth_v2.bandwidth_1", "name", "bandwidth_1"),
				),
			},
			{
				Config: testAccVpcBandWidthV2_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcBandWidthV2Exists("huaweicloud_vpc_bandwidth_v2.bandwidth_1", &bandwidth),
					resource.TestCheckResourceAttr(
						"huaweicloud_vpc_bandwidth_v2.bandwidth_1", "name", "bandwidth_1_updated"),
					resource.TestCheckResourceAttr(
						"huaweicloud_vpc_bandwidth_v2.bandwidth_1", "size", "6"),
				),
			},
		},
	})
}

func testAccCheckVpcBandWidthV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_vpc_bandwidth_v2" {
			continue
		}

		_, err := bandwidths.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("BandWidth still exists")
		}
	}

	return nil
}

func testAccCheckVpcBandWidthV2Exists(n string, bandwidth *bandwidths.BandWidth) resource.TestCheckFunc {
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
			return fmt.Errorf("Error creating huaweicloud networking client: %s", err)
		}

		found, err := bandwidths.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("bandwidth not found")
		}

		*bandwidth = found

		return nil
	}
}

const testAccVpcBandWidthV2_basic = `
resource "huaweicloud_vpc_bandwidth_v2" "bandwidth_1" {
	name = "bandwidth_1"
	size = 5
}`

const testAccVpcBandWidthV2_update = `
resource "huaweicloud_vpc_bandwidth_v2" "bandwidth_1" {
	name = "bandwidth_1_updated"
	size = 6
}`
