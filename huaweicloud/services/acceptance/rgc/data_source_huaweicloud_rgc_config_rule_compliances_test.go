package rgc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceConfigRuleCompliances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_config_rule_compliances.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRGCAccountID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConfigRuleCompliances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "config_rule_compliances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "config_rule_compliances.0.rule_name"),
					resource.TestCheckResourceAttrSet(dataSource, "config_rule_compliances.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "config_rule_compliances.0.region"),
					resource.TestCheckResourceAttrSet(dataSource, "config_rule_compliances.0.control_id"),
				),
			},
		},
	})
}

func testAccDataSourceConfigRuleCompliances_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rgc_config_rule_compliances" "test" {
  managed_account_id = "%[1]s"
}
`, acceptance.HW_RGC_ACCOUNT_ID)
}
