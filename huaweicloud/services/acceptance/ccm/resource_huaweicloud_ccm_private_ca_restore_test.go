package ccm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCaRestore_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please configure a private CA ID which in `DELETED` status.
			acceptance.TestAccPreCheckCCMPrivateCaID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceResourceCaRestore_basic(),
			},
		},
	})
}

func testResourceResourceCaRestore_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_private_ca_restore" "test" {
  ca_id = "%s"
}
`, acceptance.HW_CCM_PRIVATE_CA_ID)
}
