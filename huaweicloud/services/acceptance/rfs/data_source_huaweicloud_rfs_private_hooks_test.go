package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRfsPrivateHooks_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_rfs_private_hooks.test"
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
				Config: testAccDataSourceRfsPrivateHooks_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "hooks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "hooks.0.hook_id"),
					resource.TestCheckResourceAttrSet(dataSource, "hooks.0.hook_name"),
					resource.TestCheckResourceAttrSet(dataSource, "hooks.0.default_version"),
					resource.TestCheckResourceAttrSet(dataSource, "hooks.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "hooks.0.update_time"),
				),
			},
		},
	})
}

func testAccDataSourceRfsPrivateHooks_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rfs_private_hooks" "test" {
  sort_key = "create_time"
  sort_dir = "desc"
}
`, testAccPrivateModule_basic(name))
}
