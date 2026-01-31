package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityLoginProtects_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_identity_login_protects.test"
		dc     = acceptance.InitDataSourceCheck(dcName)

		rName    = acceptance.RandomAccResourceName()
		password = acceptance.RandomPassword()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testIdentityLoginProtects,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dcName, "login_protects.#"),
				),
			},
			{
				Config: testIdentityLoginProtectsWithUserId(rName, password),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dcName, "login_protects.#"),
					resource.TestCheckResourceAttrSet(dcName, "login_protects.0.user_id"),
					resource.TestCheckResourceAttr(dcName, "login_protects.0.enabled", "true"),
					resource.TestCheckResourceAttr(dcName, "login_protects.0.verification_method", "email"),
				),
			},
		},
	})
}

const testIdentityLoginProtects = `
data "huaweicloud_identity_login_protects" "test" {}
`

func testIdentityLoginProtectsWithUserId(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_login_protects" "test" {
  user_id = huaweicloud_identity_user.test.id
}
`, testAccIdentityUser_basic(rName, password))
}
