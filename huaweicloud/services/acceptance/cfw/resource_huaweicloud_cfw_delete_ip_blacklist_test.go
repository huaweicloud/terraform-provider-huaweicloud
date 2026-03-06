package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDeleteIpBlacklist_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDeleteIpBlacklist_basic(),
			},
		},
	})
}

func testDeleteIpBlacklist_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cfw_delete_ip_blacklist" "test" {
  fw_instance_id = "%s"
  effect_scope   = [1, 2]
}
`, acceptance.HW_CFW_INSTANCE_ID)
}
