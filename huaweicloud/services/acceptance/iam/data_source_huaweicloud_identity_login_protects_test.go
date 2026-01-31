package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataLoginProtects_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_identity_login_protects.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byUserId   = "data.huaweicloud_identity_login_protects.filter_by_user_id"
		dcByUserId = acceptance.InitDataSourceCheck(byUserId)
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
				Config: testAccDataLoginProtects_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "login_protects.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByUserId.CheckResourceExists(),
					resource.TestCheckOutput("is_filter_by_user_id_useful", "true"),
				),
			},
		},
	})
}

func testAccDataLoginProtects_base(name string) string {
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
  description = "tested by terraform"

  login_protect_verification_method = "email"
}
`, name)
}

func testAccDataLoginProtects_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# All
data "huaweicloud_identity_login_protects" "all" {
  # Waiting for the login protect to be created
  depends_on = [huaweicloud_identity_user.test]
}

# Filter by user ID
locals {
  user_id = huaweicloud_identity_user.test.id
}

data "huaweicloud_identity_login_protects" "filter_by_user_id" {
  user_id = local.user_id
}

locals {
  login_protects = data.huaweicloud_identity_login_protects.filter_by_user_id.login_protects
}

output "is_filter_by_user_id_useful" {
  value = length(local.login_protects) > 0 && alltrue(
    [for v in local.login_protects[*].user_id : v == local.user_id]
  )
}
`, testAccDataLoginProtects_base(name))
}
