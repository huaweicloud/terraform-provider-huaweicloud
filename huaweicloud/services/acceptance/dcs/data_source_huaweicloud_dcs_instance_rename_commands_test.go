package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcsInstanceRenameCommands_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dcs_instance_rename_commands.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDcsInstanceRenameCommands_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "rename_commands.0"),
				),
			},
		},
	})
}

func testAccDataSourceDcsInstanceRenameCommands_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dcs_instance_rename_commands" "test" {
  instance_id = huaweicloud_dcs_instance.test.id
}
`, testAccDcsV1Instance_basic(name))
}
