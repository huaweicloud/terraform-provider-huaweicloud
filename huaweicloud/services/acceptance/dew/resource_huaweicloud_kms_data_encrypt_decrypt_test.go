package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEncryptDecrypt_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	// Avoid CheckDestroy because this resource is a action resource and there is nothing in the destroy method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// The action resource do not need to be checked and no processing is performed in the Read method.
				Config: testAccDataEncryptDecrypt_encrypt(rName),
			},
			{
				// The action resource do not need to be checked and no processing is performed in the Read method.
				Config: testAccDataEncryptDecrypt_decrypt(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("is_equals", "true"),
				),
			},
		},
	})
}

func testAccDataEncryptDecrypt_encrypt(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_kms_data_encrypt_decrypt" "test" {
  key_id     = huaweicloud_kms_key.test.id
  action     = "encrypt"
  plain_text = "abc"
}
`, testAccKmsKey_basic(name))
}

func testAccDataEncryptDecrypt_decrypt(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_kms_data_encrypt_decrypt" "try" {
  action      = "decrypt"
  cipher_text = huaweicloud_kms_data_encrypt_decrypt.test.cipher_data
}

output "is_equals" {
  value = (huaweicloud_kms_data_encrypt_decrypt.test.plain_text == huaweicloud_kms_data_encrypt_decrypt.try.plain_data)
}
`, testAccDataEncryptDecrypt_encrypt(name))
}
