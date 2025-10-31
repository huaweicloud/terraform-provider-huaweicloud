package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityUserTokenInfo_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identity_user_token_info.test"
	userName := acceptance.RandomAccResourceName()
	initPassword := acceptance.RandomPassword()

	config := testAccIdentityUserTokenInfo_basic(userName, initPassword)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSource, "methods.0", "password"),
					resource.TestCheckResourceAttr(dataSource, "user.#", "1"),
				),
			},
		},
	})
}

func testAccIdentityUserTokenInfo_basic(userName, initPassword string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_user_token" "test" {
  account_name = "%[2]s"
  user_name    = huaweicloud_identity_user.user_1.name
  password     = "%[3]s"
}
data "huaweicloud_identity_user_token_info" "test" {
  token = huaweicloud_identity_user_token.test.token
  no_catalog = "false"
}

`, testAccIdentityUser_basic(userName, initPassword), acceptance.HW_DOMAIN_NAME, initPassword)
}
