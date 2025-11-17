package cdn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDomainOwnerVerify_basic(t *testing.T) {
	// Avoid CheckDestroy, because there is nothing in the resource destroy method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCDN(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainOwnerVerify_basic(),
			},
			{
				Config: testAccDomainOwnerVerify_basic_step1(),
			},
		},
	})
}

func testAccDomainOwnerVerify_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain_owner_verify" "test" {
  domain_name = "%[1]s"
}
`, acceptance.HW_CDN_DOMAIN_NAME)
}

func testAccDomainOwnerVerify_basic_step1() string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain_owner_verify" "test_with_dns" {
  domain_name = "%[1]s"
  verify_type = "dns"
}
`, acceptance.HW_CDN_DOMAIN_NAME)
}
