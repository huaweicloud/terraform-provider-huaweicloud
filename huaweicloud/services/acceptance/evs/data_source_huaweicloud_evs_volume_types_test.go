package evs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEvsVolumeTypes_basic(t *testing.T) {
	dataSource := "data.huaweicloud_evs_volume_types.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceEvsVolumeTypes_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "types.#"),
					resource.TestCheckResourceAttrSet(dataSource, "types.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "types.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "types.0.extra_specs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "types.0.extra_specs.0.available_availability_zones"),
				),
			},
		},
	})
}

func testDataSourceDataSourceEvsVolumeTypes_basic() string {
	return `
data "huaweicloud_evs_volume_types" "test" {}
`
}
