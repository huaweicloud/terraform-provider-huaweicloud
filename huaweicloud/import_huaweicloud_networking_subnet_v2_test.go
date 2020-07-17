package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNetworkingV2Subnet_importBasic(t *testing.T) {
	resourceName := "huaweicloud_networking_subnet_v2.subnet_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Subnet_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
