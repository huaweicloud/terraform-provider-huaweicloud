package das

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccLockBlockingSwitch_basic(t *testing.T) {
	var (
		rName = "huaweicloud_das_lock_blocking_switch.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDasInstanceIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccLockBlockingSwitch_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "switch_on", "false"),
				),
			},
			{
				Config: testAccLockBlockingSwitch_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "switch_on", "true"),
					resource.TestCheckResourceAttr(rName, "retention_hours", "300"),
				),
			},
		},
	})
}

func testAccLockBlockingSwitch_base() string {
	return fmt.Sprintf(`
locals {
  instance_ids = split(",", "%[1]s")
}
`, acceptance.HW_DAS_INSTANCE_IDS)
}

func testAccLockBlockingSwitch_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_lock_blocking_switch" "test" {
  instance_id = local.instance_ids[0]
  switch_on   = false
}
`, testAccLockBlockingSwitch_base())
}

func testAccLockBlockingSwitch_basic_update() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_lock_blocking_switch" "test" {
  instance_id     = local.instance_ids[0]
  switch_on       = true
  retention_hours = 300

  enable_force_new = "true"
}
`, testAccLockBlockingSwitch_base())
}
