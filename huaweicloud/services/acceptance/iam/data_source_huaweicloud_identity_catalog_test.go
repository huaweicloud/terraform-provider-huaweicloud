package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataCatalog_basic(t *testing.T) {
	var (
		userName = acceptance.RandomAccResourceName()
		password = acceptance.RandomPassword()

		all = "data.huaweicloud_identity_catalog.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPrecheckDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataCatalog_basic(userName, password),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "catalog.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "catalog.0.id"),
					resource.TestCheckResourceAttrSet(all, "catalog.0.name"),
					resource.TestCheckResourceAttrSet(all, "catalog.0.type"),
				),
			},
		},
	})
}

func testAccDataCatalog_base(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user" "test" {
  name        = "%[1]s"
  password    = "%[2]s"
  enabled     = true
  email       = "%[1]s@abc.com"
  description = "tested by terraform"

  login_protect_verification_method = "email"
}

resource "huaweicloud_identity_user_token" "test" {
  user_name    = huaweicloud_identity_user.test.name
  password     = "%[2]s"
  account_name = "%[3]s"
  project_name = "%[4]s"
}
`, name, password, acceptance.HW_DOMAIN_NAME, acceptance.HW_REGION_NAME)
}

func testAccDataCatalog_basic(name, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_catalog" "all" {
  project_token = huaweicloud_identity_user_token.test.token
}
`, testAccDataCatalog_base(name, password))
}
