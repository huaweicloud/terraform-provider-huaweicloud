package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKmsDatakeyWithoutPlaintext_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare a KMS Key ID, then config it to the environment variable.
			acceptance.TestAccPreCheckKmsKeyID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsDatakeyWithoutPlaintext_datakeyLengthOnly(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_kms_datakey_without_plaintext.test1", "key_id", acceptance.HW_KMS_KEY_ID),
					resource.TestCheckResourceAttr("huaweicloud_kms_datakey_without_plaintext.test1", "datakey_length", "64"),
					resource.TestCheckResourceAttrSet("huaweicloud_kms_datakey_without_plaintext.test1", "cipher_text"),
				),
			},
			{
				Config: testAccKmsDatakeyWithoutPlaintext_keySpecOnly(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_kms_datakey_without_plaintext.test2", "key_id", acceptance.HW_KMS_KEY_ID),
					resource.TestCheckResourceAttr("huaweicloud_kms_datakey_without_plaintext.test2", "key_spec", "AES_256"),
					resource.TestCheckResourceAttrSet("huaweicloud_kms_datakey_without_plaintext.test2", "cipher_text"),
				),
			},
			{
				Config: testAccKmsDatakeyWithoutPlaintext_both(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_kms_datakey_without_plaintext.test3", "key_id", acceptance.HW_KMS_KEY_ID),
					resource.TestCheckResourceAttr("huaweicloud_kms_datakey_without_plaintext.test3", "datakey_length", "128"),
					resource.TestCheckResourceAttr("huaweicloud_kms_datakey_without_plaintext.test3", "key_spec", "AES_128"),
					resource.TestCheckResourceAttrSet("huaweicloud_kms_datakey_without_plaintext.test3", "cipher_text"),
				),
			},
			{
				Config: testAccKmsDatakeyWithoutPlaintext_none(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_kms_datakey_without_plaintext.test4", "key_id", acceptance.HW_KMS_KEY_ID),
					resource.TestCheckResourceAttrSet("huaweicloud_kms_datakey_without_plaintext.test4", "cipher_text"),
				),
			},
			{
				Config: testAccKmsDatakeyWithoutPlaintext_datakeyLengthAndSeq(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_kms_datakey_without_plaintext.test5", "key_id", acceptance.HW_KMS_KEY_ID),
					resource.TestCheckResourceAttr("huaweicloud_kms_datakey_without_plaintext.test5", "datakey_length", "256"),
					resource.TestCheckResourceAttr("huaweicloud_kms_datakey_without_plaintext.test5", "sequence",
						"12345678-1234-1234-1234-123456789abc"),
					resource.TestCheckResourceAttrSet("huaweicloud_kms_datakey_without_plaintext.test5", "cipher_text"),
				),
			},
			{
				Config: testAccKmsDatakeyWithoutPlaintext_keySpecAndSeq(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_kms_datakey_without_plaintext.test6", "key_id", acceptance.HW_KMS_KEY_ID),
					resource.TestCheckResourceAttr("huaweicloud_kms_datakey_without_plaintext.test6", "key_spec", "AES_128"),
					resource.TestCheckResourceAttr("huaweicloud_kms_datakey_without_plaintext.test6", "sequence",
						"abcdefab-1234-5678-9abc-abcdefabcdef"),
					resource.TestCheckResourceAttrSet("huaweicloud_kms_datakey_without_plaintext.test6", "cipher_text"),
				),
			},
		},
	})
}

func testAccKmsDatakeyWithoutPlaintext_datakeyLengthOnly() string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_datakey_without_plaintext" "test1" {
  key_id         = "%s"
  datakey_length = "64"
}
`, acceptance.HW_KMS_KEY_ID)
}

func testAccKmsDatakeyWithoutPlaintext_keySpecOnly() string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_datakey_without_plaintext" "test2" {
  key_id   = "%s"
  key_spec = "AES_256"
}
`, acceptance.HW_KMS_KEY_ID)
}

func testAccKmsDatakeyWithoutPlaintext_both() string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_datakey_without_plaintext" "test3" {
  key_id         = "%s"
  datakey_length = "128"
  key_spec       = "AES_128"
}
`, acceptance.HW_KMS_KEY_ID)
}

func testAccKmsDatakeyWithoutPlaintext_none() string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_datakey_without_plaintext" "test4" {
  key_id = "%s"
}
`, acceptance.HW_KMS_KEY_ID)
}

func testAccKmsDatakeyWithoutPlaintext_datakeyLengthAndSeq() string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_datakey_without_plaintext" "test5" {
  key_id         = "%s"
  datakey_length = "256"
  sequence       = "12345678-1234-1234-1234-123456789abc"
}
`, acceptance.HW_KMS_KEY_ID)
}

func testAccKmsDatakeyWithoutPlaintext_keySpecAndSeq() string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_datakey_without_plaintext" "test6" {
  key_id   = "%s"
  key_spec = "AES_128"
  sequence = "abcdefab-1234-5678-9abc-abcdefabcdef"
}
`, acceptance.HW_KMS_KEY_ID)
}
