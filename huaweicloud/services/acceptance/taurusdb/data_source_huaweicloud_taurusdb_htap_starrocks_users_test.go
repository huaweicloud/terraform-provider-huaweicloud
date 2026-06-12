package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBHtapStarrocksUsers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_htap_starrocks_users.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBHtapInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTaurusDBHtapStarrocksUsers_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "user_details.#"),
					resource.TestCheckResourceAttrSet(dataSource, "user_details.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "user_details.0.dml"),
					resource.TestCheckResourceAttrSet(dataSource, "user_details.0.ddl"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceTaurusDBHtapStarrocksUsers_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_taurusdb_htap_starrocks_users" "test" {
  instance_id = "%[1]s"
}

locals {
  user_name = data.huaweicloud_taurusdb_htap_starrocks_users.test.user_details[0].user_name
}

data "huaweicloud_taurusdb_htap_starrocks_users" "filter" {
  instance_id = "%[1]s"
  user_name   = local.user_name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_taurusdb_htap_starrocks_users.filter.user_details) > 0 && alltrue(
    [for v in data.huaweicloud_taurusdb_htap_starrocks_users.filter.user_details : v.user_name == local.user_name]
  )
}
`, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID)
}
