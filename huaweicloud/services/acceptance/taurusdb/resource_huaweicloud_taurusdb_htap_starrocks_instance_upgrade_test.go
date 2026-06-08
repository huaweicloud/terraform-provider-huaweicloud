package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTaurusDBHtapStarrocksInstanceUpgrade_basic(t *testing.T) {
	resourceName := "huaweicloud_taurusdb_htap_starrocks_instance_upgrade.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
			acceptance.TestAccPreCheckTaurusDBHtapInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBHtapStarrocksInstanceUpgrade_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccTaurusDBHtapStarrocksInstanceUpgrade_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_taurusdb_htap_starrocks_instance_upgrade" "test" {
  instance_id           = "%[1]s"
  starrocks_instance_id = "%[2]s"
  delay                 = "true"
  is_skip_validate      = "true"
}
`, acceptance.HW_TAURUSDB_INSTANCE_ID, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID)
}
