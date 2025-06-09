package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKmsDecryptDatakey_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare a KMS Key ID and cipher text, then config it to the environment variable.
			acceptance.TestAccPreCheckKmsKeyID(t)
			acceptance.TestAccPreCheckKmsKeyCiphertext(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsDecryptDatakey_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_kms_decrypt_datakey.test", "key_id", acceptance.HW_KMS_KEY_ID),
					resource.TestCheckResourceAttrSet("huaweicloud_kms_decrypt_datakey.test", "data_key"),
					resource.TestCheckResourceAttrSet("huaweicloud_kms_decrypt_datakey.test", "datakey_length"),
					resource.TestCheckResourceAttrSet("huaweicloud_kms_decrypt_datakey.test", "datakey_dgst"),
				),
			},
		},
	})
}

func testAccKmsDecryptDatakey_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_decrypt_datakey" "test" {
  key_id                = "%[1]s"
  cipher_text           = "%[2]s"
  datakey_cipher_length = "%[3]s"
}
`, acceptance.HW_KMS_KEY_ID, acceptance.HW_KMS_KEY_CIPHER_TEXT, acceptance.HW_KMS_KEY_CIPHER_TEXT_LEN)
}
