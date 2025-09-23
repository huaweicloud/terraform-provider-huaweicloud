package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAppStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_app_statistics.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAppStatistics_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.app_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.num"),

					resource.TestCheckOutput("is_app_name_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

const testDataSourceAppStatistics_basic string = `
data "huaweicloud_hss_app_statistics" "test" {}

# Filter using app name
locals {
  app_name = data.huaweicloud_hss_app_statistics.test.data_list[0].app_name
}

data "huaweicloud_hss_app_statistics" "app_name_filter" {
  app_name = local.app_name
}

output "is_app_name_filter_useful" {
  value = length(data.huaweicloud_hss_app_statistics.app_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_app_statistics.app_name_filter.data_list[*].app_name : v == local.app_name]
  )
}

# Filter using non-existent app name
data "huaweicloud_hss_app_statistics" "not_found" {
  app_name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_app_statistics.not_found.data_list) == 0
}
`
