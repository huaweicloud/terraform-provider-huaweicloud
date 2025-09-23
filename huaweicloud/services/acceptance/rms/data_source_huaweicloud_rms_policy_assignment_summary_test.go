package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsPolicyAssignmentSummary_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rms_policy_assignment_summary.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRmsPolicyAssignmentSummary_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_state"),
					resource.TestCheckResourceAttrSet(dataSource, "results.#"),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.resource_details.#"),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.resource_details.0.compliant_count"),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.resource_details.0.non_compliant_count"),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.assignment_details.#"),
				),
			},
		},
	})
}

func testDataSourceRmsPolicyAssignmentSummary_base(name string) string {
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
    tag_key           = "policy_filter_tag_key"
    tag_value         = "policy_filter_tag_value"
  }

  parameters = {
    listOfAllowedFlavors = "[\"${data.huaweicloud_compute_flavors.test.ids[0]}\"]"
  }

  tags = {
    foo = "bar"
  }
}
`, testAccPolicyAssignment_ecsConfig(name), name, acceptance.HW_REGION_NAME)
}

func testDataSourceRmsPolicyAssignmentSummary_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rms_policy_assignment_summary" "test" {
  policy_assignment_id = huaweicloud_rms_policy_assignment.test.id
}
`, testDataSourceRmsPolicyAssignmentSummary_base(name))
}
