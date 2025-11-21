package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccApplicationVisibilityBatchAction_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_workspace_application_visibility_batch_action.test"
		name         = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationVisibilityBatchAction_basic(name, "enable"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "action", "enable"),
					resource.TestCheckResourceAttr(resourceName, "app_ids.#", "2"),
				),
			},
			{
				Config: testAccApplicationVisibilityBatchAction_basic(name, "disable"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "action", "disable"),
					resource.TestCheckResourceAttr(resourceName, "app_ids.#", "2"),
				),
			},
		},
	})
}

func testAccApplicationVisibilityBatchAction_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_application_catalogs" "test" {}

resource "huaweicloud_workspace_application" "test" {
  count = 2

  name               = "%[1]s-${count.index}"
  version            = "1.0.0"
  description        = "Created by terraform script"
  authorization_type = "ALL_USER"
  install_type       = "QUIET_INSTALL"
  support_os         = "Windows"
  catalog_id         = try(data.huaweicloud_workspace_application_catalogs.test.catalogs[0].id, "NOT_FOUND")

  application_file_store {
    store_type = "LINK"
    file_link  = "https://www.huaweicloud.com/TerraformTest.msi"
  }
}
`, name)
}

func testAccApplicationVisibilityBatchAction_basic(name, action string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_application_visibility_batch_action" "test" {
  action  = "%[2]s"
  app_ids = huaweicloud_workspace_application.test[*].id

  enable_force_new = "true"
}
`, testAccApplicationVisibilityBatchAction_base(name), action)
}
