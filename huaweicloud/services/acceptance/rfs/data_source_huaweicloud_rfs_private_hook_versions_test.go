package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePrivateHookVersions_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_rfs_private_hook_versions.test"
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
				Config: testDataSourcePrivateHookVersions_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "versions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.hook_name"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.hook_id"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.hook_version"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.create_time"),
				),
			},
		},
	})
}

func testDataSourcePrivateHookVersions_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rfs_private_hook" "test" {
  name                = "%[1]s"
  version             = "1.0.0"
  version_description = "acc test hook version"
  policy_body         = <<EOT
package policy

import rego.v1

hook_result := {
  "is_passed": true,
  "err_msg":   "",
}
EOT

  configuration {
    failure_mode  = "WARN"
    target_stacks = "ALL"
  }
}
`, name)
}

func testDataSourcePrivateHookVersions_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rfs_private_hook_versions" "test" {
  hook_name = huaweicloud_rfs_private_hook.test.name
  sort_key  = "create_time"
  sort_dir  = "desc"
}
`, testDataSourcePrivateHookVersions_base(name))
}
