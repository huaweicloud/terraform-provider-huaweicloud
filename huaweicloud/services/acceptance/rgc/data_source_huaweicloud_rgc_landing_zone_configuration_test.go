package rgc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceLandingZoneConfiguration_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_landing_zone_configuration.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLandingZoneConfiguration_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "common_configuration.#"),
					resource.TestCheckResourceAttrSet(dataSource, "common_configuration.0.home_region"),
					resource.TestCheckResourceAttrSet(dataSource, "common_configuration.0.cloud_trail_type"),
					resource.TestCheckResourceAttrSet(dataSource, "common_configuration.0.identity_center_status"),
					resource.TestCheckResourceAttrSet(dataSource, "common_configuration.0.organization_structure_type"),
					resource.TestCheckResourceAttrSet(dataSource, "logging_configuration.#"),
					resource.TestCheckResourceAttrSet(dataSource, "logging_configuration.0.access_logging_bucket.#"),
					resource.TestCheckResourceAttrSet(dataSource, "logging_configuration.0.access_logging_bucket.0.retention_days"),
					resource.TestCheckResourceAttrSet(dataSource, "logging_configuration.0.access_logging_bucket.0.enable_multi_az"),
					resource.TestCheckResourceAttrSet(dataSource, "logging_configuration.0.logging_bucket.#"),
					resource.TestCheckResourceAttrSet(dataSource, "logging_configuration.0.logging_bucket.0.retention_days"),
					resource.TestCheckResourceAttrSet(dataSource, "logging_configuration.0.logging_bucket.0.enable_multi_az"),
					resource.TestCheckResourceAttrSet(dataSource, "organization_structure.#"),
					resource.TestCheckResourceAttrSet(dataSource, "organization_structure.0.organizational_unit_name"),
					resource.TestCheckResourceAttrSet(dataSource, "organization_structure.0.organizational_unit_type"),
					resource.TestCheckResourceAttrSet(dataSource, "organization_structure.0.accounts.#"),
					resource.TestCheckResourceAttrSet(dataSource, "organization_structure.0.accounts.0.account_id"),
					resource.TestCheckResourceAttrSet(dataSource, "organization_structure.0.accounts.0.account_type"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.region"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.region_configuration_status"),
				),
			},
		},
	})
}

const testAccDataSourceLandingZoneConfiguration_basic = `
data "huaweicloud_rgc_landing_zone_configuration" "test" {
}
`
