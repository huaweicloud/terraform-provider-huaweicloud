package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePrivateModules_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_rfs_private_modules.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		name       = acceptance.RandomAccResourceName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePrivateModules_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "modules.#"),
					resource.TestCheckResourceAttrSet(dataSource, "modules.0.module_name"),
					resource.TestCheckResourceAttrSet(dataSource, "modules.0.module_id"),
					resource.TestCheckResourceAttrSet(dataSource, "modules.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "modules.0.update_time"),
				),
			},
		},
	})
}

func testAccDataSourcePrivateModules_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rfs_private_modules" "test" {
  sort_key = "create_time"
  sort_dir = "desc"
}
`, testAccPrivateModule_basic(name))
}
