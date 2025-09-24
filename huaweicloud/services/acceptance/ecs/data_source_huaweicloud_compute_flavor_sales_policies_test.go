package ecs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEcsComputeFlavorSalesPolicies_basic(t *testing.T) {
	dataSource := "data.huaweicloud_compute_flavor_sales_policies.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEcsComputeFlavorSalesPolicies_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "sell_policies.#"),
					resource.TestCheckResourceAttrSet(dataSource, "sell_policies.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "sell_policies.0.flavor_id"),
					resource.TestCheckResourceAttrSet(dataSource, "sell_policies.0.sell_status"),
					resource.TestCheckResourceAttrSet(dataSource, "sell_policies.0.availability_zone_id"),
					resource.TestCheckResourceAttrSet(dataSource, "sell_policies.0.sell_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "sell_policies.0.spot_options.#"),
					resource.TestCheckResourceAttrSet(dataSource, "sell_policies.0.spot_options.0.interruption_policy"),
					resource.TestCheckResourceAttrSet(dataSource, "sell_policies.0.spot_options.0.longest_spot_duration_hours"),
					resource.TestCheckResourceAttrSet(dataSource, "sell_policies.0.spot_options.0.largest_spot_duration_count"),
					resource.TestCheckOutput("flavor_id_filter_is_useful", "true"),
					resource.TestCheckOutput("sell_status_filter_is_useful", "true"),
					resource.TestCheckOutput("sell_mode_filter_is_useful", "true"),
					resource.TestCheckOutput("availability_zone_id_filter_is_useful", "true"),
					resource.TestCheckOutput("longest_spot_duration_hours_filter_is_useful", "true"),
					resource.TestCheckOutput("largest_spot_duration_count_filter_is_useful", "true"),
					resource.TestCheckOutput("interruption_policy_filter_is_useful", "true"),
					resource.TestCheckOutput("longest_spot_duration_hours_gt_filter_is_useful", "true"),
					resource.TestCheckOutput("largest_spot_duration_count_gt_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceEcsComputeFlavorSalesPolicies_basic() string {
	return `
data "huaweicloud_compute_flavor_sales_policies" "test" {}

locals {
  policies = data.huaweicloud_compute_flavor_sales_policies.test.sell_policies
}

data "huaweicloud_compute_flavor_sales_policies" "flavor_id_filter" {
  flavor_id = local.policies[0].flavor_id
}
locals {
  flavor_id = local.policies[0].flavor_id
}
output "flavor_id_filter_is_useful" {
  value = length(data.huaweicloud_compute_flavor_sales_policies.flavor_id_filter.sell_policies) > 0 && alltrue(
  [for v in data.huaweicloud_compute_flavor_sales_policies.flavor_id_filter.sell_policies[*] : v.flavor_id == local.flavor_id]
  )
}

data "huaweicloud_compute_flavor_sales_policies" "sell_status_filter" {
  sell_status = local.policies[0].sell_status
}
locals {
  sell_status = local.policies[0].sell_status
}
output "sell_status_filter_is_useful" {
  value = length(data.huaweicloud_compute_flavor_sales_policies.sell_status_filter.sell_policies) > 0 && alltrue(
  [for v in data.huaweicloud_compute_flavor_sales_policies.sell_status_filter.sell_policies[*] : v.sell_status == local.sell_status]
  )
}

data "huaweicloud_compute_flavor_sales_policies" "sell_mode_filter" {
  sell_mode = local.policies[0].sell_mode
}
locals {
  sell_mode = local.policies[0].sell_mode
}
output "sell_mode_filter_is_useful" {
  value = length(data.huaweicloud_compute_flavor_sales_policies.sell_mode_filter.sell_policies) > 0 && alltrue(
  [for v in data.huaweicloud_compute_flavor_sales_policies.sell_mode_filter.sell_policies[*] : v.sell_mode == local.sell_mode]
  )
}

data "huaweicloud_compute_flavor_sales_policies" "availability_zone_id_filter" {
  availability_zone_id = local.policies[0].availability_zone_id
}
locals {
  availability_zone_id = local.policies[0].availability_zone_id
}
output "availability_zone_id_filter_is_useful" {
  value = length(data.huaweicloud_compute_flavor_sales_policies.availability_zone_id_filter.sell_policies) > 0 && alltrue(
  [for v in data.huaweicloud_compute_flavor_sales_policies.availability_zone_id_filter.sell_policies[*] :
  v.availability_zone_id == local.availability_zone_id]
  )
}

data "huaweicloud_compute_flavor_sales_policies" "longest_spot_duration_hours_filter" {
  longest_spot_duration_hours = local.policies[0].spot_options[0].longest_spot_duration_hours
}
locals {
  longest_spot_duration_hours = local.policies[0].spot_options[0].longest_spot_duration_hours
}
output "longest_spot_duration_hours_filter_is_useful" {
  value = length(data.huaweicloud_compute_flavor_sales_policies.longest_spot_duration_hours_filter.sell_policies) > 0 && alltrue(
  [for v in data.huaweicloud_compute_flavor_sales_policies.longest_spot_duration_hours_filter.sell_policies[*] :
  v.spot_options[0].longest_spot_duration_hours == local.longest_spot_duration_hours]
  )
}

data "huaweicloud_compute_flavor_sales_policies" "largest_spot_duration_count_filter" {
  largest_spot_duration_count = local.policies[0].spot_options[0].largest_spot_duration_count
}
locals {
  largest_spot_duration_count = local.policies[0].spot_options[0].largest_spot_duration_count
}
output "largest_spot_duration_count_filter_is_useful" {
  value = length(data.huaweicloud_compute_flavor_sales_policies.largest_spot_duration_count_filter.sell_policies) > 0 && alltrue(
  [for v in data.huaweicloud_compute_flavor_sales_policies.largest_spot_duration_count_filter.sell_policies[*] :
  v.spot_options[0].largest_spot_duration_count == local.largest_spot_duration_count]
  )
}

data "huaweicloud_compute_flavor_sales_policies" "interruption_policy_filter" {
  interruption_policy = local.policies[0].spot_options[0].interruption_policy
}
locals {
  interruption_policy = local.policies[0].spot_options[0].interruption_policy
}
output "interruption_policy_filter_is_useful" {
  value = length(data.huaweicloud_compute_flavor_sales_policies.interruption_policy_filter.sell_policies) > 0 && alltrue(
  [for v in data.huaweicloud_compute_flavor_sales_policies.interruption_policy_filter.sell_policies[*] :
  v.spot_options[0].interruption_policy == local.interruption_policy]
  )
}

data "huaweicloud_compute_flavor_sales_policies" "longest_spot_duration_hours_gt_filter" {
  longest_spot_duration_hours_gt = 1
}
locals {
  longest_spot_duration_hours_gt = 1
}
output "longest_spot_duration_hours_gt_filter_is_useful" {
  value = length(data.huaweicloud_compute_flavor_sales_policies.longest_spot_duration_hours_gt_filter.sell_policies) > 0 && alltrue(
  [for v in data.huaweicloud_compute_flavor_sales_policies.longest_spot_duration_hours_gt_filter.sell_policies[*] :
  v.spot_options[0].longest_spot_duration_hours > local.longest_spot_duration_hours_gt]
  )
}

data "huaweicloud_compute_flavor_sales_policies" "largest_spot_duration_count_gt_filter" {
  largest_spot_duration_count_gt = 1
}
locals {
  largest_spot_duration_count_gt = 1
}
output "largest_spot_duration_count_gt_filter_is_useful" {
  value = length(data.huaweicloud_compute_flavor_sales_policies.largest_spot_duration_count_gt_filter.sell_policies) > 0 && alltrue(
  [for v in data.huaweicloud_compute_flavor_sales_policies.largest_spot_duration_count_gt_filter.sell_policies[*] :
  v.spot_options[0].largest_spot_duration_count > local.largest_spot_duration_count_gt]
  )
}
`
}
