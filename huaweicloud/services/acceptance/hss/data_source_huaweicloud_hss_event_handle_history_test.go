package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEventHandleHistory_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_event_handle_history.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running the test, you need to prepare a host with premium edition host protection enabled.
			// And manually handle a event alarm.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEventHandleHistory_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.event_abstract"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.attack_tag"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.private_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.asset_value"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.occur_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.handle_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.notes"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.event_class_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.event_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.handle_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.operate_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.severity"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.user_name"),

					resource.TestCheckOutput("is_severity_filter_useful", "true"),
					resource.TestCheckOutput("is_attack_tag_filter_useful", "true"),
					resource.TestCheckOutput("is_asset_value_filter_useful", "true"),
					resource.TestCheckOutput("is_event_name_filter_useful", "true"),
					resource.TestCheckOutput("is_host_name_filter_useful", "true"),
					resource.TestCheckOutput("is_handle_status_filter_useful", "true"),
					resource.TestCheckOutput("is_private_ip_filter_useful", "true"),
					resource.TestCheckOutput("is_sort_filter_useful", "true"),
					resource.TestCheckOutput("is_eps_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceEventHandleHistory_basic = `
data "huaweicloud_hss_event_handle_history" "test" {}

locals {
  severity = data.huaweicloud_hss_event_handle_history.test.data_list[0].severity
}

data "huaweicloud_hss_event_handle_history" "severity_filter" {
  severity = local.severity
}

output "is_severity_filter_useful" {
  value = length(data.huaweicloud_hss_event_handle_history.severity_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_event_handle_history.severity_filter.data_list[*].severity : v == local.severity]
  )
}

locals {
  attack_tag = data.huaweicloud_hss_event_handle_history.test.data_list[0].attack_tag
}

data "huaweicloud_hss_event_handle_history" "attack_tag_filter" {	
  attack_tag = local.attack_tag
}

output "is_attack_tag_filter_useful" {
  value = length(data.huaweicloud_hss_event_handle_history.attack_tag_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_event_handle_history.attack_tag_filter.data_list[*].attack_tag : v == local.attack_tag]
  )
}

locals {
  asset_value = data.huaweicloud_hss_event_handle_history.test.data_list[0].asset_value
}

data "huaweicloud_hss_event_handle_history" "asset_value_filter" {
  asset_value = local.asset_value
}

output "is_asset_value_filter_useful" {
  value = length(data.huaweicloud_hss_event_handle_history.asset_value_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_event_handle_history.asset_value_filter.data_list[*].asset_value : v == local.asset_value]
  )
}

locals {
  event_name = data.huaweicloud_hss_event_handle_history.test.data_list[0].event_name
}

data "huaweicloud_hss_event_handle_history" "event_name_filter" {
  event_name = local.event_name
}

output "is_event_name_filter_useful" {
  value = length(data.huaweicloud_hss_event_handle_history.event_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_event_handle_history.event_name_filter.data_list[*].event_name : v == local.event_name]
  )
}

locals {
  host_name = data.huaweicloud_hss_event_handle_history.test.data_list[0].host_name
}

data "huaweicloud_hss_event_handle_history" "host_name_filter" {
  host_name = local.host_name
}

output "is_host_name_filter_useful" {
  value = length(data.huaweicloud_hss_event_handle_history.host_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_event_handle_history.host_name_filter.data_list[*].host_name : v == local.host_name]
  )
}

locals {
  handle_status = data.huaweicloud_hss_event_handle_history.test.data_list[0].handle_status
}	

data "huaweicloud_hss_event_handle_history" "handle_status_filter" {
  handle_status = local.handle_status
}

output "is_handle_status_filter_useful" {
  value = length(data.huaweicloud_hss_event_handle_history.handle_status_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_event_handle_history.handle_status_filter.data_list[*].handle_status : v == local.handle_status]
  )
}

locals {
  private_ip = data.huaweicloud_hss_event_handle_history.test.data_list[0].private_ip
}

data "huaweicloud_hss_event_handle_history" "private_ip_filter" {
  private_ip = local.private_ip
}

output "is_private_ip_filter_useful" {
  value = length(data.huaweicloud_hss_event_handle_history.private_ip_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_event_handle_history.private_ip_filter.data_list[*].private_ip : v == local.private_ip]
  )
}

data "huaweicloud_hss_event_handle_history" "sort_filter" {
  sort_key = "handle_time"
  sort_dir = "asc"
}

output "is_sort_filter_useful" {
  value = length(data.huaweicloud_hss_event_handle_history.sort_filter.data_list) > 0 
}

data "huaweicloud_hss_event_handle_history" "eps_filter" {
  enterprise_project_id = "all_granted_eps"
}

output "is_eps_filter_useful" {
  value = length(data.huaweicloud_hss_event_handle_history.eps_filter.data_list) > 0
}
`
