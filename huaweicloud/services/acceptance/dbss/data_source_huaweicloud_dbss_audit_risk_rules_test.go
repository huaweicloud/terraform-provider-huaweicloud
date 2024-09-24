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
		name           = acceptance.RandomAccResourceName()
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byName   = "data.huaweicloud_dbss_audit_risk_rules.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byRiskLevel   = "data.huaweicloud_dbss_audit_risk_rules.filter_by_risk_level"
		dcByRiskLevel = acceptance.InitDataSourceCheck(byRiskLevel)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAuditRiskRules_basic(name),
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

func testDataSourceAuditRiskRules_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dbss_audit_risk_rules" "test" {
  instance_id = huaweicloud_dbss_instance.test.instance_id
}

locals {
  name = data.huaweicloud_dbss_audit_risk_rules.test.rules[0].name
}

data "huaweicloud_dbss_audit_risk_rules" "filter_by_name" {
  instance_id = huaweicloud_dbss_instance.test.instance_id
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
  instance_id = huaweicloud_dbss_instance.test.instance_id
  risk_level  = local.risk_level
}

output "risk_level_filter_useful" {
  value = length(data.huaweicloud_dbss_audit_risk_rules.filter_by_risk_level.rules) > 0 && alltrue(
    [for v in data.huaweicloud_dbss_audit_risk_rules.filter_by_risk_level.rules[*].risk_level : v == local.risk_level]
  )
}
`, testInstance_basic(name))
}
