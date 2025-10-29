package dew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAvailabilityZones_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cpcs_availability_zones.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDewFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{

			{
				Config: testAccDataSourceAvailabilityZones_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zone.#"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zone.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zone.0.display_name"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zone.0.locales.#"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zone.0.locales.0.en_us"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zone.0.locales.0.zh_cn"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zone.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zone.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zone.0.status"),
				),
			},
		},
	})
}

const testAccDataSourceAvailabilityZones_basic = `data "huaweicloud_cpcs_availability_zones" "test" {}`
