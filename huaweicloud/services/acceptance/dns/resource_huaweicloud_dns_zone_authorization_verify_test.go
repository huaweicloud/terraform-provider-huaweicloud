package dns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccZoneAuthorizationVerify_basic(t *testing.T) {
	var resourceName = "huaweicloud_dns_zone_authorization_verify.test"

	// Avoid CheckDestroy because this resource is a one-time action resource and there is nothing in the destroy
	// method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDnsZoneNames(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccZoneAuthorizationVerify_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "status", regexp.MustCompile(`(CREATED|VERIFIED)`)),
				),
			},
		},
	})
}

func testAccZoneAuthorizationVerify_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_zone_authorization" "test" {
  zone_name = format("dev.%%s", try(split(",", "%[1]s")[0], "terraform.example.com"))
}

resource "huaweicloud_dns_zone_authorization_verify" "test" {
  authorization_id = huaweicloud_dns_zone_authorization.test.id
}
`, acceptance.HW_DNS_ZONE_NAMES)
}
