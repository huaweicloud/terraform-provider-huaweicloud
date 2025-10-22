package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRecycleBinVolumeDelete_basic(t *testing.T) {
	// Avoid CheckDestroy because this resource is a one-time action resource and there is nothing in the destroy
	// method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a volume ID that is in the recycle bin.
			acceptance.TestAccPreCheckEVSRecycleBinVolumeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// One-time action resource do not need to be checked and no processing is performed in the Read method.
				Config: testRecycleBinVolumeDelete_basic(),
			},
		},
	})
}

func testRecycleBinVolumeDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_evs_recycle_bin_volume_delete" "test" {
  volume_id = "%s"
}
`, acceptance.HW_EVS_RECYCLE_BIN_VOLUME_ID)
}
