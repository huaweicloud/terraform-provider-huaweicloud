package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsPolicyAssignmentEvaluateStatus_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rms_policy_assignment_evaluate_status.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRmsPolicyAssignmentEvaluateStatus_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(dataSource, "policy_assignment_id",
						"huaweicloud_rms_policy_assignment.test", "id"),
					resource.TestCheckResourceAttrSet(dataSource, "start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "state"),
				),
			},
		},
	})
}

func testDataSourceRmsPolicyAssignmentEvaluateStatus_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_policy_definitions" "test" {
  name = "allowed-ecs-flavors"
}

resource "huaweicloud_rms_policy_assignment" "test" {
  name                 = "%[2]s"
  description          = "An ECS is noncompliant if its flavor is not in the specified flavor list (filter by resource ID)."
  policy_definition_id = try(data.huaweicloud_rms_policy_definitions.test.definitions[0].id, "")

  policy_filter {
    region            = "%[3]s"
    resource_provider = "ecs"
    resource_type     = "cloudservers"
    resource_id       = huaweicloud_compute_instance.test.id
  }

  parameters = {
    listOfAllowedFlavors = "[\"${data.huaweicloud_compute_flavors.test.ids[0]}\"]"
  }
}

resource "huaweicloud_rms_policy_assignment_evaluate" "test" {
  policy_assignment_id = huaweicloud_rms_policy_assignment.test.id
}
`, testAccPolicyAssignment_ecsConfig(name), name, acceptance.HW_REGION_NAME)
}

func testDataSourceRmsPolicyAssignmentEvaluateStatus_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rms_policy_assignment_evaluate_status" "test" {
  depends_on = [huaweicloud_rms_policy_assignment_evaluate.test]

  policy_assignment_id = huaweicloud_rms_policy_assignment.test.id
}
`, testDataSourceRmsPolicyAssignmentEvaluateStatus_base(name))
}
