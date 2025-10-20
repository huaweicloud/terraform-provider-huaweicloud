package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAppWhitelistAssociateHosts_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_app_whitelist_associate_hosts.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case need to set a host ID that has enabled premium edition host protection
			// and associated with process whitelist policy.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppWhitelistAssociateHosts_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.private_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.os_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.asset_value"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.learning_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.apply_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.policy_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.policy_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.policy_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.intercept"),

					resource.TestCheckOutput("policy_name_filter_useful", "true"),
					resource.TestCheckOutput("learning_status_filter_useful", "true"),
					resource.TestCheckOutput("asset_value_filter_useful", "true"),
					resource.TestCheckOutput("host_name_filter_useful", "true"),
					resource.TestCheckOutput("private_ip_filter_useful", "true"),
					resource.TestCheckOutput("os_type_filter_useful", "true"),
					resource.TestCheckOutput("policy_id_filter_useful", "true"),
					resource.TestCheckOutput("eps_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceAppWhitelistAssociateHosts_basic = `
data "huaweicloud_hss_app_whitelist_associate_hosts" "test" {}

locals {
  policy_name = data.huaweicloud_hss_app_whitelist_associate_hosts.test.data_list[0].policy_name
}

data "huaweicloud_hss_app_whitelist_associate_hosts" "policy_name_filter" {
  policy_name = local.policy_name
}

output "policy_name_filter_useful" {
  value = length(data.huaweicloud_hss_app_whitelist_associate_hosts.policy_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_app_whitelist_associate_hosts.policy_name_filter.data_list[*].policy_name : v == local.policy_name]
  )
}

locals {
  learning_status = data.huaweicloud_hss_app_whitelist_associate_hosts.test.data_list[0].learning_status
}

data "huaweicloud_hss_app_whitelist_associate_hosts" "learning_status_filter" {
  learning_status = local.learning_status
}

output "learning_status_filter_useful" {
  value = length(data.huaweicloud_hss_app_whitelist_associate_hosts.learning_status_filter.data_list) > 0 && alltrue(	
    [for v in data.huaweicloud_hss_app_whitelist_associate_hosts.learning_status_filter.data_list[*].learning_status : v == local.learning_status]
  )
}

locals {
  asset_value = data.huaweicloud_hss_app_whitelist_associate_hosts.test.data_list[0].asset_value
}

data "huaweicloud_hss_app_whitelist_associate_hosts" "asset_value_filter" {	
  asset_value = local.asset_value
}

output "asset_value_filter_useful" {
  value = length(data.huaweicloud_hss_app_whitelist_associate_hosts.asset_value_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_app_whitelist_associate_hosts.asset_value_filter.data_list[*].asset_value : v == local.asset_value]
  )
}

locals {
  host_name = data.huaweicloud_hss_app_whitelist_associate_hosts.test.data_list[0].host_name
}

data "huaweicloud_hss_app_whitelist_associate_hosts" "host_name_filter" {	
  host_name = local.host_name
}

output "host_name_filter_useful" {
  value = length(data.huaweicloud_hss_app_whitelist_associate_hosts.host_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_app_whitelist_associate_hosts.host_name_filter.data_list[*].host_name : v == local.host_name]
  )
}

locals {
  private_ip = data.huaweicloud_hss_app_whitelist_associate_hosts.test.data_list[0].private_ip
}

data "huaweicloud_hss_app_whitelist_associate_hosts" "private_ip_filter" {
  private_ip = local.private_ip
}

output "private_ip_filter_useful" {
  value = length(data.huaweicloud_hss_app_whitelist_associate_hosts.private_ip_filter.data_list) > 0 && alltrue(	
    [for v in data.huaweicloud_hss_app_whitelist_associate_hosts.private_ip_filter.data_list[*].private_ip : v == local.private_ip]
  )
}

locals {
  os_type = data.huaweicloud_hss_app_whitelist_associate_hosts.test.data_list[0].os_type
}

data "huaweicloud_hss_app_whitelist_associate_hosts" "os_type_filter" {
  os_type = local.os_type
}

output "os_type_filter_useful" {
  value = length(data.huaweicloud_hss_app_whitelist_associate_hosts.os_type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_app_whitelist_associate_hosts.os_type_filter.data_list[*].os_type : v == local.os_type]
  )
}

locals {
  policy_id = data.huaweicloud_hss_app_whitelist_associate_hosts.test.data_list[0].policy_id
}

data "huaweicloud_hss_app_whitelist_associate_hosts" "policy_id_filter" {	
  policy_id = local.policy_id
}

output "policy_id_filter_useful" {
  value = length(data.huaweicloud_hss_app_whitelist_associate_hosts.policy_id_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_app_whitelist_associate_hosts.policy_id_filter.data_list[*].policy_id : v == local.policy_id]
  )
}

data "huaweicloud_hss_app_whitelist_associate_hosts" "eps_filter" {
  enterprise_project_id = "all_granted_eps"
}

output "eps_filter_useful" {
  value = length(data.huaweicloud_hss_app_whitelist_associate_hosts.eps_filter.data_list) > 0
}
`
