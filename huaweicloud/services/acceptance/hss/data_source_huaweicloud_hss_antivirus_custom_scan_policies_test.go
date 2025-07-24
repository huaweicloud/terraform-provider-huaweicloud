package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAntivirusCustomScanPolicies_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_antivirus_custom_scan_policies.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case must ensure that the host has enabled protection and created a custom scan strategy.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAntivirusCustomScanPolicies_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "total_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.policy_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.policy_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.start_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.scan_period_date"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.scan_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.scan_hour"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.scan_minute"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.next_start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.action"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.invalidate"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_info_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_info_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_info_list.0.host_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_info_list.0.private_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_info_list.0.asset_value"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.whether_paid_task"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.file_type_list.#"),

					resource.TestCheckOutput("is_policy_name_filter_useful", "true"),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

const testAccDataSourceAntivirusCustomScanPolicies_basic string = `
data "huaweicloud_hss_antivirus_custom_scan_policies" "test" {}

# Filter using policy_name.
locals {
  policy_name = data.huaweicloud_hss_antivirus_custom_scan_policies.test.data_list[0].policy_name
}

data "huaweicloud_hss_antivirus_custom_scan_policies" "policy_name_filter" {
  policy_name = local.policy_name
}

output "is_policy_name_filter_useful" {
  value = length(data.huaweicloud_hss_antivirus_custom_scan_policies.policy_name_filter.data_list) > 0
}

# Filter using enterprise_project_id.
data "huaweicloud_hss_antivirus_custom_scan_policies" "enterprise_project_id_filter" {
  enterprise_project_id = "all_granted_eps"
}

output "is_enterprise_project_id_filter_useful" {
  value = length(data.huaweicloud_hss_antivirus_custom_scan_policies.enterprise_project_id_filter.data_list) > 0
}

# Filter using non existent policy_name.
data "huaweicloud_hss_antivirus_custom_scan_policies" "not_found" {
  policy_name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_antivirus_custom_scan_policies.not_found.data_list) == 0
}
`
