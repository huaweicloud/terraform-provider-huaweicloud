package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVerifyMac_basic(t *testing.T) {
	rName := "huaweicloud_kms_verify_mac.test"
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare a KMS Key with `key_usage` value is **GENERATE_VERIFY_MAC**.
			acceptance.TestAccPreCheckKmsKeyID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVerifyMac_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "mac_valid"),

					resource.TestCheckOutput("is_valid", "true"),
				),
			},
		},
	})
}

func testAccVerifyMac_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_kms_verify_mac" "test" {
  key_id        = "%[2]s"
  mac_algorithm = huaweicloud_kms_generate_mac.test.mac_algorithm
  message       = huaweicloud_kms_generate_mac.test.message
  mac           = huaweicloud_kms_generate_mac.test.mac
}

output "is_valid" {
  value = (huaweicloud_kms_verify_mac.test.mac_valid == true)
}
`, testAccGenerateMac_basic(), acceptance.HW_KMS_KEY_ID)
}
