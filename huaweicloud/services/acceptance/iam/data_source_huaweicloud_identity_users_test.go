package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityUsersDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identity_users.test"
	userName := acceptance.RandomAccResourceName()
	initPassword := acceptance.RandomPassword()
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUsersDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.name"),
				),
			},
			{
				Config: testAccUsersDataSourceBasic1(userName, initPassword),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "users.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "users.0.name", userName),
					resource.TestCheckResourceAttr(dataSourceName, "users.0.description", "tested by terraform"),
					resource.TestCheckResourceAttr(dataSourceName, "users.0.enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "users.0.password_status", "true"),
				),
			},
			{
				Config: testAccUsersDataSourceBasic2(userName, initPassword),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "users.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "users.0.name", userName),
					resource.TestCheckResourceAttr(dataSourceName, "users.0.description", "tested by terraform"),
					resource.TestCheckResourceAttr(dataSourceName, "users.0.enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "users.0.password_status", "true"),
				),
			},
		},
	})
}

const testAccUsersDataSourceBasic string = `
data "huaweicloud_identity_users" "test" {
}
`

func testAccUsersDataSourceBasic1(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user" "user_1" {
  name                              = "%[1]s"
  password                          = "%[2]s"
  enabled                           = true
  email                             = "%[1]s@abc.com"
  description                       = "tested by terraform"
  login_protect_verification_method = "email"
}

data "huaweicloud_identity_users" "test" {
  name = huaweicloud_identity_user.user_1.name
}

`, name, password)
}

func testAccUsersDataSourceBasic2(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user" "user_1" {
  name                              = "%[1]s"
  password                          = "%[2]s"
  enabled                           = true
  email                             = "%[1]s@abc.com"
  description                       = "tested by terraform"
  login_protect_verification_method = "email"
}

data "huaweicloud_identity_users" "test" {
  user_id = huaweicloud_identity_user.user_1.id
}

`, name, password)
}
