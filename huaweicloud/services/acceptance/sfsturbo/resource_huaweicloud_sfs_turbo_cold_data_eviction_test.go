package sfsturbo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceColdDataEviction_basic(t *testing.T) {
	// Avoid CheckDestroy because this resource is a one-time action resource and there is nothing in the destroy method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running test, please prepare a SFS Turbo file system which bound a storage backend.
			acceptance.TestAccPrecheckSFSTurboShareId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccColdDataEviction_basic(),
			},
		},
	})
}

func testAccColdDataEviction_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_sfs_turbo_cold_data_eviction" "test" {
  share_id = "%s"
  action   = "config_gc_time"
  gc_time  = 10
}
`, acceptance.HW_SFS_TURBO_SHARE_ID)
}
