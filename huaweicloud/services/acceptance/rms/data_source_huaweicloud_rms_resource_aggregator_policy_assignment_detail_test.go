package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAggregatorPolicyAssignmentDetail_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rms_resource_aggregator_policy_assignment_detail.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAggregatorPolicyAssignmentDetail_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "policy_assignment_type"),
					resource.TestCheckResourceAttrSet(dataSource, "name"),
					resource.TestCheckResourceAttrSet(dataSource, "description"),
					resource.TestCheckResourceAttrSet(dataSource, "policy_filter.#"),
					resource.TestCheckResourceAttrSet(dataSource, "policy_filter_v2.#"),
					resource.TestCheckResourceAttrSet(dataSource, "state"),
					resource.TestCheckResourceAttrSet(dataSource, "created"),
					resource.TestCheckResourceAttrSet(dataSource, "updated"),
					resource.TestCheckResourceAttrSet(dataSource, "custom_policy.#"),
					resource.TestCheckResourceAttrSet(dataSource, "custom_policy.0.function_urn"),
					resource.TestCheckResourceAttrSet(dataSource, "custom_policy.0.auth_type"),
					resource.TestCheckResourceAttrSet(dataSource, "custom_policy.0.auth_value.%"),
					resource.TestCheckResourceAttrSet(dataSource, "parameters.%"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "created_by"),
				),
			},
		},
	})
}

func testDataSourceAggregatorPolicyAssignmentDetail_base(name, accountId string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rms_resource_aggregator" "test" {
  name        = "%[1]s"
  type        = "ACCOUNT"
  account_ids = ["%[2]s"]
}

data "huaweicloud_rms_resource_aggregator_policy_assignments" "test" {
  aggregator_id = huaweicloud_rms_resource_aggregator.test.id
}
`, name, accountId)
}

func testDataSourceAggregatorPolicyAssignmentDetail_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_resource_aggregator_policy_assignment_detail" "test" {
  aggregator_id        = huaweicloud_rms_resource_aggregator.test.id
  account_id           = "%[2]s"
  policy_assignment_id = data.huaweicloud_rms_resource_aggregator_policy_assignments.test.assignments[0].policy_assignment_id
}
`, testDataSourceAggregatorPolicyAssignmentDetail_base(name, acceptance.HW_DOMAIN_ID), acceptance.HW_DOMAIN_ID)
}
