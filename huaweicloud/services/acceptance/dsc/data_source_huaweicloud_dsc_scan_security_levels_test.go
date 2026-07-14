package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscSecurityLevels_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_dsc_scan_security_levels.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byName   = "data.huaweicloud_dsc_scan_security_levels.name_filter"
		dcByName = acceptance.InitDataSourceCheck(byName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscEnableFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscSecurityLevels_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_levels.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_levels.0.level_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_levels.0.security_level_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_levels.0.category"),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_levels.0.color_number"),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_levels.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_levels.0.sort_weight"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceDscSecurityLevels_basic = `
data "huaweicloud_dsc_scan_security_levels" "test" {
}

# Filter by name
locals {
  name = data.huaweicloud_dsc_scan_security_levels.test.security_levels[0].security_level_name
}

data "huaweicloud_dsc_scan_security_levels" "name_filter" {
  name = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dsc_scan_security_levels.name_filter.security_levels[*].security_level_name : v == local.name
  ]
}

output "is_name_filter_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}
`
