package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBaselineSecurityChecksDetails_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_baseline_security_checks_details.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBaselineSecurityChecksDetails_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "check_details.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "check_details.0.key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "check_details.0.check_rule_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "check_details.0.check_rule_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "check_details.0.check_rule_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "check_details.0.check_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "check_details.0.severity"),
					resource.TestCheckResourceAttrSet(dataSourceName, "check_details.0.level"),
					resource.TestCheckResourceAttrSet(dataSourceName, "check_details.0.checked"),
					resource.TestCheckResourceAttrSet(dataSourceName, "check_details.0.rule_params.#"),

					resource.TestCheckOutput("check_type_filter_is_useful", "true"),
					resource.TestCheckOutput("check_rule_name_filter_is_useful", "true"),
					resource.TestCheckOutput("checked_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceBaselineSecurityChecksDetails_basic() string {
	return `
data "huaweicloud_hss_baseline_security_checks_details" "test" {
  support_os            = "Linux"
  standard              = "hw_standard"
  enterprise_project_id = "all_granted_eps"
  severity              = "Low"
}

# Filter using check_type.
locals {
  check_type = data.huaweicloud_hss_baseline_security_checks_details.test.check_details[0].check_type
}

data "huaweicloud_hss_baseline_security_checks_details" "check_type_filter" {
  support_os = "Linux"
  standard   = "hw_standard"
  check_type = local.check_type
}

output "check_type_filter_is_useful" {
  value = length(data.huaweicloud_hss_baseline_security_checks_details.check_type_filter.check_details) > 0 && alltrue(
    [for v in data.huaweicloud_hss_baseline_security_checks_details.check_type_filter.check_details[*].check_type : v == local.check_type]
  )
}

# Filter using check_rule_name.
locals {
  check_rule_name = data.huaweicloud_hss_baseline_security_checks_details.test.check_details[0].check_rule_name
}

data "huaweicloud_hss_baseline_security_checks_details" "check_rule_name_filter" {
  support_os      = "Linux"
  standard        = "hw_standard"
  check_rule_name = local.check_rule_name
}

output "check_rule_name_filter_is_useful" {
  value = length(data.huaweicloud_hss_baseline_security_checks_details.check_rule_name_filter.check_details) > 0 && alltrue(
    [for v in data.huaweicloud_hss_baseline_security_checks_details.check_rule_name_filter.check_details[*].check_rule_name : v
    == local.check_rule_name]
  )
}

# Filter using checked.
locals {
  checked = data.huaweicloud_hss_baseline_security_checks_details.test.check_details[0].checked
}

data "huaweicloud_hss_baseline_security_checks_details" "checked_filter" {
  support_os = "Linux"
  standard   = "hw_standard"
  checked    = local.checked
}

output "checked_filter_is_useful" {
  value = length(data.huaweicloud_hss_baseline_security_checks_details.checked_filter.check_details) > 0 && alltrue(
    [for v in data.huaweicloud_hss_baseline_security_checks_details.checked_filter.check_details[*].checked : v == local.checked]
  )
}
`
}
