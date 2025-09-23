package aom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceMultiAccountAggregationRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aom_multi_account_aggregation_rules.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccountAggregationRuleEnable(t)
			acceptance.TestAccPrecheckDomainId(t)
			acceptance.TestAccPrecheckDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceMultiAccountAggregationRules_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.accounts.#"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.services.#"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.send_to_source_account"),
				),
			},
		},
	})
}

func testDataSourceMultiAccountAggregationRules_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_aom_multi_account_aggregation_rules" "test" {
  depends_on = [huaweicloud_aom_multi_account_aggregation_rule.test]

  enterprise_project_id = huaweicloud_aom_multi_account_aggregation_rule.test.enterprise_project_id
}
`, testMultiAccountAggregationRule_basic(name))
}
