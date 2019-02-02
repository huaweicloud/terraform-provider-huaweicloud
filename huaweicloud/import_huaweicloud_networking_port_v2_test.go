package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNetworkingV2Port_importBasic(t *testing.T) {
	resourceName := "huaweicloud_networking_port_v2.port_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Port_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"fixed_ip",
				},
			},
		},
	})
}

func TestAccNetworkingV2Port_importDHCPOpts(t *testing.T) {
	resourceName := "huaweicloud_networking_port_v2.port_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Port_createExtraDHCPOpts,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"fixed_ip",
				},
			},
		},
	})
}
