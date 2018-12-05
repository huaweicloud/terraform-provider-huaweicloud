package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccCTSTrackerV1_importBasic(t *testing.T) {
	resourceName := "huaweicloud_cts_tracker_v1.tracker_v1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCTSTrackerV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCTSTrackerV1_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
