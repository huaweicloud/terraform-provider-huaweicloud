package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDesktopMaintenanceBatchManage_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_workspace_desktop_maintenance_batch_manage.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceDesktopIds(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDesktopMaintenanceBatchManage_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "desktop_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "in_maintenance_mode", "true"),
				),
			},
			{
				Config: testAccDesktopMaintenanceBatchManage_basic_step2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "desktop_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "in_maintenance_mode", "false"),
				),
			},
		},
	})
}
func testAccDesktopMaintenanceBatchManage_basic_step1() string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_desktop_maintenance_batch_manage" "test" {
  desktop_ids         = split(",", "%[1]s")
  in_maintenance_mode = true

  enable_force_new = "true"
}
`, acceptance.HW_WORKSPACE_DESKTOP_IDS)
}

func testAccDesktopMaintenanceBatchManage_basic_step2() string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_desktop_maintenance_batch_manage" "test" {
  desktop_ids         = split(",", "%[1]s")
  in_maintenance_mode = false

  enable_force_new = "true"
}
`, acceptance.HW_WORKSPACE_DESKTOP_IDS)
}
