package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIpBlacklistRetry_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
			acceptance.TestAccPreCheckCfwIpBlacklistName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testIpBlacklistRetry_basic(),
			},
		},
	})
}

func testIpBlacklistRetry_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cfw_ip_blacklist_retry" "test" {
  fw_instance_id = "%[1]s"
  name           = "%[2]s"
}
`, acceptance.HW_CFW_INSTANCE_ID, acceptance.HW_CFW_IP_BLACKLIST_NAME)
}
