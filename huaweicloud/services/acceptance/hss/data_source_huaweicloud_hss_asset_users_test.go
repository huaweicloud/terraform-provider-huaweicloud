package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssetUsers_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_asset_users.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAssetUsers_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.agent_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.login_permission"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.root_permission"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.user_group_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.user_home_dir"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.shell"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.recent_scan_time"),

					resource.TestCheckOutput("is_host_id_filter_useful", "true"),
					resource.TestCheckOutput("is_user_name_filter_useful", "true"),
					resource.TestCheckOutput("is_host_name_filter_useful", "true"),
					resource.TestCheckOutput("is_private_ip_filter_useful", "true"),
					resource.TestCheckOutput("is_login_permission_filter_useful", "true"),
					resource.TestCheckOutput("is_root_permission_filter_useful", "true"),
					resource.TestCheckOutput("is_user_group_filter_useful", "true"),
					resource.TestCheckOutput("is_category_filter_useful", "true"),
					resource.TestCheckOutput("is_part_match_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

const testDataSourceAssetUsers_basic string = `
data "huaweicloud_hss_asset_users" "test" {}

# Filter using host id.
locals {
  host_id = data.huaweicloud_hss_asset_users.test.data_list[0].host_id
}

data "huaweicloud_hss_asset_users" "host_id_filter" {
  host_id = local.host_id
}

output "is_host_id_filter_useful" {
  value = length(data.huaweicloud_hss_asset_users.host_id_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_users.host_id_filter.data_list[*].host_id : v == local.host_id]
  )
}

# Filter using user name.
locals {
  user_name = data.huaweicloud_hss_asset_users.test.data_list[0].user_name
}

data "huaweicloud_hss_asset_users" "user_name_filter" {
  user_name = local.user_name
}

output "is_user_name_filter_useful" {
  value = length(data.huaweicloud_hss_asset_users.user_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_users.user_name_filter.data_list[*].user_name : v == local.user_name]
  )
}

# Filter using host name.
locals {
  host_name = data.huaweicloud_hss_asset_users.test.data_list[0].host_name
}

data "huaweicloud_hss_asset_users" "host_name_filter" {
  host_name = local.host_name
}

output "is_host_name_filter_useful" {
  value = length(data.huaweicloud_hss_asset_users.host_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_users.host_name_filter.data_list[*].host_name : v == local.host_name]
  )
}

# Filter using private IP.
locals {
  private_ip = data.huaweicloud_hss_asset_users.test.data_list[0].host_ip
}

data "huaweicloud_hss_asset_users" "private_ip_filter" {
  private_ip = local.private_ip
}

output "is_private_ip_filter_useful" {
  value = length(data.huaweicloud_hss_asset_users.private_ip_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_users.private_ip_filter.data_list[*].host_ip : v == local.private_ip]
  )
}

# Filter using login permission.
data "huaweicloud_hss_asset_users" "login_permission_filter" {
  login_permission = true
}

output "is_login_permission_filter_useful" {
  value = length(data.huaweicloud_hss_asset_users.login_permission_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_users.login_permission_filter.data_list[*].login_permission : v == true]
  )
}

# Filter using root permission.
data "huaweicloud_hss_asset_users" "root_permission_filter" {
  root_permission = true
}

output "is_root_permission_filter_useful" {
  value = length(data.huaweicloud_hss_asset_users.root_permission_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_users.root_permission_filter.data_list[*].root_permission : v == true]
  )
}

# Filter using user group.
locals {
  user_group = data.huaweicloud_hss_asset_users.test.data_list[0].user_group_name
}

data "huaweicloud_hss_asset_users" "user_group_filter" {
  user_group = local.user_group
}

output "is_user_group_filter_useful" {
  value = length(data.huaweicloud_hss_asset_users.user_group_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_users.user_group_filter.data_list[*].user_group_name : v == local.user_group]
  )
}

# Filter using category.
data "huaweicloud_hss_asset_users" "category_filter" {
  category = "host"
}

output "is_category_filter_useful" {
  value = length(data.huaweicloud_hss_asset_users.category_filter.data_list) > 0
}

# Filter using part match.
data "huaweicloud_hss_asset_users" "part_match_filter" {
  user_name  = substr(data.huaweicloud_hss_asset_users.test.data_list[0].user_name, 0, 3)
  part_match = true
}

output "is_part_match_filter_useful" {
  value = length(data.huaweicloud_hss_asset_users.part_match_filter.data_list) > 0
}

# Filter using non-existent user name.
data "huaweicloud_hss_asset_users" "not_found" {
  user_name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_asset_users.not_found.data_list) == 0
}
`
