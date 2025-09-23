package sfsturbo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceSFSTurboChangeChargeMode_basic(t *testing.T) {
	// Avoid CheckDestroy because this resource is a one-time action resource and there is nothing in the destroy method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running test, please prepare a SFS Turbo file system which is postpaid.
			acceptance.TestAccPrecheckSFSTurboShareId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurboChangeChargeMode_basic(),
			},
		},
	})
}

func testAccSFSTurboChangeChargeMode_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_sfs_turbo_change_charge_mode" "test" {
  share_id      = "%s"
  period_num    = 1
  period_type   = 2
  is_auto_renew = 0
}
`, acceptance.HW_SFS_TURBO_SHARE_ID)
}
