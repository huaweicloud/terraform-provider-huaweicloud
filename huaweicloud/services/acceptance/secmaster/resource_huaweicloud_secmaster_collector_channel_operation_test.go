package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCollectorChannelOperation_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterChannelID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccCollectorChannelOperation_basic(),
			},
		},
	})
}

func testAccCollectorChannelOperation_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_collector_channel_operation" "test" {
  workspace_id = "%[1]s"
  channel_id   = "%[2]s"
  action       = "STOP"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_CHANNEL_ID)
}
