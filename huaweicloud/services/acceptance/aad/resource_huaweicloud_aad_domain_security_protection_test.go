package antiddos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceDomainSecurityProtection_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare an AAD domain ID and set it to an environment variable.
			acceptance.TestAccPreCheckAadDomainID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDomainSecurityProtection_basic(),
			},
		},
	})
}

func testDomainSecurityProtection_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_aad_domain_security_protection" "test" {
  domain_id  = "%s"
  waf_switch = 0
  cc_switch  = 0
}
`, acceptance.HW_AAD_DOMAIN_ID)
}
