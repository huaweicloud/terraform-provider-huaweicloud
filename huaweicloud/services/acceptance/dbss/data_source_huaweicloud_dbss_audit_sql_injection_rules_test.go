package dbss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAuditSqlInjectionRules_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_dbss_audit_sql_injection_rules.test"
		name           = acceptance.RandomAccResourceName()
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byRiskLevels   = "data.huaweicloud_dbss_audit_sql_injection_rules.filter_by_risk_levels"
		dcByRiskLevels = acceptance.InitDataSourceCheck(byRiskLevels)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAuditSqlInjectionRules_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.risk_level"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.feature"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.rank"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.regex"),

					dcByRiskLevels.CheckResourceExists(),
					resource.TestCheckOutput("risk_levels_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceAuditSqlInjectionRules_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dbss_audit_sql_injection_rules" "test" {
  instance_id = huaweicloud_dbss_instance.test.instance_id
}

locals {
  risk_levels = data.huaweicloud_dbss_audit_sql_injection_rules.test.rules[0].risk_level
}

data "huaweicloud_dbss_audit_sql_injection_rules" "filter_by_risk_levels" {
  instance_id = huaweicloud_dbss_instance.test.instance_id
  risk_levels = local.risk_levels
}

output "risk_levels_filter_useful" {
  value = length(data.huaweicloud_dbss_audit_sql_injection_rules.filter_by_risk_levels.rules) > 0 && alltrue(
    [for v in data.huaweicloud_dbss_audit_sql_injection_rules.filter_by_risk_levels.rules[*].risk_level : 
    v == local.risk_levels]
  )
}
`, testInstance_basic(name))
}
