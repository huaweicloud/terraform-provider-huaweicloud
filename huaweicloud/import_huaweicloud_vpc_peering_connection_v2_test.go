package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOTCVpcPeeringConnectionV1_importBasic(t *testing.T) {
	resourceName := "huaweicloud_vpc_peering_connection_v2.peering_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcPeeringConnectionV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpcPeeringConnectionV2_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
