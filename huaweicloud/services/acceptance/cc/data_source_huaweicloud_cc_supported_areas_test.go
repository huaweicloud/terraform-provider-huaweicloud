package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcSupportedAreas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_supported_areas.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcSupportedAreas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "areas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "areas.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "areas.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "areas.0.en_name"),
					resource.TestCheckResourceAttrSet(dataSource, "areas.0.station"),
				),
			},
		},
	})
}

func testDataSourceCcSupportedAreas_basic() string {
	return `
data "huaweicloud_cc_supported_areas" "test" {}
`
}
