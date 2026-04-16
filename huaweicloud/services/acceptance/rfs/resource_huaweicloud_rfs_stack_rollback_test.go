package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccStackRollback_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRfsStackName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStackRollback_basic(),
			},
		},
	})
}

func testAccStackRollback_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_rfs_stack_rollback" "test" {
  stack_name = "%s"
}
`, acceptance.HW_RFS_STACK_NAME)
}
