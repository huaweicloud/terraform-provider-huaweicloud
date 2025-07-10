package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEventSystemUserWhiteLists_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_event_system_user_white_lists.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID that has enabled host protection and added the system user
			// whitelists.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEventSystemUserWhiteLists_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "total_num"),
					resource.TestCheckResourceAttrSet(dataSource, "remain_num"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.enterprise_project_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.private_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.system_user_name_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.update_time"),

					resource.TestCheckOutput("is_host_id_filter_useful", "true"),
					resource.TestCheckOutput("is_host_name_filter_useful", "true"),
					resource.TestCheckOutput("is_private_ip_filter_useful", "true"),
					resource.TestCheckOutput("is_public_ip_filter_useful", "true"),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

const testDataSourceEventSystemUserWhiteLists_basic string = `
data "huaweicloud_hss_event_system_user_white_lists" "test" {}

# Filter using host_id.
locals {
  host_id = data.huaweicloud_hss_event_system_user_white_lists.test.data_list[0].host_id
}

data "huaweicloud_hss_event_system_user_white_lists" "host_id_filter" {
  host_id = local.host_id
}

output "is_host_id_filter_useful" {
  value = length(data.huaweicloud_hss_event_system_user_white_lists.host_id_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_event_system_user_white_lists.host_id_filter.data_list[*].host_id : v == local.host_id]
  )
}

# Filter using host_name.
locals {
  host_name = data.huaweicloud_hss_event_system_user_white_lists.test.data_list[0].host_name
}

data "huaweicloud_hss_event_system_user_white_lists" "host_name_filter" {
  host_name = local.host_name
}

output "is_host_name_filter_useful" {
  value = length(data.huaweicloud_hss_event_system_user_white_lists.host_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_event_system_user_white_lists.host_name_filter.data_list[*].host_name : v == local.host_name]
  )
}

# Filter using private_ip.
locals {
  private_ip = data.huaweicloud_hss_event_system_user_white_lists.test.data_list[0].private_ip
}

data "huaweicloud_hss_event_system_user_white_lists" "private_ip_filter" {
  private_ip = local.private_ip
}

output "is_private_ip_filter_useful" {
  value = length(data.huaweicloud_hss_event_system_user_white_lists.private_ip_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_event_system_user_white_lists.private_ip_filter.data_list[*].private_ip : v == local.private_ip]
  )
}

# Filter using public_ip.
locals {
  public_ip = data.huaweicloud_hss_event_system_user_white_lists.test.data_list[0].public_ip
}

data "huaweicloud_hss_event_system_user_white_lists" "public_ip_filter" {
  public_ip = local.public_ip
}

output "is_public_ip_filter_useful" {
  value = length(data.huaweicloud_hss_event_system_user_white_lists.public_ip_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_event_system_user_white_lists.public_ip_filter.data_list[*].public_ip : v == local.public_ip]
  )
}

# Filter using system_user_name.
locals {
  system_user_name = data.huaweicloud_hss_event_system_user_white_lists.test.data_list[0].system_user_name_list[0]
}

data "huaweicloud_hss_event_system_user_white_lists" "system_user_name_filter" {
  system_user_name = local.system_user_name
}

output "is_system_user_name_filter_useful" {
  value = length(data.huaweicloud_hss_event_system_user_white_lists.system_user_name_filter.data_list) > 0
}

# Filter using enterprise_project_id.
data "huaweicloud_hss_event_system_user_white_lists" "enterprise_project_id_filter" {
  enterprise_project_id = "all_granted_eps"
}

output "is_enterprise_project_id_filter_useful" {
  value = length(data.huaweicloud_hss_event_system_user_white_lists.enterprise_project_id_filter.data_list) > 0
}

# Filter using non existent host_name.
data "huaweicloud_hss_event_system_user_white_lists" "not_found" {
  host_name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_event_system_user_white_lists.not_found.data_list) == 0
}
`
