package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccApplyExecutionPlan_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccApplyExecutionPlan_basic(name),
			},
		},
	})
}

func testAccApplyExecutionPlan_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rfs_apply_execution_plan" "test" {
  depends_on = [huaweicloud_rfs_execution_plan_v2.test]

  stack_name          = "%[2]s"
  execution_plan_name = "%[2]s"
}
`, testAccExecutionPlanV2_basic(name), name)
}
