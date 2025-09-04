package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKmsSign_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare a KMS Key ID and one message based on Base64 encoding, then config it to the environment variable.
			acceptance.TestAccPreCheckKmsKeyID(t)
			acceptance.TestAccPreCheckKmsKeyMessage(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsSign_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_kms_sign.test", "key_id", acceptance.HW_KMS_KEY_ID),
					resource.TestCheckResourceAttrSet("huaweicloud_kms_sign.test", "signature"),
				),
			},
		},
	})
}

func testAccKmsSign_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_sign" "test" {
  key_id            = "%[1]s"
  message           = "%[2]s"
  signing_algorithm = "RSASSA_PSS_SHA_256"
  message_type      = "RAW"
}
`, acceptance.HW_KMS_KEY_ID, acceptance.HW_KMS_KEY_MESSAGE)
}
