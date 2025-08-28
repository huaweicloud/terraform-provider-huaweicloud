package antiddos

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Note: Due to the lack of a test environment, this test case only verifies the expected error message.
func TestAccAadDomainCertificate_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAadDomainID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testDomainCertificate_basic(),
				ExpectError: regexp.MustCompile(`证书不存在。`),
			},
		},
	})
}

func testDomainCertificate_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_aad_domain_certificate" "test" {
  domain_id = "%[1]s"
  op_type   = 1
  cert_name = "test_cert_name"
}
`, acceptance.HW_AAD_DOMAIN_ID)
}
