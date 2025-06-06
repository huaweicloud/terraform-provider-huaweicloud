package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDomainAssociateCertificate_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running the test case, please ensure that there is at least one WAF instance in the current region.
			// Prepare a WAF certificate and a WAF domain.
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckWafCertID(t)
			acceptance.TestAccPreCheckWafDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainAssociateCertificate_basic(),
			},
		},
	})
}

func testAccDomainAssociateCertificate_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_domain_associate_certificate" "test" {
  certificate_id = "%[1]s"
  cloud_host_ids = ["%[2]s"]
}
`, acceptance.HW_WAF_CERTIFICATE_ID, acceptance.HW_WAF_DOMAIN_ID)
}
