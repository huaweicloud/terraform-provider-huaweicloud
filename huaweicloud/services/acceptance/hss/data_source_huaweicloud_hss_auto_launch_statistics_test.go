package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAutoLaunchStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_auto_launch_statistics.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAutoLaunchStatistics_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.num"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

const testDataSourceAutoLaunchStatistics_basic string = `
data "huaweicloud_hss_auto_launch_statistics" "test" {}

# Filter using name
locals {
  name = data.huaweicloud_hss_auto_launch_statistics.test.data_list[0].name
}

data "huaweicloud_hss_auto_launch_statistics" "name_filter" {
  name = local.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_hss_auto_launch_statistics.name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_auto_launch_statistics.name_filter.data_list[*].name : v == local.name]
  )
}

# Filter using type
locals {
  type = data.huaweicloud_hss_auto_launch_statistics.test.data_list[0].type
}

data "huaweicloud_hss_auto_launch_statistics" "type_filter" {
  type = local.type
}

output "is_type_filter_useful" {
  value = length(data.huaweicloud_hss_auto_launch_statistics.type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_auto_launch_statistics.type_filter.data_list[*].type : v == local.type]
  )
}

# Filter using non-existent name
data "huaweicloud_hss_auto_launch_statistics" "not_found" {
  name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_auto_launch_statistics.not_found.data_list) == 0
}
`
