package organizations

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAccounts_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_organizations_accounts.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byParentId   = "data.huaweicloud_organizations_accounts.filter_by_parent_id"
		dcByParentId = acceptance.InitDataSourceCheck(byParentId)

		byName   = "data.huaweicloud_organizations_accounts.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byWithRegisterContactInfo   = "data.huaweicloud_organizations_accounts.filter_by_with_register_contact_info"
		dcByWithRegisterContactInfo = acceptance.InitDataSourceCheck(byWithRegisterContactInfo)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAccounts_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "accounts.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'parent_id' parameter.
					dcByParentId.CheckResourceExists(),
					resource.TestCheckOutput("is_parent_id_filter_useful", "true"),
					// Filter by 'name' parameter.
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(byName, "accounts.0.id"),
					resource.TestCheckResourceAttrSet(byName, "accounts.0.name"),
					resource.TestCheckResourceAttrSet(byName, "accounts.0.urn"),
					resource.TestCheckResourceAttrSet(byName, "accounts.0.status"),
					resource.TestCheckResourceAttrSet(byName, "accounts.0.join_method"),
					resource.TestCheckResourceAttrSet(byName, "accounts.0.joined_at"),
					resource.TestCheckResourceAttrSet(byName, "accounts.0.description"),
					// Filter by 'with_register_contact_info' parameter.
					dcByWithRegisterContactInfo.CheckResourceExists(),
					resource.TestCheckOutput("is_with_register_contact_info_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataAccounts_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_organizations_organization" "test" {}

resource "huaweicloud_organizations_account" "test" {
  name        = "%[1]s"
  parent_id   = data.huaweicloud_organizations_organization.test.root_id
  phone       = "13245678978"
  description = "Created by terraform script"
}
`, name)
}

func testAccDataAccounts_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameters.
data "huaweicloud_organizations_accounts" "test" {}

# Filter by 'parent_id' parameter.
data "huaweicloud_organizations_accounts" "filter_by_parent_id" {
  parent_id = data.huaweicloud_organizations_organization.test.root_id
}

output "is_parent_id_filter_useful" {
  value = length(data.huaweicloud_organizations_accounts.filter_by_parent_id.accounts) > 0
}

# Filter by 'name' parameter.
locals {
  name = huaweicloud_organizations_account.test.name
}

data "huaweicloud_organizations_accounts" "filter_by_name" {
  name = local.name

  depends_on = [huaweicloud_organizations_account.test]
}

locals {
  name_filter_result = [for v in data.huaweicloud_organizations_accounts.filter_by_name.accounts[*].name : v == local.name]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by 'with_register_contact_info' parameter.
data "huaweicloud_organizations_accounts" "filter_by_with_register_contact_info" {
  parent_id                  = data.huaweicloud_organizations_organization.test.root_id
  with_register_contact_info = true
}

locals {
  with_register_contact_info_filter_result = [for v in data.huaweicloud_organizations_accounts.filter_by_with_register_contact_info.accounts :
  v.email != "" || v.mobile_phone != ""]
}

output "is_with_register_contact_info_filter_useful" {
  value = length(local.with_register_contact_info_filter_result) > 0 && alltrue(local.with_register_contact_info_filter_result)
}
`, testAccDataAccounts_base(name))
}
