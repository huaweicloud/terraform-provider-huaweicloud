package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePolicyGroups_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_policy_groups.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID with host protection enabled.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePolicyGroups_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.group_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.deletable"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.default_group"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.support_os"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.support_version"),

					resource.TestCheckOutput("is_group_id_filter_useful", "true"),
					resource.TestCheckOutput("is_group_name_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

const testDataSourcePolicyGroups_basic string = `
data "huaweicloud_hss_policy_groups" "test" {}

# Filter using group id.
locals {
  group_id = data.huaweicloud_hss_policy_groups.test.data_list[0].group_id
}

data "huaweicloud_hss_policy_groups" "group_id_filter" {
  group_id = local.group_id
}

output "is_group_id_filter_useful" {
  value = length(data.huaweicloud_hss_policy_groups.group_id_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_policy_groups.group_id_filter.data_list[*].group_id : v == local.group_id]
  )
}

# Filter using group name.
locals {
  group_name = data.huaweicloud_hss_policy_groups.test.data_list[0].group_name
}

data "huaweicloud_hss_policy_groups" "group_name_filter" {
  group_name = local.group_name
}

output "is_group_name_filter_useful" {
  value = length(data.huaweicloud_hss_policy_groups.group_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_policy_groups.group_name_filter.data_list[*].group_name : v == local.group_name]
  )
}

# Filter using non group name.
data "huaweicloud_hss_policy_groups" "not_found" {
  group_name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_policy_groups.not_found.data_list) == 0
}
`
