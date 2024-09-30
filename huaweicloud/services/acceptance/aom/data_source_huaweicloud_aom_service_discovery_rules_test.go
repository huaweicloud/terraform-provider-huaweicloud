package aom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceServiceDiscoveryRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aom_service_discovery_rules.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceServiceDiscoveryRules_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "rules.#"),

					resource.TestCheckOutput("rule_id_filter_validation", "true"),
				),
			},
		},
	})
}

func testDataSourceServiceDiscoveryRules_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_aom_service_discovery_rules" "test" {}

data "huaweicloud_aom_service_discovery_rules" "id_filter" {
  rule_id = huaweicloud_aom_service_discovery_rule.test.rule_id
}

output "rule_id_filter_validation" {
  value = alltrue([for v in data.huaweicloud_aom_service_discovery_rules.id_filter.rules[*].id : 
    v == huaweicloud_aom_service_discovery_rule.test.rule_id])
}
`, testAOMServiceDiscoveryRule_basic(name))
}
