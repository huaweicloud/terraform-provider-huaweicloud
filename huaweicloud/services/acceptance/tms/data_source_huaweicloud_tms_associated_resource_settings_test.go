package tms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTmsAssociatedResourceSettings_basic(t *testing.T) {
	dataSource := "data.huaweicloud_tms_associated_resource_settings.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// The step covers the query behavior (GET /v2/tms/associated-resources/settings)
				// and verifies the computed attributes of the returned settings.
				Config: testAccTmsAssociatedResourceSettingsDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "settings.#"),
					resource.TestCheckResourceAttrSet(dataSource, "settings.0.setting_name"),
					resource.TestCheckResourceAttrSet(dataSource, "settings.0.master_service"),
					resource.TestCheckResourceAttrSet(dataSource, "settings.0.master_resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "settings.0.associated_service"),
					resource.TestCheckResourceAttrSet(dataSource, "settings.0.associated_resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "settings.0.support_existing_resource"),
					resource.TestCheckResourceAttrSet(dataSource, "settings.0.support_auto_delete"),
					resource.TestCheckResourceAttrSet(dataSource, "settings.0.region_ids.#"),
					resource.TestCheckOutput("is_useful", "true"),
				),
			},
		},
	})
}

func TestAccDataSourceTmsAssociatedResourceSettings_filter(t *testing.T) {
	dataSource := "data.huaweicloud_tms_associated_resource_settings.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// The step covers the query behavior filtered by `region_id` and returns
				// only the settings applicable to the specified region.
				Config: testAccTmsAssociatedResourceSettingsDataSource_filter(acceptance.HW_REGION_NAME),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "settings.#"),
					resource.TestCheckOutput("is_useful", "true"),
				),
			},
		},
	})
}

func testAccTmsAssociatedResourceSettingsDataSource_basic() string {
	return `
data "huaweicloud_tms_associated_resource_settings" "test" {}

locals {
  settings = data.huaweicloud_tms_associated_resource_settings.test.settings
}

output "is_useful" {
  value = length(local.settings) > 0
}
`
}

func testAccTmsAssociatedResourceSettingsDataSource_filter(regionId string) string {
	return fmt.Sprintf(`
data "huaweicloud_tms_associated_resource_settings" "test" {
  region_id = "%s"
}

locals {
  settings = data.huaweicloud_tms_associated_resource_settings.test.settings
}

output "is_useful" {
  value = length(local.settings) > 0
}
`, regionId)
}
