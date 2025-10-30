package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAntivirusPayPerScanHosts_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_antivirus_pay_per_scan_hosts.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires the preparation of a host ID that has completed the paid virus quick scan.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAntivirusPayPerScanHosts_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.private_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.agent_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.os_type"),

					resource.TestCheckOutput("is_host_id_filter_useful", "true"),
					resource.TestCheckOutput("is_host_name_filter_useful", "true"),
					resource.TestCheckOutput("is_private_ip_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testDataSourceAntivirusPayPerScanHosts_basic() string {
	return `
data "huaweicloud_hss_antivirus_pay_per_scan_hosts" "test" {
  scan_type  = "quick"
  start_type = "now"
}

# Filter using host_id.
locals {
  host_id = data.huaweicloud_hss_antivirus_pay_per_scan_hosts.test.data_list[0].host_id
}

data "huaweicloud_hss_antivirus_pay_per_scan_hosts" "host_id_filter" {
  scan_type  = "quick"
  start_type = "now"
  host_id    = local.host_id
}

output "is_host_id_filter_useful" {
  value = length(data.huaweicloud_hss_antivirus_pay_per_scan_hosts.host_id_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_antivirus_pay_per_scan_hosts.host_id_filter.data_list[*].host_id : v == local.host_id]
  )
}

# Filter using host_name.
locals {
  host_name = data.huaweicloud_hss_antivirus_pay_per_scan_hosts.test.data_list[0].host_name
}

data "huaweicloud_hss_antivirus_pay_per_scan_hosts" "host_name_filter" {
  scan_type  = "quick"
  start_type = "now"
  host_name  = local.host_name
}

output "is_host_name_filter_useful" {
  value = length(data.huaweicloud_hss_antivirus_pay_per_scan_hosts.host_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_antivirus_pay_per_scan_hosts.host_name_filter.data_list[*].host_name : v == local.host_name]
  )
}

# Filter using private_ip.
locals {
  private_ip = data.huaweicloud_hss_antivirus_pay_per_scan_hosts.test.data_list[0].private_ip
}

data "huaweicloud_hss_antivirus_pay_per_scan_hosts" "private_ip_filter" {
  scan_type  = "quick"
  start_type = "now"
  private_ip = local.private_ip
}

output "is_private_ip_filter_useful" {
  value = length(data.huaweicloud_hss_antivirus_pay_per_scan_hosts.private_ip_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_antivirus_pay_per_scan_hosts.private_ip_filter.data_list[*].private_ip : v == local.private_ip]
  )
}

# Filter using non existent host_name.
data "huaweicloud_hss_antivirus_pay_per_scan_hosts" "not_found" {
  scan_type  = "quick"
  start_type = "now"
  host_name  = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_antivirus_pay_per_scan_hosts.not_found.data_list) == 0
}
`
}
