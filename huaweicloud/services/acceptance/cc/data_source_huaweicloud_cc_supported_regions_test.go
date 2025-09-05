package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcSupportedRegions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_supported_regions.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcSupportedRegions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "regions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.area_id"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.area_name"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.used_scenes.#"),
				),
			},
		},
	})
}

func testDataSourceCcSupportedRegions_basic() string {
	return `
data "huaweicloud_cc_supported_regions" "test" {}
`
}
