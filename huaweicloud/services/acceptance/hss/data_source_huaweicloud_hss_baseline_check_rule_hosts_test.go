package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBaselineCheckRuleHosts_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_baseline_check_rule_hosts.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSCheckRuleId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBaselineCheckRuleHosts_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.host_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.check_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.baseline_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.host_private_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.scan_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.failed_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.passed_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.enable_fix"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.enable_verify"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.enable_click"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.cancel_ignore_enable_click"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.result_type"),

					resource.TestCheckOutput("is_check_name_filter_useful", "true"),
					resource.TestCheckOutput("is_result_type_filter_useful", "true"),
					resource.TestCheckOutput("is_host_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceBaselineCheckRuleHosts_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_baseline_check_rule_hosts" "test" {
  check_rule_id         = "%[1]s"
  standard              = "hw_standard"
  enterprise_project_id = "all_granted_eps"
}

# Filter using check_name.
locals {
  check_name = data.huaweicloud_hss_baseline_check_rule_hosts.test.data_list[0].check_name
}

data "huaweicloud_hss_baseline_check_rule_hosts" "check_name_filter" {
  check_rule_id         = "%[1]s"
  standard              = "hw_standard"
  enterprise_project_id = "all_granted_eps"
  check_name            = local.check_name
}

output "is_check_name_filter_useful" {
  value = length(data.huaweicloud_hss_baseline_check_rule_hosts.check_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_baseline_check_rule_hosts.check_name_filter.data_list[*].check_name : v == local.check_name]
  )
}

# Filter using result_type.
locals {
  result_type = data.huaweicloud_hss_baseline_check_rule_hosts.test.data_list[0].result_type
}

data "huaweicloud_hss_baseline_check_rule_hosts" "result_type_filter" {
  check_rule_id         = "%[1]s"
  standard              = "hw_standard"
  enterprise_project_id = "all_granted_eps"
  result_type           = local.result_type
}

output "is_result_type_filter_useful" {
  value = length(data.huaweicloud_hss_baseline_check_rule_hosts.result_type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_baseline_check_rule_hosts.result_type_filter.data_list[*].result_type : v == local.result_type]
  )
}

# Filter using host_name.
locals {
  host_name = data.huaweicloud_hss_baseline_check_rule_hosts.test.data_list[0].host_name
}

data "huaweicloud_hss_baseline_check_rule_hosts" "host_name_filter" {
  check_rule_id         = "%[1]s"
  standard              = "hw_standard"
  enterprise_project_id = "all_granted_eps"
  host_name             = local.host_name
}

output "is_host_name_filter_useful" {
  value = length(data.huaweicloud_hss_baseline_check_rule_hosts.host_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_baseline_check_rule_hosts.host_name_filter.data_list[*].host_name : v == local.host_name]
  )
}
`, acceptance.HW_HSS_CHECK_RULE_ID)
}
