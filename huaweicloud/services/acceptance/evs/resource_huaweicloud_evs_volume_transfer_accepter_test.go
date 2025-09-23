package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVolumeTransferAccepter_basic(t *testing.T) {
	// Avoid CheckDestroy because this resource is a one-time action resource and there is nothing in the destroy
	// method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckEVSTransferAccepter(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// One-time action resource do not need to be checked and no processing is performed in the Read method.
				Config: testAccVolumeTransferAccepter_basic(),
			},
		},
	})
}

func testAccVolumeTransferAccepter_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_evs_volume_transfer_accepter" "test" { 
  transfer_id = "%[1]s"
  auth_key    = "%[2]s"
}
`, acceptance.HW_EVS_TRANSFER_ID, acceptance.HW_EVS_TRANSFER_AUTH_KEY)
}
