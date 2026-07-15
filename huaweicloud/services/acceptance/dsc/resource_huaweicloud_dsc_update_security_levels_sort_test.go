package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceDscUpdateSecurityLevelsSort_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscSecurityLevelId(t)
			acceptance.TestAccPreCheckDscSecurityLevelIdTarget(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDscUpdateSecurityLevelsSort_basic(),
			},
		},
	})
}

func testAccResourceDscUpdateSecurityLevelsSort_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_update_security_levels_sort" "test" {
  level_id        = "%[1]s"
  target_level_id = "%[2]s"
}
`, acceptance.HW_DSC_SECURITY_LEVEL_ID, acceptance.HW_DSC_SECURITY_LEVEL_ID_TARGET)
}
