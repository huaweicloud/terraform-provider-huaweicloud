package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKmsCancelKeyDeletion_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare a KMS Key ID in the scheduled deletion state, then config it to the environment variable.
			acceptance.TestAccPreCheckKmsKeyID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsCancelKeyDeletion_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_kms_cancel_key_deletion.test", "key_id", acceptance.HW_KMS_KEY_ID),
					resource.TestCheckResourceAttrSet("huaweicloud_kms_cancel_key_deletion.test", "key_state"),
				),
			},
		},
	})
}

func testAccKmsCancelKeyDeletion_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_cancel_key_deletion" "test" {
  key_id = "%s"
}
`, acceptance.HW_KMS_KEY_ID)
}
