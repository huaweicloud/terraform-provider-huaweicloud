package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceMaskAlgorithms_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_mask_algorithms.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byName   = "data.huaweicloud_dsc_mask_algorithms.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceMaskAlgorithms_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "algorithms.#"),
					resource.TestCheckResourceAttrSet(dataSource, "algorithms.0.algorithm"),
					resource.TestCheckResourceAttrSet(dataSource, "algorithms.0.algorithm_id"),
					resource.TestCheckResourceAttrSet(dataSource, "algorithms.0.algorithm_name"),
					resource.TestCheckResourceAttrSet(dataSource, "algorithms.0.algorithm_name_en"),
					resource.TestCheckResourceAttrSet(dataSource, "algorithms.0.category"),
					resource.TestCheckResourceAttrSet(dataSource, "algorithms.0.data"),
					resource.TestCheckResourceAttrSet(dataSource, "algorithms.0.parameter"),
					resource.TestCheckResourceAttrSet(dataSource, "algorithms.0.processed_data"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceMaskAlgorithms_basic() string {
	return `
data "huaweicloud_dsc_mask_algorithms" "test" {}

locals {
  algorithm_name = data.huaweicloud_dsc_mask_algorithms.test.algorithms[0].algorithm_name
}

data "huaweicloud_dsc_mask_algorithms" "filter_by_name" {
  name = local.algorithm_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dsc_mask_algorithms.filter_by_name.algorithms[*].algorithm_name : v == local.algorithm_name
  ]
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}
`
}
