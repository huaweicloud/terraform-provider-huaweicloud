package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVulnerabilities_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_vulnerabilities.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVulnerabilities_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.vul_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.vul_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.repair_necessity"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.severity_level"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.label_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.unhandle_host_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.solution_detail"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.url"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_id_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.cve_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.repair_priority"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.hosts_num.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.hosts_num.0.common"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.repair_success_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.repair_priority_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.repair_priority_list.0.repair_priority"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.scan_time"),

					resource.TestCheckOutput("is_vul_id_filter_useful", "true"),
					resource.TestCheckOutput("is_vul_name_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckOutput("is_repair_priority_filter_useful", "true"),
					resource.TestCheckOutput("is_asset_value_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_handle_status_filter_useful", "true"),
					resource.TestCheckOutput("is_eps_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceVulnerabilities_basic = `
data "huaweicloud_hss_vulnerabilities" "test" {}

locals {
  vul_id = data.huaweicloud_hss_vulnerabilities.test.data_list[0].vul_id
}

data "huaweicloud_hss_vulnerabilities" "vul_id_filter" {
  vul_id = local.vul_id
}

output "is_vul_id_filter_useful" {
  value = length(data.huaweicloud_hss_vulnerabilities.vul_id_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_vulnerabilities.vul_id_filter.data_list[*].vul_id : v == local.vul_id]
  )
}

locals {
  vul_name = data.huaweicloud_hss_vulnerabilities.test.data_list[0].vul_name
}

data "huaweicloud_hss_vulnerabilities" "vul_name_filter" {	
  vul_name = local.vul_name
}

output "is_vul_name_filter_useful" {
  value = length(data.huaweicloud_hss_vulnerabilities.vul_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_vulnerabilities.vul_name_filter.data_list[*].vul_name : v == local.vul_name]
  )
}

locals {
  type = data.huaweicloud_hss_vulnerabilities.test.data_list[0].type
}

data "huaweicloud_hss_vulnerabilities" "type_filter" {
  type = local.type
}

output "is_type_filter_useful" {
  value = length(data.huaweicloud_hss_vulnerabilities.type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_vulnerabilities.type_filter.data_list[*].type : v == local.type]
  )
}

locals {
  repair_priority = data.huaweicloud_hss_vulnerabilities.test.data_list[0].repair_priority
}

data "huaweicloud_hss_vulnerabilities" "repair_priority_filter" {
  repair_priority = local.repair_priority
}

output "is_repair_priority_filter_useful" {
  value = length(data.huaweicloud_hss_vulnerabilities.repair_priority_filter.data_list) > 0 && alltrue(	
    [for v in data.huaweicloud_hss_vulnerabilities.repair_priority_filter.data_list[*].repair_priority : v == local.repair_priority]
  )
}

data "huaweicloud_hss_vulnerabilities" "asset_value_filter" {
  asset_value = "common"
}

output "is_asset_value_filter_useful" {
  value = length(data.huaweicloud_hss_vulnerabilities.asset_value_filter.data_list) > 0
}

data "huaweicloud_hss_vulnerabilities" "status_filter" {
  status = "vul_status_unfix"	
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_hss_vulnerabilities.status_filter.data_list) > 0
}

data "huaweicloud_hss_vulnerabilities" "handle_status_filter" {
  handle_status = "unhandled"
}

output "is_handle_status_filter_useful" {
  value = length(data.huaweicloud_hss_vulnerabilities.handle_status_filter.data_list) > 0
}

data "huaweicloud_hss_vulnerabilities" "eps_filter" {
  enterprise_project_id = "all_granted_eps"
}

output "is_eps_filter_useful" {
  value = length(data.huaweicloud_hss_vulnerabilities.eps_filter.data_list) > 0
}
`
