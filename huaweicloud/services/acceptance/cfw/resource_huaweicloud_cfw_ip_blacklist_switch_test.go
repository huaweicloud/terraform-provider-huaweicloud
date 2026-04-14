package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIpBlacklistSwitch_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testIpBlacklistSwitch_basic(),
			},
		},
	})
}

func testIpBlacklistSwitch_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cfw_ip_blacklist_switch" "test" {
  fw_instance_id = "%s"
  status         = 1
}
`, acceptance.HW_CFW_INSTANCE_ID)
}
