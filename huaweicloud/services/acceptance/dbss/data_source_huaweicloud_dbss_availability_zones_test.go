package dbss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAvailabilityZones_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dbss_availability_zones.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAvailabilityZones_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "availability_zones.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "availability_zones.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "availability_zones.0.number"),
					resource.TestCheckResourceAttrSet(dataSourceName, "availability_zones.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "availability_zones.0.alias"),
					resource.TestCheckResourceAttrSet(dataSourceName, "availability_zones.0.alias_us"),
				),
			},
		},
	})
}

const testDataSourceAvailabilityZones_basic = `data "huaweicloud_dbss_availability_zones" "test" {}`
