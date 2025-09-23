package ccm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// This test case will cause resource residue, please call it with caution
// The resource is a one-time action resource and there is nothing in the destroy method.
func TestAccResourceCertificateApply_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCMSSLCertificateId(t)
			acceptance.TestAccPreCheckCCMEnableFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceResourceCertificateApply_basic(),
			},
		},
	})
}

func testResourceResourceCertificateApply_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate_apply" "test" {
  certificate_id  = "%s"
  domain          = "www.example.com"
  applicant_name  = "Emily"
  applicant_phone = "13212345678"
  applicant_email = "example@huawei.com"
  domain_method   = "dns"
}
`, acceptance.HW_CCM_SSL_CERTIFICATE_ID)
}
