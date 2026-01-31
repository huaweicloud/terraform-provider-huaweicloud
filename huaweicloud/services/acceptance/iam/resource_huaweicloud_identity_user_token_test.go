package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityUserToken_basic(t *testing.T) {
	userName := acceptance.RandomAccResourceName()
	initPassword := acceptance.RandomPassword()
	resourceName := "huaweicloud_identity_user_token.test"

	// Avoid CheckDestroy because the token can not be destroyed.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityUserToken_basic(userName, initPassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "token"),
					resource.TestCheckResourceAttrSet(resourceName, "expires_at"),
				),
			},
		},
	})
}

func testAccIdentityUserToken_basic(userName, initPassword string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_user_token" "test" {
  account_name = "%[2]s"
  user_name    = huaweicloud_identity_user.test.name
  password     = "%[3]s"
}
`, testAccIdentityUser_basic(userName, initPassword), acceptance.HW_DOMAIN_NAME, initPassword)
}

func TestAccIdentityUserToken_project(t *testing.T) {
	userName := acceptance.RandomAccResourceName()
	initPassword := acceptance.RandomPassword()
	resourceName := "huaweicloud_identity_user_token.test"

	// Avoid CheckDestroy because the token can not be destroyed.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityUserToken_project(userName, initPassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "token"),
					resource.TestCheckResourceAttrSet(resourceName, "expires_at"),
				),
			},
		},
	})
}

func testAccIdentityUserToken_project(userName, initPassword string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_user_token" "test" {
  account_name = "%[2]s"
  user_name    = huaweicloud_identity_user.test.name
  password     = "%[3]s"
  project_name = "%[4]s"
}
`, testAccIdentityUser_basic(userName, initPassword), acceptance.HW_DOMAIN_NAME, initPassword, acceptance.HW_REGION_NAME)
}
