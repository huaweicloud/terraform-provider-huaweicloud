package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKeyUpdatePrimaryRegion_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare a Multi-region KMS key.
			acceptance.TestAccPreCheckKmsKeyID(t)
			acceptance.TestAccPreCheckKmsKeyPrimaryRegion(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyUpdatePrimaryRegion_basic(),
			},
		},
	})
}

func testAccKeyUpdatePrimaryRegion_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key_update_primary_region" "test" {
  key_id         = "%[1]s"
  primary_region = "%[2]s"
}
`, acceptance.HW_KMS_KEY_ID, acceptance.HW_KMS_KEY_PRIMARY_REGION)
}
