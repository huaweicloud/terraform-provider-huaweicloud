package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running the test case, please ensure that there is at least one WAF instance in the current region.
// Due to testing environment limitations, this test case can only test the scenario with empty `items`.
func TestAccDataSourceRulesThreatIntelligence_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_waf_rules_threat_intelligence.test"
		rName          = acceptance.RandomAccResourceName()
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRulesThreatIntelligence_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.#"),
				),
			},
		},
	})
}

func testDataSourceRulesThreatIntelligence_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_policy" "test" {
  name                  = "%[1]s"
  level                 = 1
  enterprise_project_id = "0"
}
`, name)
}

func testDataSourceRulesThreatIntelligence_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_waf_rules_threat_intelligence" "test" {
  policy_id             = huaweicloud_waf_policy.test.id
  enterprise_project_id = "0"
}

`, testDataSourceRulesThreatIntelligence_base(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
