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
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "hashicorp/random",
				// The version of the random provider must be greater than 3.3.0 to support the 'numeric' parameter.
				VersionConstraint: "3.3.0",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataCatalog_basic(userName),
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

func testAccDataCatalog_base(name string) string {
	return fmt.Sprintf(`
resource "random_string" "test" {
  length           = 10
  min_numeric      = 1
  min_special      = 1
  min_lower        = 1
  override_special = "@!"
}

resource "huaweicloud_identity_user" "test" {
  name        = "%[1]s"
  password    = random_string.test.result
  enabled     = true
  email       = "%[1]s@abc.com"

  login_protect_verification_method = "email"
}

resource "huaweicloud_identity_user_token" "test" {
  user_name    = huaweicloud_identity_user.test.name
  password     = random_string.test.result
  account_name = "%[2]s"
  project_name = "%[3]s"

  # Waiting for the user to be created.
  depends_on = [huaweicloud_identity_user.test]
}
`, name, acceptance.HW_DOMAIN_NAME, acceptance.HW_REGION_NAME)
}

func testAccDataCatalog_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_catalog" "all" {
  project_token = huaweicloud_identity_user_token.test.token

  # The token is created asynchronously, so we need to depend on it.
  depends_on = [huaweicloud_identity_user_token.test]
}
`, testAccDataCatalog_base(name))
}
