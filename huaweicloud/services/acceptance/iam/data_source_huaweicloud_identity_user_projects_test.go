package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityUserProjects_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identity_user_projects.test"
	userName := acceptance.RandomAccResourceName()
	password := acceptance.RandomPassword()
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testTestDataSourceIdentityUserProjects1,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "projects.#"),
				),
			},
			{
				Config: testTestDataSourceIdentityUserProjects2(userName, password),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "projects.#"),
				),
			},
		},
	})
}

const testTestDataSourceIdentityUserProjects1 = `
data "huaweicloud_identity_user_projects" "test" {}
`

func testTestDataSourceIdentityUserProjects2(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user" "user_1" {
  name                              = "%[1]s"
  password                          = "%[2]s"
  enabled                           = true
  email                             = "%[1]s@abc.com"
  description                       = "tested by terraform"
  login_protect_verification_method = "email"
}

data "huaweicloud_identity_user_projects" "test" {
  user_id = huaweicloud_identity_user.user_1.id
}
`, name, password)
}
