package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceExecutionPlanPrices_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_rfs_execution_plan_prices.test"
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
				Config: testDataSourceExecutionPlanPrices_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "items.#"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.supported"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource_price.#"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource_price.0.charge_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource_price.0.sale_price"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource_price.0.discount"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource_price.0.original_price"),
				),
			},
		},
	})
}

func testDataSourceExecutionPlanPrices_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rfs_stack" "test" {
  name = "%[1]s"
}

resource "huaweicloud_rfs_execution_plan_v2" "test" {
  stack_name          = huaweicloud_rfs_stack.test.name
  execution_plan_name = "%[1]s"
  description         = "test execution plan"
  stack_id            = huaweicloud_rfs_stack.test.id
  template_body       = %[2]s
  vars_body           = %[3]s

  vars_structure {
    var_key   = "vpc_name"
    var_value = "%[1]s-vpc"
  }

  vars_structure {
    var_key   = "subnet_name"
    var_value = "%[1]s-subnet"
  }

  vars_structure {
    var_key   = "instance_count"
    var_value = "2"
  }
}
`, name, updateTemplateInJsonFormat(), basicVariablesInVarsFormat(name))
}

func testDataSourceExecutionPlanPrices_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rfs_execution_plan_prices" "test" {
  depends_on = [huaweicloud_rfs_execution_plan_v2.test]

  stack_name          = huaweicloud_rfs_stack.test.name
  execution_plan_name = huaweicloud_rfs_execution_plan_v2.test.execution_plan_name
  stack_id            = huaweicloud_rfs_stack.test.id
  execution_plan_id   = huaweicloud_rfs_execution_plan_v2.test.execution_plan_id
}
`, testDataSourceExecutionPlanPrices_base(name))
}
