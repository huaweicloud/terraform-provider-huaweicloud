package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceStacks_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_rfs_stacks.test"
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
				Config: testDataSourceStacks_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "stacks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "stacks.0.stack_name"),
					resource.TestCheckResourceAttrSet(dataSource, "stacks.0.stack_id"),
					resource.TestCheckResourceAttrSet(dataSource, "stacks.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "stacks.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "stacks.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "stacks.0.update_time"),
				),
			},
		},
	})
}

func testDataSourceStacks_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rfs_stack" "test" {
  name        = "%[1]s"
  description = "Create by acc test for stacks data source"
}
`, name)
}

func testDataSourceStacks_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rfs_stacks" "test" {
  depends_on = [huaweicloud_rfs_stack.test]
}
`, testDataSourceStacks_base(name))
}
