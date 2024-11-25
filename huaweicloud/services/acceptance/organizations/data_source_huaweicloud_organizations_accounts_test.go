package organizations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceAccounts_basic(t *testing.T) {
	rName := "data.huaweicloud_organizations_accounts.all"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceAccounts_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "accounts.#"),
					resource.TestCheckResourceAttrSet(rName, "accounts.0.id"),
					resource.TestCheckResourceAttrSet(rName, "accounts.0.name"),
					resource.TestCheckResourceAttrSet(rName, "accounts.0.urn"),
					resource.TestCheckResourceAttrSet(rName, "accounts.0.description"),
					resource.TestCheckResourceAttrSet(rName, "accounts.0.status"),
					resource.TestCheckResourceAttrSet(rName, "accounts.0.join_method"),
					resource.TestCheckResourceAttrSet(rName, "accounts.0.joined_at"),
					resource.TestCheckOutput("parent_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func TestAccDatasourceAccounts_name(t *testing.T) {
	rName := "data.huaweicloud_organizations_accounts.name_filter"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsAccountName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceAccounts_name(acceptance.HW_ORGANIZATIONS_ACCOUNT_NAME),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "accounts.#"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceAccounts_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_organizations_accounts" "all" {}

data "huaweicloud_organizations_accounts" "parent_id_filter" {
  parent_id = data.huaweicloud_organizations_organization.test.root_id
}

output "parent_id_filter_is_useful" {
  value = length(data.huaweicloud_organizations_accounts.parent_id_filter.accounts) > 0
}
`, testAccDatasourceOrganization_basic())
}

func testAccDatasourceAccounts_name(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_organizations_accounts" "name_filter" {
  name = "%[1]s"
}

locals {
  name_filter_result = [for v in data.huaweicloud_organizations_accounts.name_filter.accounts[*].name : v == "%[1]s"]
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}
`, name)
}
