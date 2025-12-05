package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataApplicationAuthorizations_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_workspace_application_authorizations.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_workspace_application_authorizations.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byTargetType   = "data.huaweicloud_workspace_application_authorizations.filter_by_target_type"
		dcByTargetType = acceptance.InitDataSourceCheck(byTargetType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataApplicationAuthorizations_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "authorizations.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(all, "authorizations.0.account"),
					resource.TestCheckResourceAttrSet(all, "authorizations.0.account_type"),
					// Filter by 'name' parameter.
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					// Filter by 'target_type' parameter.
					dcByTargetType.CheckResourceExists(),
					resource.TestCheckOutput("is_target_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataApplicationAuthorizations_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_application_catalogs" "test" {}

resource "huaweicloud_workspace_application" "test" {
  name               = "%[1]s"
  version            = "1.0.0"
  authorization_type = "ALL_USER"
  install_type       = "QUIET_INSTALL"
  support_os         = "Windows"
  catalog_id         = try(data.huaweicloud_workspace_application_catalogs.test.catalogs[0].id, "NOT_FOUND")
  install_command    = "terraform test install"
  description        = "Created by terraform script"

  application_file_store {
    store_type = "LINK"
    file_link  = "https://www.huaweicloud.com/TerraformTest.msi"
  }

  lifecycle {
    ignore_changes = [
      authorization_type
    ]
  }
}

resource "huaweicloud_workspace_user" "test" {
  name  = "%[1]s_user"
  email = "test@example.com"
  phone = "+8613800000000"
}

resource "huaweicloud_workspace_application_batch_authorize" "test" {
  depends_on = [
    huaweicloud_workspace_user.test
  ]

  app_ids            = [huaweicloud_workspace_application.test.id]
  authorization_type = "ASSIGN_USER"

  users {
    account       = huaweicloud_workspace_user.test.name
    account_type  = "SIMPLE"
    platform_type = "LOCAL"
  }

  enable_force_new = "true"
}
`, name)
}

func testAccDataApplicationAuthorizations_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_workspace_application_authorizations" "all" {
  app_id = huaweicloud_workspace_application.test.id

  depends_on = [
    huaweicloud_workspace_application_batch_authorize.test
  ]
}

locals {
  authorized_account      = huaweicloud_workspace_user.test.name
  authorized_account_type = "SIMPLE"
}

# Without any filter parameter.
data "huaweicloud_workspace_application_authorizations" "filter_by_name" {
  app_id = huaweicloud_workspace_application.test.id
  name   = local.authorized_account

  depends_on = [
    huaweicloud_workspace_application_batch_authorize.test
  ]
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_workspace_application_authorizations.filter_by_name.authorizations[*].account :
    strcontains(v, local.authorized_account)
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by 'target_type' parameter.
data "huaweicloud_workspace_application_authorizations" "filter_by_target_type" {
  app_id     = huaweicloud_workspace_application.test.id
  target_type = local.authorized_account_type

  depends_on = [
    huaweicloud_workspace_application_batch_authorize.test
  ]
}

locals {
  target_type_filter_result = [
    for v in data.huaweicloud_workspace_application_authorizations.filter_by_target_type.authorizations[*].account_type :
    v == local.authorized_account_type
  ]
}

output "is_target_type_filter_useful" {
  value = length(local.target_type_filter_result) > 0 && alltrue(local.target_type_filter_result)
}
`, testAccDataApplicationAuthorizations_base(name))
}
