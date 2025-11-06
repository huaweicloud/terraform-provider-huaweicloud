package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKpsExportPrivateKey_basic(t *testing.T) {
	rscName := "huaweicloud_kps_export_private_key.test"

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please ensure that there is at least one available keypair in the test environment.
			acceptance.TestAccPreCheckKPSKeyPairName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKpsExportPrivateKey_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rscName, "private_key"),
				),
			},
		},
	})
}

func testAccKpsExportPrivateKey_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_kps_export_private_key" "test" {
  name = "%s"
}
`, acceptance.HW_KPS_KEYPAIR_NAME)
}
