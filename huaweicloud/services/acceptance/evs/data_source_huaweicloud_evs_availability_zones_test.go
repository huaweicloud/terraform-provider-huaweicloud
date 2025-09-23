package evs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEvsAvailabilityZones_basic(t *testing.T) {
	dataSource := "data.huaweicloud_evs_availability_zones.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceEvsAvailabilityZones_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.#"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.is_available"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.name"),
				),
			},
		},
	})
}

func testDataSourceDataSourceEvsAvailabilityZones_basic() string {
	return `
data "huaweicloud_evs_availability_zones" "test" {}
`
}
