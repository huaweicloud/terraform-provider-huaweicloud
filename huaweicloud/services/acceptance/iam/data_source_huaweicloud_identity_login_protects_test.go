package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityLoginProtects_basic(t *testing.T) {
	userName := acceptance.RandomAccResourceName()
	initPassword := acceptance.RandomPassword()
	dataSourceName := "data.huaweicloud_identity_login_protects.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testTestDataSourceIdentityLoginProtects,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "login_protects.#"),
				),
			},
			{
				Config: testTestDataSourceIdentityLoginProtectsWithUserId(userName, initPassword),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "login_protects.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "login_protects.0.user_id"),
					resource.TestCheckResourceAttr(dataSourceName, "login_protects.0.enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "login_protects.0.verification_method", "email"),
				),
			},
		},
	})
}

const testTestDataSourceIdentityLoginProtects = `
data "huaweicloud_identity_login_protects" "test" {}
`

func testTestDataSourceIdentityLoginProtectsWithUserId(name, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_login_protects" "test" {
  user_id = huaweicloud_identity_user.user_1.id
}
`, testAccIdentityUser_basic(name, password))
}
