package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Please ensure that the user executing the acceptance test has 'admin' permission.
func TestAccV5UserPassword_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		rName = "huaweicloud_identityv5_user_password.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "hashicorp/random",
				// The version of the random provider must be greater than 3.3.0 to support the 'numeric' parameter.
				VersionConstraint: "3.3.0",
			},
		},
		// This is a one-time resource without delete logic.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccV5UserPassword_basic_step1(name),
			},
			{
				Config: testAccV5UserPassword_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(rName, "old_password", "random_string.test.0", "result"),
					resource.TestCheckResourceAttrPair(rName, "new_password", "random_string.test.1", "result"),
				),
			},
		},
	})
}

func testAccV5UserPassword_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "test" {
  name = "%[1]s"
}

resource "huaweicloud_identityv5_access_key" "test" {
  user_id = huaweicloud_identityv5_user.test.id
}

resource "random_string" "test" {
  count = 2

  length           = 10
  min_numeric      = 1
  min_special      = 2
  min_lower        = 1
  override_special = "@!#&"
}

resource "huaweicloud_identityv5_login_profile" "test" {
  user_id  = huaweicloud_identityv5_user.test.id
  password = random_string.test[0].result
}
`, name)
}

// lintignore:AT004
func testAccV5UserPassword_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

provider "huaweicloud" {
  alias  = "new_user"
  region = "%[2]s"

  access_key = huaweicloud_identityv5_access_key.test.id
  secret_key = huaweicloud_identityv5_access_key.test.secret_access_key
}

resource "huaweicloud_identityv5_user_password" "test" {
  provider = huaweicloud.new_user

  new_password = random_string.test[1].result
  old_password = huaweicloud_identityv5_login_profile.test.password

  depends_on = [huaweicloud_identityv5_login_profile.test]
}
`, testAccV5UserPassword_basic_step1(name), acceptance.HW_REGION_NAME)
}
