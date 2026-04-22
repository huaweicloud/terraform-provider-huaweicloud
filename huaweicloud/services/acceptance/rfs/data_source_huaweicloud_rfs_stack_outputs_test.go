package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceStackOutputs_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_rfs_stack_outputs.test"
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
				Config: testDataSourceStackOutputs_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "outputs.#"),
				),
			},
		},
	})
}

func testDataSourceStackOutputs_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rfs_stack" "test" {
  name        = "%[1]s"
  description = "Create by acc test for stack outputs data source"
}
`, name)
}

func testDataSourceStackOutputs_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rfs_stack_outputs" "test" {
  stack_name = huaweicloud_rfs_stack.test.name
  stack_id   = huaweicloud_rfs_stack.test.id

  depends_on = [huaweicloud_rfs_stack.test]
}
`, testDataSourceStackOutputs_base(name))
}
