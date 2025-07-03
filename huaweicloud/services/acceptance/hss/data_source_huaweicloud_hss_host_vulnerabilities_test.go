package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHostVulnerabilities_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_host_vulnerabilities.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID with host protection enabled.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceHostVulnerabilities_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.vul_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.vul_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.repair_necessity"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.severity_level"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.repair_priority"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.repair_success_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.repair_cmd"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.app_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.app_list.0.app_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.app_list.0.app_version"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.app_list.0.upgrade_version"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.cve_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.cve_list.0.cve_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.cve_list.0.cvss"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.first_scan_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.scan_time"),

					resource.TestCheckOutput("is_vul_name_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_repair_priority_filter_useful", "true"),
					resource.TestCheckOutput("is_handle_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceHostVulnerabilities_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_host_vulnerabilities" "test" {
  host_id = "%[1]s"
}

locals {
  vul_name = data.huaweicloud_hss_host_vulnerabilities.test.data_list[0].vul_name
}

data "huaweicloud_hss_host_vulnerabilities" "vul_name_filter" {
  host_id  = "%[1]s"
  vul_name = local.vul_name
}

output "is_vul_name_filter_useful" {
  value = length(data.huaweicloud_hss_host_vulnerabilities.vul_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_host_vulnerabilities.vul_name_filter.data_list[*].vul_name : v == local.vul_name]
  )
}

locals {
  type = data.huaweicloud_hss_host_vulnerabilities.test.data_list[0].type
}

data "huaweicloud_hss_host_vulnerabilities" "type_filter" {
  host_id = "%[1]s"
  type    = local.type
}

output "is_type_filter_useful" {
  value = length(data.huaweicloud_hss_host_vulnerabilities.type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_host_vulnerabilities.type_filter.data_list[*].type : v == local.type]
  )
}

locals {
  status = data.huaweicloud_hss_host_vulnerabilities.test.data_list[0].status
}

data "huaweicloud_hss_host_vulnerabilities" "status_filter" {
  host_id = "%[1]s"
  status  = local.status
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_hss_host_vulnerabilities.status_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_host_vulnerabilities.status_filter.data_list[*].status : v == local.status]
  )
}

locals {
  repair_priority = data.huaweicloud_hss_host_vulnerabilities.test.data_list[0].repair_priority
}

data "huaweicloud_hss_host_vulnerabilities" "repair_priority_filter" {
  host_id         = "%[1]s"
  repair_priority = local.repair_priority
}

output "is_repair_priority_filter_useful" {
  value = length(data.huaweicloud_hss_host_vulnerabilities.repair_priority_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_host_vulnerabilities.repair_priority_filter.data_list[*].repair_priority : v == local.repair_priority]
  )
}

data "huaweicloud_hss_host_vulnerabilities" "handle_status_filter" {
  host_id       = "%[1]s"
  handle_status = "unhandled"
}

output "is_handle_status_filter_useful" {
  value = length(data.huaweicloud_hss_host_vulnerabilities.handle_status_filter.data_list) > 0
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
