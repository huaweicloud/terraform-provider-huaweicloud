package dns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccZoneAuthorization_basic(t *testing.T) {
	var resourceName = "huaweicloud_dns_zone_authorization.test"

	// Avoid CheckDestroy because there is nothing in the destroy method for this resource.
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
				Config: testAccZoneAuthorization_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "zone_name"),
					resource.TestCheckResourceAttrSet(resourceName, "second_level_zone_name"),
					resource.TestCheckResourceAttrSet(resourceName, "record.#"),
					resource.TestCheckResourceAttrSet(resourceName, "record.0.host"),
					resource.TestCheckResourceAttrSet(resourceName, "record.0.value"),
					resource.TestMatchResourceAttr(resourceName, "status", regexp.MustCompile(`(CREATED|VERIFIED)`)),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(resourceName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccZoneAuthorization_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_zone_authorization" "test" {
  zone_name = format("dev.%%s", try(split(",", "%[1]s")[0], "terraform.example.com"))
}
`, acceptance.HW_DNS_ZONE_NAMES)
}
