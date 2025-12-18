package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBaselineRiskConfigCheckRules_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_baseline_risk_config_check_rules.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running test, prepare a host with host enterprise edition protection enabled.
			// Configure corresponding policy and manual check.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBaselineRiskConfigCheckRules_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.severity"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.check_rule_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.check_rule_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.standard"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.check_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.scan_result"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.status"),

					resource.TestCheckOutput("severity_filter_useful", "true"),
					resource.TestCheckOutput("result_type_filter_useful", "true"),
					resource.TestCheckOutput("eps_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceBaselineRiskConfigCheckRules_basic = `
data "huaweicloud_hss_baseline_risk_config_check_rules" "test" {
  check_name = "SSH"
  standard   = "hw_standard"
}

locals {
  severity = data.huaweicloud_hss_baseline_risk_config_check_rules.test.data_list[0].severity
}

data "huaweicloud_hss_baseline_risk_config_check_rules" "severity_filter" {
  check_name = "SSH"
  standard   = "hw_standard"
  severity   = local.severity
}

output "severity_filter_useful" {
  value = length(data.huaweicloud_hss_baseline_risk_config_check_rules.severity_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_baseline_risk_config_check_rules.severity_filter.data_list[*].severity : v == local.severity]
  )
}

data "huaweicloud_hss_baseline_risk_config_check_rules" "result_type_filter" {	
  check_name  = "SSH"
  standard    = "hw_standard"
  result_type = "unhandled"
}

output "result_type_filter_useful" {
  value = length(data.huaweicloud_hss_baseline_risk_config_check_rules.result_type_filter.data_list) > 0
}

data "huaweicloud_hss_baseline_risk_config_check_rules" "eps_filter" {
  check_name            = "SSH"
  standard              = "hw_standard"
  enterprise_project_id = "0"
}

output "eps_filter_useful" {
  value = length(data.huaweicloud_hss_baseline_risk_config_check_rules.eps_filter.data_list) > 0
}
`
