package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityUserTokenInfo_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_identity_user_token_info.test"
		dc  = acceptance.InitDataSourceCheck(all)

		userName = acceptance.RandomAccResourceName()
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
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityUserTokenInfo_basic(userName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "methods.0", "password"),
					resource.TestCheckResourceAttr(all, "user.#", "1"),
				),
			},
		},
	})
}

func testAccIdentityUserTokenInfo_basic(userName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_user_token" "test" {
  account_name = "%[2]s"
  user_name    = huaweicloud_identity_user.test.name
  password     = random_string.test.result
}

data "huaweicloud_identity_user_token_info" "test" {
  token      = huaweicloud_identity_user_token.test.token
  no_catalog = "false"
}
`, testAccIdentityUser_basic(userName), acceptance.HW_DOMAIN_NAME)
}
