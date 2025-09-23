package eps

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAction_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare an enterprise project and the status is enable.
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAction_basic_step1(),
			},
			{
				Config: testAccAction_basic_step2(),
			},
		},
	})
}

func testAccAction_basic_step1() string {
	return fmt.Sprintf(`
resource "huaweicloud_enterprise_project_action" "disable" {
  enterprise_project_id = "%[1]s"
  action                = "disable"
}
`, acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST)
}

func testAccAction_basic_step2() string {
	return fmt.Sprintf(`
resource "huaweicloud_enterprise_project_action" "enable" {
  enterprise_project_id = "%[1]s"
  action                = "enable"
}
`, acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST)
}
