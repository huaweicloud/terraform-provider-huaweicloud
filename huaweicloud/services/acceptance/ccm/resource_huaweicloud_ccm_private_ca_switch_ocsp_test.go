package ccm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourcePrivateCaSwitchOcsp_basic(t *testing.T) {
	// This resource does not support Read and Destroy, so CheckExist and CheckDestroy are ignored.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Prepare private CA in the runtime environment before running test cases.
			acceptance.TestAccPreCheckCCMPrivateCaID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePrivateCaSwitchOcsp_basic(),
			},
			{
				Config: testAccResourcePrivateCaSwitchOcsp_update(),
			},
		},
	})
}

func testAccResourcePrivateCaSwitchOcsp_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_private_ca_switch_ocsp" "test" {
  ca_id       = "%s"
  ocsp_switch = true
}
`, acceptance.HW_CCM_PRIVATE_CA_ID)
}

func testAccResourcePrivateCaSwitchOcsp_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_private_ca_switch_ocsp" "test" {
  ca_id       = "%[1]s"
  ocsp_switch = false
}

resource "huaweicloud_ccm_private_ca_switch_ocsp" "test_duplicate_situation" {
  depends_on = [huaweicloud_ccm_private_ca_switch_ocsp.test]

  ca_id       = "%[1]s"
  ocsp_switch = false
}
`, acceptance.HW_CCM_PRIVATE_CA_ID)
}
