package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccPoliciesDataSource_basic(t *testing.T) {
	var (
		rName          = acceptance.RandomAccResourceName()
		dataSourceName = "data.huaweicloud_as_policies.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPoliciesDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.scaling_group_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.scaling_resource_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.scaling_resource_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.alarm_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.action.0.operation"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.action.0.instance_number"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.cool_down_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.created_at"),
					resource.TestCheckOutput("is_scaling_policy_id_filter_useful", "true"),
					resource.TestCheckOutput("is_scaling_policy_name_filter_useful", "true"),
					resource.TestCheckOutput("is_scaling_policy_type_filter_useful", "true"),
					resource.TestCheckOutput("is_scaling_group_id_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_alarm_id_filter_useful", "true"),
				),
			},
		},
	})
}

func TestAccPoliciesDataSource_bandwidth(t *testing.T) {
	var (
		rName          = acceptance.RandomAccResourceName()
		dataSourceName = "data.huaweicloud_as_policies.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPoliciesDataSource_bandwidth(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.scaling_resource_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.scaling_resource_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.alarm_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.action.0.operation"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.action.0.instance_number"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.action.0.limits"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.cool_down_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.created_at"),
					resource.TestCheckOutput("is_scaling_policy_id_filter_useful", "true"),
					resource.TestCheckOutput("is_scaling_policy_name_filter_useful", "true"),
					resource.TestCheckOutput("is_scaling_resource_id_filter_useful", "true"),
					resource.TestCheckOutput("is_scaling_resource_type_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccPoliciesDataSource_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_as_policies" "test" {
  depends_on = [
    huaweicloud_as_policy.acc_as_policy
  ]
}

// Filter using scaling policy ID.
locals {
  scaling_policy_id = data.huaweicloud_as_policies.test.policies[0].id
}

data "huaweicloud_as_policies" "scaling_policy_id_filter" {
  scaling_policy_id = local.scaling_policy_id
}

output "is_scaling_policy_id_filter_useful" {
  value = length(data.huaweicloud_as_policies.scaling_policy_id_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_as_policies.scaling_policy_id_filter.policies[*].id : v == local.scaling_policy_id]
  )
}

// Filter using scaling policy name.
locals {
  scaling_policy_name = data.huaweicloud_as_policies.test.policies[0].name
}

data "huaweicloud_as_policies" "scaling_policy_name_filter" {
  scaling_policy_name = local.scaling_policy_name
}

output "is_scaling_policy_name_filter_useful" {
  value = length(data.huaweicloud_as_policies.scaling_policy_name_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_as_policies.scaling_policy_name_filter.policies[*].name : v == local.scaling_policy_name]
  )
}

// Filter using scaling policy type.
locals {
  scaling_policy_type = data.huaweicloud_as_policies.test.policies[0].type
}

data "huaweicloud_as_policies" "scaling_policy_type_filter" {
  scaling_policy_type = local.scaling_policy_type
}

output "is_scaling_policy_type_filter_useful" {
  value = length(data.huaweicloud_as_policies.scaling_policy_type_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_as_policies.scaling_policy_type_filter.policies[*].type : v == local.scaling_policy_type]
  )
}

// Filter using scaling group ID.
locals {
  scaling_group_id = data.huaweicloud_as_policies.test.policies[0].scaling_group_id
}

data "huaweicloud_as_policies" "scaling_group_id_filter" {
  scaling_group_id = local.scaling_group_id
}

output "is_scaling_group_id_filter_useful" {
  value = length(data.huaweicloud_as_policies.scaling_group_id_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_as_policies.scaling_group_id_filter.policies[*].scaling_group_id : v == local.scaling_group_id]
  )
}

// Filter using status.
locals {
  status = data.huaweicloud_as_policies.test.policies[0].status
}

data "huaweicloud_as_policies" "status_filter" {
  status = local.status
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_as_policies.status_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_as_policies.status_filter.policies[*].status : v == local.status]
  )
}

// Filter using alarm ID.
locals {
  alarm_id = data.huaweicloud_as_policies.test.policies[0].alarm_id
}

data "huaweicloud_as_policies" "alarm_id_filter" {
  alarm_id = local.alarm_id
}

output "is_alarm_id_filter_useful" {
  value = length(data.huaweicloud_as_policies.alarm_id_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_as_policies.alarm_id_filter.policies[*].alarm_id : v == local.alarm_id]
  )
}
`, testASPolicy_alarm(name))
}

func testAccPoliciesDataSource_bandwidth(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_as_policies" "test" {
  depends_on = [
    huaweicloud_as_bandwidth_policy.test
  ]
}

// Filter using scaling policy ID.
locals {
  scaling_policy_id = data.huaweicloud_as_policies.test.policies[0].id
}

data "huaweicloud_as_policies" "scaling_policy_id_filter" {
  scaling_policy_id = local.scaling_policy_id
}

output "is_scaling_policy_id_filter_useful" {
  value = length(data.huaweicloud_as_policies.scaling_policy_id_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_as_policies.scaling_policy_id_filter.policies[*].id : v == local.scaling_policy_id]
  )
}

// Filter using scaling policy name.
locals {
  scaling_policy_name = data.huaweicloud_as_policies.test.policies[0].name
}

data "huaweicloud_as_policies" "scaling_policy_name_filter" {
  scaling_policy_name = local.scaling_policy_name
}

output "is_scaling_policy_name_filter_useful" {
  value = length(data.huaweicloud_as_policies.scaling_policy_name_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_as_policies.scaling_policy_name_filter.policies[*].name : v == local.scaling_policy_name]
  )
}

// Filter using scaling resource ID
locals {
  scaling_resource_id = data.huaweicloud_as_policies.test.policies[0].scaling_resource_id
}

data "huaweicloud_as_policies" "scaling_resource_id_filter" {
  scaling_resource_id = local.scaling_resource_id
}

output "is_scaling_resource_id_filter_useful" {
  value = length(data.huaweicloud_as_policies.scaling_resource_id_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_as_policies.scaling_resource_id_filter.policies[*].scaling_resource_id : 
      v == local.scaling_resource_id]
  )
}

// Filter using scaling resource type
locals {
  scaling_resource_type = data.huaweicloud_as_policies.test.policies[0].scaling_resource_type
}

data "huaweicloud_as_policies" "scaling_resource_type_filter" {
  scaling_resource_type = local.scaling_resource_type
}

output "is_scaling_resource_type_filter_useful" {
  value = length(data.huaweicloud_as_policies.scaling_resource_type_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_as_policies.scaling_resource_type_filter.policies[*].scaling_resource_type : 
      v == local.scaling_resource_type]
  )
}

// Filter using status
locals {
  status = data.huaweicloud_as_policies.test.policies[0].status
}

data "huaweicloud_as_policies" "status_filter" {
  status = local.status
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_as_policies.status_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_as_policies.status_filter.policies[*].status : v == local.status]
  )
}
`, testASBandWidthPolicy_alarm(name))
}
