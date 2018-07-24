package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOTCVpcRouteV2_importBasic(t *testing.T) {
	resourceName := "huaweicloud_vpc_route_v2.route_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouteV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRouteV2_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
