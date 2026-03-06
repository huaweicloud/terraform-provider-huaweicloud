package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccExportIpBlacklist_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testExportIpBlacklist_basic(),
			},
		},
	})
}

func testExportIpBlacklist_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cfw_export_ip_blacklist" "test" {
  fw_instance_id = "%s"
  name           = "ip-blacklist-eip.txt"
}
`, acceptance.HW_CFW_INSTANCE_ID)
}
