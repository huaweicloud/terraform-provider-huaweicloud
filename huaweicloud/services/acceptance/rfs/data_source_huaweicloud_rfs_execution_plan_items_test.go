package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceExecutionPlanItems_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_rfs_execution_plan_items.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		name       = acceptance.RandomAccResourceNameWithDash()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceExecutionPlanItems_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "execution_plan_items.#"),
					resource.TestCheckResourceAttrSet(dataSource, "execution_plan_items.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "execution_plan_items.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "execution_plan_items.0.action"),
					resource.TestCheckResourceAttrSet(dataSource, "execution_plan_items.0.mode"),
					resource.TestCheckResourceAttrSet(dataSource, "execution_plan_items.0.provider_name"),
					resource.TestCheckResourceAttrSet(dataSource, "execution_plan_items.0.attributes.#"),
				),
			},
		},
	})
}

func testDataSourceExecutionPlanItems_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rfs_execution_plan_items" "test" {
  stack_name          = huaweicloud_rfs_stack.test.name
  execution_plan_name = huaweicloud_rfs_execution_plan_v2.test.execution_plan_name
  stack_id            = huaweicloud_rfs_stack.test.id
  execution_plan_id   = huaweicloud_rfs_execution_plan_v2.test.execution_plan_id
}
`, testAccExecutionPlanV2_basic(name))
}
