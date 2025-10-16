package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAppWhitelistPolicies_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_app_whitelist_policies.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case need to set a host ID that has enabled premium edition host protection
			// and create a process whitelist policy.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppWhitelistPolicies_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.policy_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.policy_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.policy_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.learning_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.learning_days"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.specified_dir"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.intercept"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.auto_detect"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.auto_confirm"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.default_policy"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_id_list.#"),

					resource.TestCheckOutput("policy_name_filter_useful", "true"),
					resource.TestCheckOutput("policy_type_filter_useful", "true"),
					resource.TestCheckOutput("learning_status_filter_useful", "true"),
					resource.TestCheckOutput("eps_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceAppWhitelistPolicies_basic = `
data "huaweicloud_hss_app_whitelist_policies" "test" {}

locals {
  policy_name = data.huaweicloud_hss_app_whitelist_policies.test.data_list[0].policy_name
}

data "huaweicloud_hss_app_whitelist_policies" "policy_name_filter" {
  policy_name = local.policy_name
}

output "policy_name_filter_useful" {
  value = length(data.huaweicloud_hss_app_whitelist_policies.policy_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_app_whitelist_policies.policy_name_filter.data_list[*].policy_name : v == local.policy_name]
  )
}

locals {
  policy_type = data.huaweicloud_hss_app_whitelist_policies.test.data_list[0].policy_type
}

data "huaweicloud_hss_app_whitelist_policies" "policy_type_filter" {	
  policy_type = local.policy_type
}

output "policy_type_filter_useful" {
  value = length(data.huaweicloud_hss_app_whitelist_policies.policy_type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_app_whitelist_policies.policy_type_filter.data_list[*].policy_type : v == local.policy_type]
  )
}

locals {
  learning_status = data.huaweicloud_hss_app_whitelist_policies.test.data_list[0].learning_status
}

data "huaweicloud_hss_app_whitelist_policies" "learning_status_filter" {
  learning_status = local.learning_status
}

output "learning_status_filter_useful" {
  value = length(data.huaweicloud_hss_app_whitelist_policies.learning_status_filter.data_list) > 0 && alltrue(	
    [for v in data.huaweicloud_hss_app_whitelist_policies.learning_status_filter.data_list[*].learning_status : v == local.learning_status]
  )
}

data "huaweicloud_hss_app_whitelist_policies" "eps_filter" {
  enterprise_project_id = "all_granted_eps"
}

output "eps_filter_useful" {
  value = length(data.huaweicloud_hss_app_whitelist_policies.eps_filter.data_list) > 0
}
`
