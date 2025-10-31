package rgc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceLandingZoneAvailableUpdates_basic(t *testing.T) {

	dataSource := "data.huaweicloud_rgc_landing_zone_available_updates.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLandingZoneAvailableUpdates_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "baseline_update_available"),
					resource.TestCheckResourceAttrSet(dataSource, "control_update_available"),
					resource.TestCheckResourceAttrSet(dataSource, "landing_zone_update_available"),
					resource.TestCheckResourceAttrSet(dataSource, "service_landing_zone_version"),
					resource.TestCheckResourceAttrSet(dataSource, "user_landing_zone_version"),
				),
			},
		},
	})
}

const testAccDataSourceLandingZoneAvailableUpdates_basic string = `
data "huaweicloud_rgc_landing_zone_available_updates" "test" {}
`
