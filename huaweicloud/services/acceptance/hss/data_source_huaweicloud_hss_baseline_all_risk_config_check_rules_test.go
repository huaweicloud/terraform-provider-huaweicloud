package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBaselineAllRiskConfigCheckRules_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_baseline_all_risk_config_check_rules.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test requires ensuring the existence of baseline check data under HSS service.
			acceptance.TestAccPreCheckHSSCheckRuleId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBaselineAllRiskConfigCheckRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.tag"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.check_rule_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.check_rule_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.severity"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.check_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.check_type_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.standard"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.host_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.failed_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.scan_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.statistics_scan_result"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.enable_fix"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.enable_click"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.cancel_ignore_enable_click"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.rule_params.#"),

					resource.TestCheckOutput("is_check_type_filter_useful", "true"),
					resource.TestCheckOutput("is_standard_filter_useful", "true"),
					resource.TestCheckOutput("is_statistics_scan_result_filter_useful", "true"),
					resource.TestCheckOutput("is_check_rule_name_filter_useful", "true"),
					resource.TestCheckOutput("is_severity_filter_useful", "true"),
					resource.TestCheckOutput("is_tag_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceBaselineAllRiskConfigCheckRules_basic() string {
	return `
data "huaweicloud_hss_baseline_all_risk_config_check_rules" "test" {
  enterprise_project_id = "all_granted_eps"
  statistics_flag       = true
}

# Filter using check_type.
locals {
  check_type = data.huaweicloud_hss_baseline_all_risk_config_check_rules.test.data_list[0].check_type
}

data "huaweicloud_hss_baseline_all_risk_config_check_rules" "check_type_filter" {
  enterprise_project_id = "all_granted_eps"
  check_type            = local.check_type
}

output "is_check_type_filter_useful" {
  value = length(data.huaweicloud_hss_baseline_all_risk_config_check_rules.check_type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_baseline_all_risk_config_check_rules.check_type_filter.data_list[*].check_type
    : v == local.check_type]
  )
}

# Filter using standard.
locals {
  standard = data.huaweicloud_hss_baseline_all_risk_config_check_rules.test.data_list[0].standard
}

data "huaweicloud_hss_baseline_all_risk_config_check_rules" "standard_filter" {
  enterprise_project_id = "all_granted_eps"
  standard              = local.standard
}

output "is_standard_filter_useful" {
  value = length(data.huaweicloud_hss_baseline_all_risk_config_check_rules.standard_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_baseline_all_risk_config_check_rules.standard_filter.data_list[*].standard : v == local.standard]
  )
}

# Filter using statistics_scan_result.
locals {
  statistics_scan_result = data.huaweicloud_hss_baseline_all_risk_config_check_rules.test.data_list[0].statistics_scan_result
}

data "huaweicloud_hss_baseline_all_risk_config_check_rules" "statistics_scan_result_filter" {
  enterprise_project_id  = "all_granted_eps"
  statistics_scan_result = local.statistics_scan_result
}

output "is_statistics_scan_result_filter_useful" {
  value = length(data.huaweicloud_hss_baseline_all_risk_config_check_rules.statistics_scan_result_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_baseline_all_risk_config_check_rules.statistics_scan_result_filter.data_list[*].statistics_scan_result
    : v == local.statistics_scan_result]
  )
}

# Filter using check_rule_name.
locals {
  check_rule_name = data.huaweicloud_hss_baseline_all_risk_config_check_rules.test.data_list[0].check_rule_name
}

data "huaweicloud_hss_baseline_all_risk_config_check_rules" "check_rule_name_filter" {
  enterprise_project_id = "all_granted_eps"
  check_rule_name       = local.check_rule_name
}

output "is_check_rule_name_filter_useful" {
  value = length(data.huaweicloud_hss_baseline_all_risk_config_check_rules.check_rule_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_baseline_all_risk_config_check_rules.check_rule_name_filter.data_list[*].check_rule_name
    : v == local.check_rule_name]
  )
}

# Filter using severity.
locals {
  severity = data.huaweicloud_hss_baseline_all_risk_config_check_rules.test.data_list[0].severity
}

data "huaweicloud_hss_baseline_all_risk_config_check_rules" "severity_filter" {
  enterprise_project_id = "all_granted_eps"
  severity              = local.severity
}

output "is_severity_filter_useful" {
  value = length(data.huaweicloud_hss_baseline_all_risk_config_check_rules.severity_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_baseline_all_risk_config_check_rules.severity_filter.data_list[*].severity : v == local.severity]
  )
}

# Filter using tag.
locals {
  tag = data.huaweicloud_hss_baseline_all_risk_config_check_rules.test.data_list[0].tag
}

data "huaweicloud_hss_baseline_all_risk_config_check_rules" "tag_filter" {
  enterprise_project_id = "all_granted_eps"
  tag                   = local.tag
}

output "is_tag_filter_useful" {
  value = length(data.huaweicloud_hss_baseline_all_risk_config_check_rules.tag_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_baseline_all_risk_config_check_rules.tag_filter.data_list[*].tag : v == local.tag]
  )
}
`
}
