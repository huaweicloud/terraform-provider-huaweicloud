package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// PASS
func TestAccRTSStackV1_importBasic(t *testing.T) {
	resourceName := "huaweicloud_rts_stack_v1.stack_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRTSStackV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRTSStackV1_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
