package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsAssignmentPackageResults_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_rms_assignment_package_results.basic"
	dataSource2 := "data.huaweicloud_rms_assignment_package_results.filter_by_policy_assignment_name"
	rName := acceptance.RandomAccResourceName()
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"null": {
				Source:            "hashicorp/null",
				VersionConstraint: "3.2.1",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRmsAssignmentPackageResults_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_policy_assignment_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRmsAssignmentPackageResults_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_assignment_package_results" "basic" {
  assignment_package_id = huaweicloud_rms_assignment_package.test.id
}

# wait 30 seconds to let the assignment package evaluate
resource "null_resource" "test" {
  provisioner "local-exec" {
    command = "sleep 30"
  }

  depends_on = [ huaweicloud_rms_assignment_package.test ]
}

data "huaweicloud_rms_assignment_package_results" "filter_by_policy_assignment_name" {
  assignment_package_id  = huaweicloud_rms_assignment_package.test.id
  policy_assignment_name = "ecs-instance-no-public-ip"

  depends_on = [ null_resource.test ]
}

locals {
  policy_assignment_name_filter_result = [for v in data.huaweicloud_rms_assignment_package_results.filter_by_policy_assignment_name.
  value[*].policy_assignment_name : strcontains(v, "ecs-instance-no-public-ip")]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_rms_assignment_package_results.basic.value) > 0
}

output "is_policy_assignment_name_filter_useful" {
  value = alltrue(local.policy_assignment_name_filter_result) && length(local.policy_assignment_name_filter_result) > 0
}
`, testAssignmentPackage_basic(name), name)
}
