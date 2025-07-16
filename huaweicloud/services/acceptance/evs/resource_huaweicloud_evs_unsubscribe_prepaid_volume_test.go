package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEvsUnsubscribePrepaidVolume_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test needs to create an EVS prepaid volume before running.
			acceptance.TestAccPreCheckEVSVolumeID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsUnsubscribePrepaidVolume_basic(),
			},
		},
	})
}

func testAccEvsUnsubscribePrepaidVolume_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_evs_unsubscribe_prepaid_volume" "test" {
  volume_ids = ["%s"]
}
`, acceptance.HW_EVS_VOLUME_ID)
}
