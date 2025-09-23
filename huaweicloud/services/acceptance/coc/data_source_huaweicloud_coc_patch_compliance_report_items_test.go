package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocPatchComplianceReportItems_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_patch_compliance_report_items.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocInstanceCompliantID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocPatchComplianceReportItems_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_items.#"),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_items.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_items.0.title"),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_items.0.compliance_level"),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_items.0.patch_detail.#"),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_items.0.patch_detail.0.installed_time"),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_items.0.patch_detail.0.patch_baseline_id"),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_items.0.patch_detail.0.patch_baseline_name"),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_items.0.patch_detail.0.patch_status"),
					resource.TestCheckOutput("title_filter_is_useful", "true"),
					resource.TestCheckOutput("sort_filter_is_useful", "true"),
					resource.TestCheckOutput("patch_status_filter_is_useful", "true"),
					resource.TestCheckOutput("classification_filter_is_useful", "true"),
					resource.TestCheckOutput("severity_level_filter_is_useful", "true"),
					resource.TestCheckOutput("compliance_level_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocPatchComplianceReportItems_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_coc_patch_compliance_report_items" "test" {
  instance_compliant_id = "%[1]s"
}

locals {
  title = [for v in data.huaweicloud_coc_patch_compliance_report_items.test.compliance_items[*].title : v if v != ""][0]
}

data "huaweicloud_coc_patch_compliance_report_items" "title_filter" {
  instance_compliant_id = "%[1]s"
  title                 = local.title
}

output "title_filter_is_useful" {
  value = length(data.huaweicloud_coc_patch_compliance_report_items.title_filter.compliance_items) > 0 && alltrue(
    [for v in data.huaweicloud_coc_patch_compliance_report_items.title_filter.compliance_items[*].title : v == local.title]
  )
}

data "huaweicloud_coc_patch_compliance_report_items" "sort_filter" {
  instance_compliant_id = "%[1]s"
  sort_dir              = "asc"
  sort_key              = "installed_time"
}

output "sort_filter_is_useful" {
  value = length(data.huaweicloud_coc_patch_compliance_report_items.sort_filter.compliance_items) > 0
}

locals {
  patch_status = [for v in data.huaweicloud_coc_patch_compliance_report_items.test.compliance_items[*].patch_detail[0].patch_status
    : v if v != ""][0]
}

data "huaweicloud_coc_patch_compliance_report_items" "patch_status_filter" {
  instance_compliant_id = "%[1]s"
  patch_status          = local.patch_status
}

output "patch_status_filter_is_useful" {
  value = length(data.huaweicloud_coc_patch_compliance_report_items.patch_status_filter.compliance_items) > 0 && alltrue(
    [for v in data.huaweicloud_coc_patch_compliance_report_items.patch_status_filter.compliance_items[*].patch_detail[0].patch_status
      : v == local.patch_status]
  )
}

locals {
  classification = [for v in data.huaweicloud_coc_patch_compliance_report_items.test.compliance_items[*].classification
    : v if v != ""][0]
}

data "huaweicloud_coc_patch_compliance_report_items" "classification_filter" {
  instance_compliant_id = "%[1]s"
  classification        = local.classification
}

output "classification_filter_is_useful" {
  value = length(data.huaweicloud_coc_patch_compliance_report_items.classification_filter.compliance_items) > 0 && alltrue(
    [for v in data.huaweicloud_coc_patch_compliance_report_items.classification_filter.compliance_items[*].classification
      : v == local.classification]
  )
}

locals {
  severity_level = [for v in data.huaweicloud_coc_patch_compliance_report_items.test.compliance_items[*].severity_level
    : v if v != ""][0]
}

data "huaweicloud_coc_patch_compliance_report_items" "severity_level_filter" {
  instance_compliant_id = "%[1]s"
  severity_level        = local.severity_level
}

output "severity_level_filter_is_useful" {
  value = length(data.huaweicloud_coc_patch_compliance_report_items.severity_level_filter.compliance_items) > 0 && alltrue(
    [for v in data.huaweicloud_coc_patch_compliance_report_items.severity_level_filter.compliance_items[*].severity_level
      : v == local.severity_level]
  )
}

locals {
  compliance_level = [for v in data.huaweicloud_coc_patch_compliance_report_items.test.compliance_items[*].compliance_level
    : v if v != ""][0]
}

data "huaweicloud_coc_patch_compliance_report_items" "compliance_level_filter" {
  instance_compliant_id = "%[1]s"
  compliance_level      = local.compliance_level
}

output "compliance_level_filter_is_useful" {
  value = length(data.huaweicloud_coc_patch_compliance_report_items.compliance_level_filter.compliance_items) > 0 && alltrue(
    [for v in data.huaweicloud_coc_patch_compliance_report_items.compliance_level_filter.compliance_items[*].compliance_level
      : v == local.compliance_level]
  )
}
`, acceptance.HW_COC_INSTANCE_COMPLIANT_ID)
}
