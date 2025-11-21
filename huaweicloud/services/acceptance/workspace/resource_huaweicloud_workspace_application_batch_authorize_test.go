package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccApplicationBatchAuthorize_basic(t *testing.T) {
	var (
		name         = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_workspace_application_batch_authorize.test"
		baseConfig   = testAccApplicationBatchAuthorize_base(name)
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
				Config: testAccApplicationBatchAuthorize_basic_step1(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "authorization_type", "ASSIGN_USER"),
					resource.TestCheckResourceAttr(resourceName, "app_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "2"),
				),
			},
			{
				Config: testAccApplicationBatchAuthorize_basic_step2(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "authorization_type", "ASSIGN_USER"),
					resource.TestCheckResourceAttr(resourceName, "app_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "del_users.#", "1"),
				),
			},
			{
				Config: testAccApplicationBatchAuthorize_basic_step3(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "authorization_type", "ALL_USER"),
					resource.TestCheckResourceAttr(resourceName, "app_ids.#", "1"),
				),
			},
		},
	})
}

func testAccApplicationBatchAuthorize_base(name string) string {
	randomPhoneNum := acctest.RandIntRange(1000000000, 1999999999)

	return fmt.Sprintf(`
data "huaweicloud_workspace_application_catalogs" "test" {}

resource "huaweicloud_workspace_application" "test" {
  count = 2

  name               = format("%%s_%%d", "%[1]s", count.index)
  version            = "1.0.0"
  authorization_type = "ALL_USER"
  install_type       = "QUIET_INSTALL"
  support_os         = "Windows"
  catalog_id         = try(data.huaweicloud_workspace_application_catalogs.test.catalogs[0].id, "NOT_FOUND")
  install_command    = "terraform test install"
  description        = "Created by terraform script"

  application_file_store {
    store_type = "LINK"
    file_link  = "https://www.huaweicloud.com/TerraformTest.msi"
  }

  lifecycle {
    ignore_changes = [
      authorization_type
    ]
  }
}

resource "huaweicloud_workspace_user" "test" {
  count = 3

  name  = format("%[1]s_%%d", count.index)
  email = format("test%%d@example.com", count.index)
  phone = format("+%[2]d%%d", count.index)
}
`, name, randomPhoneNum)
}

func testAccApplicationBatchAuthorize_basic_step1(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_application_batch_authorize" "test" {
  depends_on = [
    huaweicloud_workspace_user.test
  ]

  app_ids            = huaweicloud_workspace_application.test[*].id
  authorization_type = "ASSIGN_USER"

  dynamic "users" {
    for_each = slice(huaweicloud_workspace_user.test, 0, 2)

    content {
      account      = users.value.name
      account_type = "SIMPLE"
    }
  }

  enable_force_new = "true"
}
`, baseConfig)
}

func testAccApplicationBatchAuthorize_basic_step2(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_application_batch_authorize" "test" {
  depends_on = [
    huaweicloud_workspace_user.test
  ]

  app_ids            = huaweicloud_workspace_application.test[*].id
  authorization_type = "ASSIGN_USER"

  dynamic "del_users" {
    for_each = slice(huaweicloud_workspace_user.test, 0, 1)

    content {
      account      = del_users.value.name
      account_type = "SIMPLE"
    }
  }

  enable_force_new = "true"
}
`, baseConfig)
}

func testAccApplicationBatchAuthorize_basic_step3(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_application_batch_authorize" "test" {
  depends_on = [
    huaweicloud_workspace_user.test
  ]

  app_ids            = slice(huaweicloud_workspace_application.test[*].id, 0, 1)
  authorization_type = "ALL_USER"

  enable_force_new = "true"
}
`, baseConfig)
}
