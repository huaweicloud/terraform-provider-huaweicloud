package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAlarmWhiteLists_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_event_alarm_white_lists.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case must ensure that the host has enabled protection and add alarm events to the whitelist.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAlarmWhiteLists_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "total_num"),
					resource.TestCheckResourceAttrSet(dataSource, "remain_num"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_num"),
					resource.TestCheckResourceAttrSet(dataSource, "event_type_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.enterprise_project_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.hash"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.event_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.white_field"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.field_value"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.judge_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.update_time"),

					resource.TestCheckOutput("is_hash_filter_useful", "true"),
					resource.TestCheckOutput("is_event_type_filter_useful", "true"),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

const testAccDataSourceAlarmWhiteLists_basic string = `
data "huaweicloud_hss_event_alarm_white_lists" "test" {}

# Filter using hash.
locals {
  hash = data.huaweicloud_hss_event_alarm_white_lists.test.data_list[0].hash
}

data "huaweicloud_hss_event_alarm_white_lists" "hash_filter" {
  hash = local.hash
}

output "is_hash_filter_useful" {
  value = length(data.huaweicloud_hss_event_alarm_white_lists.hash_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_event_alarm_white_lists.hash_filter.data_list[*].hash : v == local.hash]
  )
}

# Filter using event_type.
locals {
  event_type = data.huaweicloud_hss_event_alarm_white_lists.test.data_list[0].event_type
}

data "huaweicloud_hss_event_alarm_white_lists" "event_type_filter" {
  event_type = local.event_type
}

output "is_event_type_filter_useful" {
  value = length(data.huaweicloud_hss_event_alarm_white_lists.event_type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_event_alarm_white_lists.event_type_filter.data_list[*].event_type : v == local.event_type]
  )
}

# Filter using enterprise_project_id.
data "huaweicloud_hss_event_alarm_white_lists" "enterprise_project_id_filter" {
  enterprise_project_id = "all_granted_eps"
}

output "is_enterprise_project_id_filter_useful" {
  value = length(data.huaweicloud_hss_event_alarm_white_lists.enterprise_project_id_filter.data_list) > 0
}

# Filter using non existent event_type.
data "huaweicloud_hss_event_alarm_white_lists" "not_found" {
  event_type = 22222
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_event_alarm_white_lists.not_found.data_list) == 0
}
`
