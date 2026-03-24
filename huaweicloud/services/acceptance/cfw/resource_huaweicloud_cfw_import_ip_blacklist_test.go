package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccImportIpBlacklist_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testImportIpBlacklist_basic(),
			},
		},
	})
}

func testImportIpBlacklist_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cfw_import_ip_blacklist" "test" {
  fw_instance_id = "%s"
  add_type       = 1
  ip_blacklist   = "100.1.1.10"
  effect_scope   = [1]
}
`, acceptance.HW_CFW_INSTANCE_ID)
}
