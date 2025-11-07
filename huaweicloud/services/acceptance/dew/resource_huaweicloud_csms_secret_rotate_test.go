package dew

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to a lack of testing conditions, only failure scenarios are being tested for the time being.
// Therefore, this test case cannot be guaranteed to be 100% reliable.
func TestAccCsmsSecretRotate_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCsmsSecretName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCsmsSecretRotate_basic(),
				ExpectError: regexp.MustCompile(`error waiting for DEW CSMS secret rotation task to be completed`),
			},
		},
	})
}

func testAccCsmsSecretRotate_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_csms_secret_rotate" "test" {
  secret_name = "%s"
}
`, acceptance.HW_CSMS_SECRET_NAME)
}
