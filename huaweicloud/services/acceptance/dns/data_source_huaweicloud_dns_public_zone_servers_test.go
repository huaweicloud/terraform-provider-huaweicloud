package dns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// The HW_DNS_ZONE_NAMES provided must be purchased from a domain registrar.
func TestAccDataPublicZoneServers_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_dns_public_zone_servers.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDnsZoneNames(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPublicZoneServers_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dcName, "include_hw_dns", "true"),
					resource.TestMatchResourceAttr(dcName, "dns_servers.#", regexp.MustCompile(`^[1-9][0-9]*$`)),
					resource.TestMatchResourceAttr(dcName, "expected_dns_servers.#", regexp.MustCompile(`^[1-9][0-9]*$`)),
				),
			},
		},
	})
}

func testAccDataPublicZoneServers_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dns_public_zone_servers" "test" {
  domain_name = split(",", "%[1]s")[0]
}`, acceptance.HW_DNS_ZONE_NAMES)
}
