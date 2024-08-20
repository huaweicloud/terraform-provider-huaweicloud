package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsAssignmentPackageScores_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_rms_assignment_package_scores.basic"
	dataSource2 := "data.huaweicloud_rms_assignment_package_scores.filter_by_assignment_package_name"
	rName := acceptance.RandomAccResourceName()
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRmsAssignmentPackageScores_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_assignment_package_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRmsAssignmentPackageScores_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_assignment_package_scores" "basic" {
  depends_on = [huaweicloud_rms_assignment_package.test]
}

data "huaweicloud_rms_assignment_package_scores" "filter_by_assignment_package_name" {
  assignment_package_name = "%[2]s"

  depends_on = [huaweicloud_rms_assignment_package.test]
}

locals {
  assignment_package_name_filter_result = [for v in data.huaweicloud_rms_assignment_package_scores.filter_by_assignment_package_name.
  value[*].name : v == "%[2]s"]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_rms_assignment_package_scores.basic.value) > 0
}

output "is_assignment_package_name_filter_useful" {
  value = alltrue(local.assignment_package_name_filter_result) && length(local.assignment_package_name_filter_result) > 0
}
`, testAssignmentPackage_basic(name), name)
}
