package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEipAllProtectionSwitch_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEipAllProtectionSwitch_basic(),
			},
		},
	})
}

func testAccEipAllProtectionSwitch_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cfw_eip_all_protection_switch" "test" {
  fw_instance_id   = "%s"
  bypass_operation = 0
}
`, acceptance.HW_CFW_INSTANCE_ID)
}
