package dbss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAuditRiskRules_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_dbss_audit_risk_rules.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byName   = "data.huaweicloud_dbss_audit_risk_rules.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byRiskLevel   = "data.huaweicloud_dbss_audit_risk_rules.filter_by_risk_level"
		dcByRiskLevel = acceptance.InitDataSourceCheck(byRiskLevel)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDbssInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAuditRiskRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.feature"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.rank"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.risk_level"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_useful", "true"),

					dcByRiskLevel.CheckResourceExists(),
					resource.TestCheckOutput("risk_level_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceAuditRiskRules_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dbss_audit_risk_rules" "test" {
  instance_id = "%[1]s"
}

locals {
  name = data.huaweicloud_dbss_audit_risk_rules.test.rules[0].name
}

data "huaweicloud_dbss_audit_risk_rules" "filter_by_name" {
  instance_id = "%[1]s"
  name        = local.name
}

output "name_filter_useful" {
  value = length(data.huaweicloud_dbss_audit_risk_rules.filter_by_name.rules) > 0 && alltrue(
    [for v in data.huaweicloud_dbss_audit_risk_rules.filter_by_name.rules[*].name : v == local.name]
  )
}

locals {
  risk_level = data.huaweicloud_dbss_audit_risk_rules.test.rules[0].risk_level
}

data "huaweicloud_dbss_audit_risk_rules" "filter_by_risk_level" {
  instance_id = "%[1]s"
  risk_level  = local.risk_level
}

output "risk_level_filter_useful" {
  value = length(data.huaweicloud_dbss_audit_risk_rules.filter_by_risk_level.rules) > 0 && alltrue(
    [for v in data.huaweicloud_dbss_audit_risk_rules.filter_by_risk_level.rules[*].risk_level : v == local.risk_level]
  )
}
`, acceptance.HW_DBSS_INSATNCE_ID)
}
