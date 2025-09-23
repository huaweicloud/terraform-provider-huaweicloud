package antiddos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWebProtectionPolicies_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aad_web_protection_policies.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare an AAD domain name and set it to an environment variable.
			acceptance.TestAccPrecheckAadDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceWebProtectionPolicies_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "level"),
					resource.TestCheckResourceAttrSet(dataSource, "mode"),
					resource.TestCheckResourceAttrSet(dataSource, "options.#"),
					resource.TestCheckResourceAttrSet(dataSource, "options.0.cc"),
					resource.TestCheckResourceAttrSet(dataSource, "options.0.custom"),
					resource.TestCheckResourceAttrSet(dataSource, "options.0.geoip"),
					resource.TestCheckResourceAttrSet(dataSource, "options.0.modulex_enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "options.0.whiteblackip"),
				),
			},
		},
	})
}

func testDataSourceWebProtectionPolicies_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_aad_web_protection_policies" "test" {
  domain_name   = "%s"
  overseas_type = 0
}
`, acceptance.HW_AAD_DOMAIN_NAME)
}
