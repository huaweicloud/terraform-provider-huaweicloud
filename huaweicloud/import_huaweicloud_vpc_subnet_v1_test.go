package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOTCVpcSubnetV1_importBasic(t *testing.T) {
	resourceName := "huaweicloud_vpc_subnet_v1.subnet_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpcSubnetV1_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
