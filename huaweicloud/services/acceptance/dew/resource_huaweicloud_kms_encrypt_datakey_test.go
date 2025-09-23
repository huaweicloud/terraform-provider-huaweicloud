package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKmsEncryptDatakey_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare a KMS Key ID and one valid plaintext, then config it to the environment variable.
			acceptance.TestAccPreCheckKmsKeyID(t)
			acceptance.TestAccPreCheckKmsKeyPlaintext(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsEncryptDatakey_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_kms_encrypt_datakey.test", "key_id", acceptance.HW_KMS_KEY_ID),
					resource.TestCheckResourceAttrSet("huaweicloud_kms_encrypt_datakey.test", "cipher_text"),
					resource.TestCheckResourceAttr("huaweicloud_kms_encrypt_datakey.test", "datakey_plain_length",
						acceptance.HW_KMS_KEY_PLAINTEXT_LEN),
				),
			},
		},
	})
}

func testAccKmsEncryptDatakey_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_encrypt_datakey" "test" {
  key_id               = "%[1]s"
  plain_text           = "%[2]s"
  datakey_plain_length = "%[3]s"
}
`, acceptance.HW_KMS_KEY_ID, acceptance.HW_KMS_KEY_PLAINTEXT, acceptance.HW_KMS_KEY_PLAINTEXT_LEN)
}
