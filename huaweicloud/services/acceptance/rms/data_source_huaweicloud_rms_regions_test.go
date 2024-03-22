package rms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsRegions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rms_regions.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRmsRegions_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_regions_not_empty", "true"),
				),
			},
			{
				Config: testDataSourceDataSourceRmsRegions_regionID,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_regions_not_empty", "true"),
					resource.TestCheckResourceAttr(dataSource, "regions.#", "1"),
					resource.TestCheckResourceAttr(dataSource, "regions.0.region_id", "cn-north-4"),
				),
			},
			{
				Config: testDataSourceDataSourceRmsRegions_displayName,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_regions_not_empty", "true"),
					resource.TestCheckResourceAttr(dataSource, "regions.#", "1"),
					resource.TestCheckResourceAttr(dataSource, "regions.0.display_name", "CN North-Beijing4"),
				),
			},
		},
	})
}

const testDataSourceDataSourceRmsRegions_basic = `
data "huaweicloud_rms_regions" "test" {}

output "is_regions_not_empty" {
  value = length(data.huaweicloud_rms_regions.test.regions) > 0
}
`

const testDataSourceDataSourceRmsRegions_regionID = `
data "huaweicloud_rms_regions" "test" {
  region_id = "cn-north-4"
}

output "is_regions_not_empty" {
  value = length(data.huaweicloud_rms_regions.test.regions) > 0
}
`

const testDataSourceDataSourceRmsRegions_displayName = `
data "huaweicloud_rms_regions" "test" {
  display_name = "CN North-Beijing4"
}

output "is_regions_not_empty" {
  value = length(data.huaweicloud_rms_regions.test.regions) > 0
}
`
