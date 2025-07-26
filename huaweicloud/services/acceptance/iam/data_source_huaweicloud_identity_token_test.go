package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityToken_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identity_token.test"
	userName := acceptance.RandomAccResourceName()
	initPassword := acceptance.RandomPassword()

	config := testAccIdentityToken_basic(userName, initPassword)

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
					resource.TestCheckResourceAttrSet(dataSource, "token"),
				),
			},
		},
	})
}

func testAccIdentityToken_basic(userName, initPassword string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_user_token" "test" {
  account_name = "%[2]s"
  user_name    = huaweicloud_identity_user.user_1.name
  password     = "%[3]s"
}
# 获取Token
data "huaweicloud_identity_token" "test" {
  token = huaweicloud_identity_user_token.test.token
}

`, testAccIdentityUser_basic(userName, initPassword), acceptance.HW_DOMAIN_NAME, initPassword)
}
