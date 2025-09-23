package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDataBacklogPolicies_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_iotda_data_backlog_policies.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This resource only supports standard and enterprise version IoTDA instances.
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataBacklogPolicies_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "policies.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.backlog_size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.backlog_time"),

					resource.TestCheckOutput("is_policy_name_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataSourceDataBacklogPolicies_basic() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%s

data "huaweicloud_iotda_data_backlog_policies" "test" {
  depends_on = [
    huaweicloud_iotda_data_backlog_policy.test,
  ]
}

# Filter using policy_name.
locals {
  policy_name = data.huaweicloud_iotda_data_backlog_policies.test.policies[0].name
}

data "huaweicloud_iotda_data_backlog_policies" "policy_name_filter" {
  policy_name = local.policy_name
}

output "is_policy_name_filter_useful" {
  value = length(data.huaweicloud_iotda_data_backlog_policies.policy_name_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_data_backlog_policies.policy_name_filter.policies[*].name : v == local.policy_name]
  )
}

# Filter using non existent name.
data "huaweicloud_iotda_data_backlog_policies" "not_found" {
  policy_name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_iotda_data_backlog_policies.not_found.policies) == 0
}
`, testDataBacklogPolicy_basic(name))
}
