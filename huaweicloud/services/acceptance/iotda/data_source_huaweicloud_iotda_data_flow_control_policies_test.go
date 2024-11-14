package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDataFlowControlPolicies_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_iotda_data_flow_control_policies.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Only standard and enterprise IoTDA instances support this resource.
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataFlowControlPolicies_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "policies.#", "2"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.scope"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.scope_value"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.limit"),

					resource.TestCheckOutput("filter_with_scope", "2"),
					resource.TestCheckOutput("filter_with_scope_value", "1"),
					// When using `policy_name` filtering, it is a fuzzy match.
					resource.TestCheckOutput("filter_with_policy_name", "2"),
					resource.TestCheckOutput("not_found_validation_pass", "0"),
				),
			},
		},
	})
}

func testDataSourceDataFlowControlPolicies_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_data_flow_control_policy" "test1" {
  name        = "%[2]s_1"
  description = "description_test"
  scope       = "CHANNEL"
  scope_value = "DMS_KAFKA_FORWARDING"
}

resource "huaweicloud_iotda_data_flow_control_policy" "test2" {
  name        = "%[2]s_2"
  description = "description_test"
  scope       = "CHANNEL"
  scope_value = "DIS_FORWARDING"
}
`, buildIoTDAEndpoint(), name)
}

func testAccDataSourceDataFlowControlPolicies_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_iotda_data_flow_control_policies" "test" {
  depends_on = [
    huaweicloud_iotda_data_flow_control_policy.test1,
    huaweicloud_iotda_data_flow_control_policy.test2,
  ]
}

# Filter using scope.
locals {
  scope = data.huaweicloud_iotda_data_flow_control_policies.test.policies[0].scope
}

data "huaweicloud_iotda_data_flow_control_policies" "scope_filter" {
  scope = local.scope
}

output "filter_with_scope" {
  value = length(data.huaweicloud_iotda_data_flow_control_policies.scope_filter.policies)
}

# Filter using scope value.
locals {
  scope_value = data.huaweicloud_iotda_data_flow_control_policies.test.policies[0].scope_value
}

data "huaweicloud_iotda_data_flow_control_policies" "scope_value_filter" {
  scope       = local.scope
  scope_value = local.scope_value
}

output "filter_with_scope_value" {
  value = length(data.huaweicloud_iotda_data_flow_control_policies.scope_value_filter.policies)
}

# Filter using policy name.
locals {
  policy_name = data.huaweicloud_iotda_data_flow_control_policies.test.policies[0].name
}

data "huaweicloud_iotda_data_flow_control_policies" "policy_name_filter" {
  policy_name = local.policy_name
}

output "filter_with_policy_name" {
  value = length(data.huaweicloud_iotda_data_flow_control_policies.scope_filter.policies)
}

# Filter using non existent name.
data "huaweicloud_iotda_data_flow_control_policies" "not_found" {
  policy_name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_iotda_data_flow_control_policies.not_found.policies)
}
`, testDataSourceDataFlowControlPolicies_base())
}
