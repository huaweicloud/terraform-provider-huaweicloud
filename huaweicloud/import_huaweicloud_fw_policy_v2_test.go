package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccFWPolicyV2_importBasic(t *testing.T) {
	resourceName := "huaweicloud_fw_policy_v2.policy_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWPolicyV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFWPolicyV2_addRules,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
