package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccExecutionPlan_basic(t *testing.T) {
	name := acceptance.RandomAccResourceNameWithDash()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccExecutionPlan_basic(name),
			},
		},
	})
}

func testAccExecutionPlan_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rfs_stack" "test" {
  name = "%[1]s"
}

resource "huaweicloud_rfs_execution_plan" "test" {
  stack_name    = huaweicloud_rfs_stack.test.name
  name          = "%[1]s"
  stack_id      = huaweicloud_rfs_stack.test.id
  description   = "Created plan by script"
  template_body = %[2]s
  vars_body     = %[3]s
}
`, name, updateTemplateInJsonFormat(), basicVariablesInVarsFormat(name))
}

func TestAccExecutionPlan_withUri(t *testing.T) {
	name := acceptance.RandomAccResourceNameWithDash()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccExecutionPlan_withUri(name),
			},
		},
	})
}

func testAccExecutionPlan_withUri(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rfs_execution_plan" "test" {
  stack_name   = huaweicloud_rfs_stack.test.name
  name         = "%[2]s"
  stack_id     = huaweicloud_rfs_stack.test.id
  description  = "Created plan by script"
  template_uri = huaweicloud_rfs_stack.test.template_uri
  vars_uri     = huaweicloud_rfs_stack.test.vars_uri
}
`, testAccStack_withUri_hclBody(name, updateTemplateInHclFormat(), basicVariablesInVarsFormat(name)), name)
}
