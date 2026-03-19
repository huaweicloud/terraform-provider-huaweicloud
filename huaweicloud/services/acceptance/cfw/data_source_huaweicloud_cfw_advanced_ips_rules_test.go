package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAdvancedIpsRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_advanced_ips_rules.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAdvancedIpsRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "advanced_ips_rules.0.action"),
					resource.TestCheckResourceAttrSet(dataSource, "advanced_ips_rules.0.ips_rule_id"),
					resource.TestCheckResourceAttrSet(dataSource, "advanced_ips_rules.0.ips_rule_type"),
					resource.TestCheckResourceAttrSet(dataSource, "advanced_ips_rules.0.param"),
					resource.TestCheckResourceAttrSet(dataSource, "advanced_ips_rules.0.status"),
				),
			},
		},
	})
}

func testDataSourceAdvancedIpsRules_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cfw_advanced_ips_rules" "test" {
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
}
`, testAccDatasourceFirewalls_basic())
}
