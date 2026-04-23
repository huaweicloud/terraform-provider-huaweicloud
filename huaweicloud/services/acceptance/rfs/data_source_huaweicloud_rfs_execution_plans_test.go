package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRfsExecutionPlans_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_rfs_execution_plans.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRfsStackName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRfsExecutionPlans_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "execution_plans.#"),
					resource.TestCheckResourceAttrSet(dataSource, "execution_plans.0.stack_name"),
					resource.TestCheckResourceAttrSet(dataSource, "execution_plans.0.stack_id"),
					resource.TestCheckResourceAttrSet(dataSource, "execution_plans.0.execution_plan_id"),
					resource.TestCheckResourceAttrSet(dataSource, "execution_plans.0.execution_plan_name"),
					resource.TestCheckResourceAttrSet(dataSource, "execution_plans.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "execution_plans.0.create_time"),
				),
			},
		},
	})
}

func testAccDataSourceRfsExecutionPlans_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rfs_execution_plans" "test" {
  stack_name = "%s"
}
`, acceptance.HW_RFS_STACK_NAME)
}
