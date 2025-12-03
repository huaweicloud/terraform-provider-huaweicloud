package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccUserAction_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		password   = acceptance.RandomPassword()
		rName      = "huaweicloud_workspace_user_action.test"
		baseConfig = testAccUserAction_base(name, password)
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
				Config: testAccUserAction_basic_step1(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "op_type", "LOCK"),
					resource.TestCheckResourceAttrSet(rName, "user_id"),
				),
			},
			{
				Config: testAccUserAction_basic_step2(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "op_type", "UNLOCK"),
					resource.TestCheckResourceAttrSet(rName, "user_id"),
				),
			},
			{
				Config: testAccUserAction_basic_step3(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "op_type", "RESET_PWD"),
					resource.TestCheckResourceAttrSet(rName, "user_id"),
				),
			},
		},
	})
}

func testAccUserAction_base(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_user" "test" {
  name        = "%[1]s"
  description = "Created by terraform script"
  active_type = "ADMIN_ACTIVATE"
  email       = "terraform@test.com"
  password    = "%[2]s"

  account_expires            = "0"
  password_never_expires     = false
  enable_change_password     = true
  next_login_change_password = true
  disabled                   = false
}`, name, password)
}

func testAccUserAction_basic_step1(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_user_action" "test" {
  user_id = huaweicloud_workspace_user.test.id
  op_type = "LOCK"

  enable_force_new = "true"
}
`, baseConfig)
}

func testAccUserAction_basic_step2(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_user_action" "test" {
  user_id = huaweicloud_workspace_user.test.id
  op_type = "UNLOCK"

  enable_force_new = "true"
}
`, baseConfig)
}

func testAccUserAction_basic_step3(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_user_action" "test" {
  user_id = huaweicloud_workspace_user.test.id
  op_type = "RESET_PWD"

  enable_force_new = "true"
}
`, baseConfig)
}
