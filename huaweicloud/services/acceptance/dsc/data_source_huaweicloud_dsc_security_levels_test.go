package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscJobSecurityLevels_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dsc_security_levels.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscScanJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscJobSecurityLevels_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "security_level_list.#"),
					resource.TestCheckOutput("asset_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceDscJobSecurityLevels_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dsc_security_levels" "test" {
  job_id = "%[1]s"
}

# Filter by asset_type

data "huaweicloud_dsc_security_levels" "filter_by_asset_type" {
  job_id     = "%[1]s"
  asset_type = "OBS"
}

output "asset_type_filter_is_useful" {
  value = length(data.huaweicloud_dsc_security_levels.filter_by_asset_type.security_level_list) > 0
}
`, acceptance.HW_DSC_SCAN_JOB_ID)
}
