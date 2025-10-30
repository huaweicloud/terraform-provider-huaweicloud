package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGenerateMac_basic(t *testing.T) {
	rName := "huaweicloud_kms_generate_mac.test"
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
				Config: testAccGenerateMac_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "mac"),
				),
			},
		},
	})
}

func testAccGenerateMac_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_generate_mac" "test" {
  key_id        = "%s"
  mac_algorithm = "HMAC_SHA_384"
  message       = "dGVzdGFiY2RlZmc="
}
`, acceptance.HW_KMS_KEY_ID)
}
