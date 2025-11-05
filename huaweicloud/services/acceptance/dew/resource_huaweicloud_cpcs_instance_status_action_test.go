package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Currently, this resource is valid only in cn-north-9 region.
func TestAccCpcsInstanceStatusAction_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running test, prepare a password service cluster instance with ENABLE status.
			acceptance.TestAccPrecheckCpcsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testCpcsInstanceStatusAction_basic(),
			},
			{
				Config: testCpcsInstanceStatusAction_basic_update(),
			},
		},
	})
}

func testCpcsInstanceStatusAction_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cpcs_instance_status_action" "test" {
  instance_id = "%s"
  action      = "disable"
}
`, acceptance.HW_CPCS_INSTANCE_ID)
}

func testCpcsInstanceStatusAction_basic_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_cpcs_instance_status_action" "test" {
  instance_id = "%s"
  action      = "enable"
}
`, acceptance.HW_CPCS_INSTANCE_ID)
}
