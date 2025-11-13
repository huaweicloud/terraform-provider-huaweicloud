package cbh

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccUpgradeInstance_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCbhServerID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testUpgradeInstance_basic(),
			},
		},
	})
}

func testUpgradeInstance_basic() string {
	timestamp := time.Now().Add(48 * time.Hour).UnixMilli()

	return fmt.Sprintf(`
resource "huaweicloud_cbh_upgrade_instance" "test" {
  server_id    = "%s"
  upgrade_time = %d
}
`, acceptance.HW_CBH_SERVER_ID, timestamp)
}
