package antiddos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGeoIpRulesDataSource_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_aad_geoip_rules.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGeoIpRulesDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.#"),
				),
			},
		},
	})
}

func testAccGeoIpRulesDataSource_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_aad_geoip_rules" "test" {
  domain_name   = "%s"
  overseas_type = "0"
}
`, acceptance.HW_DOMAIN_NAME)
}
