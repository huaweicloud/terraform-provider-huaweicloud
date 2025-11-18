package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceImageBaselineRiskConfigs_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_image_baseline_risk_configs.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceImageBaselineRiskConfigs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.severity"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.check_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.check_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.standard"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.check_rule_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.failed_rule_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.check_type_desc"),

					resource.TestCheckOutput("is_check_name_filter_useful", "true"),
					resource.TestCheckOutput("is_severity_filter_useful", "true"),
					resource.TestCheckOutput("is_standard_filter_useful", "true"),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceImageBaselineRiskConfigs_basic() string {
	return `
data "huaweicloud_hss_image_baseline_risk_configs" "test" {
  image_type = "private_image"
}

# Filter using check_name.
locals {
  check_name = data.huaweicloud_hss_image_baseline_risk_configs.test.data_list[0].check_name
}

data "huaweicloud_hss_image_baseline_risk_configs" "check_name_filter" {
  image_type = "private_image"
  check_name = local.check_name
}

output "is_check_name_filter_useful" {
  value = length(data.huaweicloud_hss_image_baseline_risk_configs.check_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_baseline_risk_configs.check_name_filter.data_list[*].check_name : v == local.check_name]
  )
}

# Filter using severity.
locals {
  severity = data.huaweicloud_hss_image_baseline_risk_configs.test.data_list[0].severity
}

data "huaweicloud_hss_image_baseline_risk_configs" "severity_filter" {
  image_type = "private_image"
  severity   = local.severity
}

output "is_severity_filter_useful" {
  value = length(data.huaweicloud_hss_image_baseline_risk_configs.severity_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_baseline_risk_configs.severity_filter.data_list[*].severity : v == local.severity]
  )
}

# Filter using standard.
locals {
  standard = data.huaweicloud_hss_image_baseline_risk_configs.test.data_list[0].standard
}

data "huaweicloud_hss_image_baseline_risk_configs" "standard_filter" {
  image_type = "private_image"
  standard   = local.standard
}

output "is_standard_filter_useful" {
  value = length(data.huaweicloud_hss_image_baseline_risk_configs.standard_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_baseline_risk_configs.standard_filter.data_list[*].standard : v == local.standard]
  )
}

# Filter using enterprise_project_id.
data "huaweicloud_hss_image_baseline_risk_configs" "enterprise_project_id_filter" {
  image_type            = "private_image"
  enterprise_project_id = "0"
}

output "is_enterprise_project_id_filter_useful" {
  value = length(data.huaweicloud_hss_image_baseline_risk_configs.enterprise_project_id_filter.data_list) > 0
}
`
}
