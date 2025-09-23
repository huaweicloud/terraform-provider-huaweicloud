package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAggregatorPolicyStatesSummary_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rms_resource_aggregator_policy_states_summary.test"
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
				Config: testDataSourceAggregatorPolicyStatesSummary_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "results.#"),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.resource_details.#"),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.resource_details.0.compliant_count"),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.resource_details.0.non_compliant_count"),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.assignment_details.#"),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.assignment_details.0.compliant_count"),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.assignment_details.0.non_compliant_count"),
					resource.TestCheckOutput("account_id_filter_is_useful", "true"),
					resource.TestCheckOutput("group_by_key_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceAggregatorPolicyStatesSummary_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rms_resource_aggregator" "test" {
  name        = "%[1]s"
  type        = "ACCOUNT"
  account_ids = ["%[2]s"]
}
`, name, acceptance.HW_DOMAIN_ID)
}

func testDataSourceAggregatorPolicyStatesSummary_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_resource_aggregator_policy_states_summary" "test" {
  aggregator_id = huaweicloud_rms_resource_aggregator.test.id
}

data "huaweicloud_rms_resource_aggregator_policy_states_summary" "account_id_filter" {
  aggregator_id = huaweicloud_rms_resource_aggregator.test.id
  account_id    = "%[2]s"
}
output "account_id_filter_is_useful" {
  value = length(data.huaweicloud_rms_resource_aggregator_policy_states_summary.account_id_filter.results) > 0 
}

data "huaweicloud_rms_resource_aggregator_policy_states_summary" "group_by_key_filter" {
  aggregator_id = huaweicloud_rms_resource_aggregator.test.id
  group_by_key  = "DOMAIN"
}
output "group_by_key_filter_is_useful" {
  value = length(data.huaweicloud_rms_resource_aggregator_policy_states_summary.group_by_key_filter.results) > 0 
}
`, testDataSourceAggregatorPolicyStatesSummary_base(name), acceptance.HW_DOMAIN_ID)
}
