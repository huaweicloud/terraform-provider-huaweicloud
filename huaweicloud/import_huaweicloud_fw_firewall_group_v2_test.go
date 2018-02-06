package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccFWFirewallV2_importBasic(t *testing.T) {
	resourceName := "huaweicloud_fw_firewall_group_v2.fw_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWFirewallGroupV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFWFirewallGroupV2_basic_1,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
