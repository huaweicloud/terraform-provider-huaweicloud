package dws

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAvailabilityZones_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dws_availability_zones.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

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
					resource.TestMatchResourceAttr(dataSource, "availability_zones.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.display_name"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.status"),
				),
			},
		},
	})
}

const testDataSourceAvailabilityZones_basic = `data "huaweicloud_dws_availability_zones" "test" {}`
