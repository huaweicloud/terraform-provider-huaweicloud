package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEvsVolumeRetype_basic(t *testing.T) {
	// Avoid CheckDestroy because this resource is a one-time action resource and there is nothing in the destroy
	// method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test needs to create an EVS SAS volume before running.
			acceptance.TestAccPreCheckEVSVolumeID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// One-time action resource do not need to be checked and no processing is performed in the Read method.
				Config: testAccEvsVolumeRetype_basic(),
			},
		},
	})
}

func testAccEvsVolumeRetype_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_evs_volume_retype" "test1" {
  volume_id = "%s"
  new_type  = "ESSD"
}
`, acceptance.HW_EVS_VOLUME_ID)
}
