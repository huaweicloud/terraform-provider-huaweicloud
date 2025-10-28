package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRsaDatakeyPair_basic(t *testing.T) {
	rName := "huaweicloud_kms_rsa_datakey_pair.test"
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare a KMS Key ID
			acceptance.TestAccPreCheckKmsKeyID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRsaDatakeyPair_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "public_key"),
					resource.TestCheckResourceAttrSet(rName, "private_key_cipher_text"),
					resource.TestCheckResourceAttrSet(rName, "private_key_plain_text"),
				),
			},
		},
	})
}

func testAccRsaDatakeyPair_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_rsa_datakey_pair" "test" {
  key_id                        = "%s"
  key_spec                      = "RSA_3072"
  with_plain_text               = true
  additional_authenticated_data = "terraform test"
  sequence                      = "919c82d4-8046-4722-9094-35c3c6524cff"
}
`, acceptance.HW_KMS_KEY_ID)
}
