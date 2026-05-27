package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTaurusDBHtapStarrocksInstanceRestart_basic(t *testing.T) {
	resourceName := "huaweicloud_taurusdb_htap_starrocks_instance_restart.test"

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
				Config: testAccTaurusDBHtapStarrocksInstanceRestart_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID),
				),
			},
		},
	})
}

func testAccTaurusDBHtapStarrocksInstanceRestart_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_taurusdb_htap_starrocks_instance_restart" "test" {
  taurusdb_instance_id  = "%[1]s"
  starrocks_instance_id = "%[2]s"
}
`, acceptance.HW_TAURUSDB_INSTANCE_ID, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID)
}
