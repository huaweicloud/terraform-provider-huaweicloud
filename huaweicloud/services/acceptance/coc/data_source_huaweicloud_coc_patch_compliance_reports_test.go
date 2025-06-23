package coc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocPatchComplianceReports_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_patch_compliance_reports.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocInstanceCompliantID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocPatchComplianceReports_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.compliant_summary.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.compliant_summary.0.compliant_count"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.compliant_summary.0.severity_summary.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.compliant_summary.0.severity_summary.0.critical_count"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.compliant_summary.0.severity_summary.0.high_count"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.compliant_summary.0.severity_summary.0.informational_count"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.compliant_summary.0.severity_summary.0.low_count"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.compliant_summary.0.severity_summary.0.medium_count"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.compliant_summary.0.severity_summary.0.unspecified_count"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.non_compliant_summary.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.non_compliant_summary.0.non_compliant_count"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.non_compliant_summary.0.severity_summary.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.non_compliant_summary.0.severity_summary.0.critical_count"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.non_compliant_summary.0.severity_summary.0.high_count"),
					resource.TestCheckResourceAttrSet(dataSource,
						"instance_compliant.0.non_compliant_summary.0.severity_summary.0.informational_count"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.non_compliant_summary.0.severity_summary.0.low_count"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.non_compliant_summary.0.severity_summary.0.medium_count"),
					resource.TestCheckResourceAttrSet(dataSource,
						"instance_compliant.0.non_compliant_summary.0.severity_summary.0.unspecified_count"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.execution_summary.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.execution_summary.0.order_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.execution_summary.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.execution_summary.0.report_time"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.ip"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.eip"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.region"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.report_scene"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.baseline_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.baseline_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.rule_type"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_compliant.0.operating_system"),
					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("ip_filter_is_useful", "true"),
					resource.TestCheckOutput("eip_filter_is_useful", "true"),
					resource.TestCheckOutput("operating_system_filter_is_useful", "true"),
					resource.TestCheckOutput("region_filter_is_useful", "true"),
					resource.TestCheckOutput("compliant_status_filter_is_useful", "true"),
					resource.TestCheckOutput("order_id_filter_is_useful", "true"),
					resource.TestCheckOutput("sort_filter_is_useful", "true"),
					resource.TestCheckOutput("report_scene_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocPatchComplianceReports_basic() string {
	return `
data "huaweicloud_coc_patch_compliance_reports" "test" {}

locals {
  enterprise_project_id = [for v in data.huaweicloud_coc_patch_compliance_reports.test.instance_compliant[*].enterprise_project_id
    : v if v != ""][0]
}

data "huaweicloud_coc_patch_compliance_reports" "enterprise_project_id_filter" {
  enterprise_project_id = local.enterprise_project_id
}

output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_patch_compliance_reports.enterprise_project_id_filter.instance_compliant) > 0 && alltrue(
    [for v in data.huaweicloud_coc_patch_compliance_reports.enterprise_project_id_filter.instance_compliant[*].enterprise_project_id
      : v == local.enterprise_project_id]
  )
}

locals {
  name = [for v in data.huaweicloud_coc_patch_compliance_reports.test.instance_compliant[*].name : v if v != ""][0]
}

data "huaweicloud_coc_patch_compliance_reports" "name_filter" {
  name = local.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_coc_patch_compliance_reports.name_filter.instance_compliant) > 0 && alltrue(
    [for v in data.huaweicloud_coc_patch_compliance_reports.name_filter.instance_compliant[*].name : v == local.name]
  )
}

locals {
  instance_id = [for v in data.huaweicloud_coc_patch_compliance_reports.test.instance_compliant[*].instance_id : v if v != ""][0]
}

data "huaweicloud_coc_patch_compliance_reports" "instance_id_filter" {
  instance_id = local.instance_id
}

output "instance_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_patch_compliance_reports.instance_id_filter.instance_compliant) > 0 && alltrue(
    [for v in data.huaweicloud_coc_patch_compliance_reports.instance_id_filter.instance_compliant[*].instance_id : v == local.instance_id]
  )
}

locals {
  ip = [for v in data.huaweicloud_coc_patch_compliance_reports.test.instance_compliant[*].ip : v if v != ""][0]
}

data "huaweicloud_coc_patch_compliance_reports" "ip_filter" {
  ip = local.ip
}

output "ip_filter_is_useful" {
  value = length(data.huaweicloud_coc_patch_compliance_reports.ip_filter.instance_compliant) > 0 && alltrue(
    [for v in data.huaweicloud_coc_patch_compliance_reports.ip_filter.instance_compliant[*].ip : v == local.ip]
  )
}

locals {
  eip = [for v in data.huaweicloud_coc_patch_compliance_reports.test.instance_compliant[*].eip : v if v != ""][0]
}

data "huaweicloud_coc_patch_compliance_reports" "eip_filter" {
  eip = local.eip
}

output "eip_filter_is_useful" {
  value = length(data.huaweicloud_coc_patch_compliance_reports.eip_filter.instance_compliant) > 0 && alltrue(
    [for v in data.huaweicloud_coc_patch_compliance_reports.eip_filter.instance_compliant[*].eip : v == local.eip]
  )
}

locals {
  operating_system = [for v in data.huaweicloud_coc_patch_compliance_reports.test.instance_compliant[*].operating_system : v if v != ""][0]
}

data "huaweicloud_coc_patch_compliance_reports" "operating_system_filter" {
  operating_system = local.operating_system
}

output "operating_system_filter_is_useful" {
  value = length(data.huaweicloud_coc_patch_compliance_reports.operating_system_filter.instance_compliant) > 0 && alltrue(
    [for v in data.huaweicloud_coc_patch_compliance_reports.operating_system_filter.instance_compliant[*].operating_system :
      v == local.operating_system]
  )
}

locals {
  region = [for v in data.huaweicloud_coc_patch_compliance_reports.test.instance_compliant[*].region : v if v != ""][0]
}

data "huaweicloud_coc_patch_compliance_reports" "region_filter" {
  region = local.region
}

output "region_filter_is_useful" {
  value = length(data.huaweicloud_coc_patch_compliance_reports.region_filter.instance_compliant) > 0 && alltrue(
    [for v in data.huaweicloud_coc_patch_compliance_reports.region_filter.instance_compliant[*].region : v == local.region]
  )
}

locals {
  compliant_status = [for v in data.huaweicloud_coc_patch_compliance_reports.test.instance_compliant[*].status : v if v != ""][0]
}

data "huaweicloud_coc_patch_compliance_reports" "compliant_status_filter" {
  compliant_status = local.compliant_status
}

output "compliant_status_filter_is_useful" {
  value = length(data.huaweicloud_coc_patch_compliance_reports.compliant_status_filter.instance_compliant) > 0 && alltrue(
    [for v in data.huaweicloud_coc_patch_compliance_reports.compliant_status_filter.instance_compliant[*].status : v == local.compliant_status]
  )
}

locals {
  order_id = [for v in data.huaweicloud_coc_patch_compliance_reports.test.instance_compliant[*].execution_summary[0].order_id : v if v != ""][0]
}

data "huaweicloud_coc_patch_compliance_reports" "order_id_filter" {
  order_id = local.order_id
}

output "order_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_patch_compliance_reports.order_id_filter.instance_compliant) > 0 && alltrue(
    [for v in data.huaweicloud_coc_patch_compliance_reports.order_id_filter.instance_compliant[*].execution_summary[0].order_id :
      v == local.order_id]
  )
}

data "huaweicloud_coc_patch_compliance_reports" "sort_filter" {
  sort_dir = "asc"
  sort_key = "report_time"
}

output "sort_filter_is_useful" {
  value = length(data.huaweicloud_coc_patch_compliance_reports.sort_filter.instance_compliant) > 0
}

locals {
  report_scene = [for v in data.huaweicloud_coc_patch_compliance_reports.test.instance_compliant[*].report_scene : v if v != ""][0]
}

data "huaweicloud_coc_patch_compliance_reports" "report_scene_filter" {
  report_scene = local.report_scene
}

output "report_scene_filter_is_useful" {
  value = length(data.huaweicloud_coc_patch_compliance_reports.report_scene_filter.instance_compliant) > 0 && alltrue(
    [for v in data.huaweicloud_coc_patch_compliance_reports.report_scene_filter.instance_compliant[*].report_scene : v == local.report_scene]
  )
}
`
}
