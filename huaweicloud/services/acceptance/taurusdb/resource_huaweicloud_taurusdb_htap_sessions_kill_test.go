package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTaurusDBHtapSessionsKill_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBHtapInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBHtapSessionsKill_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_taurusdb_htap_sessions_kill.test", "status", "success"),
				),
			},
		},
	})
}

func testAccTaurusDBHtapSessionsKill_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_taurusdb_htap_sessions" "test"{
  instance_id = "%[1]s"
}

resource "huaweicloud_taurusdb_htap_sessions_kill" "test" {
  instance_id  = "%[1]s"
  process_list = [data.huaweicloud_taurusdb_htap_sessions.test.process_list[0].id]
}
`, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID)
}
