package rms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsAssignmentPackageSummary_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rms_assignment_package_summary.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRmsAssignmentPackageSummary_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "value.#"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.compliance"),
					resource.TestCheckOutput("conformance_pack_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRmsAssignmentPackageSummary_basic() string {
	return `
data "huaweicloud_rms_assignment_package_summary" "test" {}

locals {
  conformance_pack_name = data.huaweicloud_rms_assignment_package_summary.test.value[0].name
}
data "huaweicloud_rms_assignment_package_summary" "conformance_pack_name_filter" {
  conformance_pack_name = data.huaweicloud_rms_assignment_package_summary.test.value[0].name
}
output "conformance_pack_name_filter_is_useful" {
  value = length(data.huaweicloud_rms_assignment_package_summary.conformance_pack_name_filter.value) > 0 && alltrue(
  [for v in data.huaweicloud_rms_assignment_package_summary.conformance_pack_name_filter.value[*].name :
  v == local.conformance_pack_name]
  )
}
`
}
