package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRetireGrant_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckKmsKeyID(t)
			acceptance.TestAccPreCheckKmsGrantID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRetireGrant_basic(),
			},
		},
	})
}

func testAccRetireGrant_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_retire_grant" "test" {
  key_id   = "%[1]s"
  grant_id = "%[2]s"
  sequence = "919c82d4-8046-4722-9094-35c3c6524cff"
}
`, acceptance.HW_KMS_KEY_ID, acceptance.HW_KMS_GRANT_ID)
}
