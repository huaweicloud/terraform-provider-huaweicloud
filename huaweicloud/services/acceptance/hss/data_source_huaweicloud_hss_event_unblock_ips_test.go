package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEventUnblockIps_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_event_unblock_ips.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID with host protection enabled,
			// and the host is under the default enterprise project.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEventUnblockIps_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.src_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.login_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.intercept_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.intercept_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.block_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.latest_time"),

					resource.TestCheckOutput("is_last_days_filter_useful", "true"),
					resource.TestCheckOutput("is_host_name_filter_useful", "true"),
					resource.TestCheckOutput("is_src_ip_filter_useful", "true"),
					resource.TestCheckOutput("is_intercept_status_filter_useful", "true"),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testDataSourceEventUnblockIps_base() string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_event_unblock_ip" "test" {
  data_list {
    host_id    = "%[1]s"
    src_ip     = "127.0.0.1"
    login_type = "mysql"
  }

  data_list {
    host_id    = "%[1]s"
    src_ip     = "127.0.0.2"
    login_type = "mysql"
  }
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}

func testDataSourceEventUnblockIps_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_hss_event_unblock_ips" "test" {
  depends_on = [huaweicloud_hss_event_unblock_ip.test]
}

# Filter using last_days.
data "huaweicloud_hss_event_unblock_ips" "last_days_filter" {
  depends_on = [huaweicloud_hss_event_unblock_ip.test]

  last_days = 1
}

output "is_last_days_filter_useful" {
  value = length(data.huaweicloud_hss_event_unblock_ips.last_days_filter.data_list) > 0
}

# Filter using host_name.
locals {
  host_name = data.huaweicloud_hss_event_unblock_ips.test.data_list[0].host_name
}

data "huaweicloud_hss_event_unblock_ips" "host_name_filter" {
  host_name = local.host_name
}

output "is_host_name_filter_useful" {
  value = length(data.huaweicloud_hss_event_unblock_ips.host_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_event_unblock_ips.host_name_filter.data_list[*].host_name : v == local.host_name]
  )
}

# Filter using src_ip.
locals {
  src_ip = data.huaweicloud_hss_event_unblock_ips.test.data_list[0].src_ip
}

data "huaweicloud_hss_event_unblock_ips" "src_ip_filter" {
  src_ip = local.src_ip
}

output "is_src_ip_filter_useful" {
  value = length(data.huaweicloud_hss_event_unblock_ips.src_ip_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_event_unblock_ips.src_ip_filter.data_list[*].src_ip : v == local.src_ip]
  )
}

# Filter using intercept_status.
locals {
  intercept_status = data.huaweicloud_hss_event_unblock_ips.test.data_list[0].intercept_status
}

data "huaweicloud_hss_event_unblock_ips" "intercept_status_filter" {
  intercept_status = local.intercept_status
}

output "is_intercept_status_filter_useful" {
  value = length(data.huaweicloud_hss_event_unblock_ips.intercept_status_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_event_unblock_ips.intercept_status_filter.data_list[*].intercept_status : v == local.intercept_status]
  )
}

# Filter using enterprise_project_id.
data "huaweicloud_hss_event_unblock_ips" "enterprise_project_id_filter" {
  depends_on = [huaweicloud_hss_event_unblock_ip.test]

  enterprise_project_id = "0"
}

output "is_enterprise_project_id_filter_useful" {
  value = length(data.huaweicloud_hss_event_unblock_ips.enterprise_project_id_filter.data_list) > 0
}

# Filter using non existent enterprise_project_id.
data "huaweicloud_hss_event_unblock_ips" "not_found" {
  depends_on = [huaweicloud_hss_event_unblock_ip.test]

  enterprise_project_id = "1"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_event_unblock_ips.not_found.data_list) == 0
}
`, testDataSourceEventUnblockIps_base())
}
