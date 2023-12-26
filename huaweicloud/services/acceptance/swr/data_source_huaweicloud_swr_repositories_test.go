package swr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRepositories_basic(t *testing.T) {
	var (
		rName          = acceptance.RandomAccResourceName()
		dataSourceName = "data.huaweicloud_swr_repositories.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byName         = "data.huaweicloud_swr_repositories.filter_by_name"
		nameNotFound   = "data.huaweicloud_swr_repositories.filter_by_name_not_found"
		dcByName       = acceptance.InitDataSourceCheck(byName)
		dcNameNotFound = acceptance.InitDataSourceCheck(nameNotFound)

		byOrganization   = "data.huaweicloud_swr_repositories.filter_by_organization"
		dcByOrganization = acceptance.InitDataSourceCheck(byOrganization)

		byCategory   = "data.huaweicloud_swr_repositories.filter_by_category"
		dcByCategory = acceptance.InitDataSourceCheck(byCategory)

		byIsPublic   = "data.huaweicloud_swr_repositories.filter_by_is_public"
		dcByIsPublic = acceptance.InitDataSourceCheck(byIsPublic)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRepositories_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),

					dcNameNotFound.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_not_found", "true"),

					dcByOrganization.CheckResourceExists(),
					resource.TestCheckOutput("organization_filter_is_useful", "true"),

					dcByCategory.CheckResourceExists(),
					resource.TestCheckOutput("category_filter_is_useful", "true"),

					dcByIsPublic.CheckResourceExists(),
					resource.TestCheckOutput("is_public_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceRepositories_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_swr_repositories" "test" {}

locals {
  name = data.huaweicloud_swr_repositories.test.repositories[0].name
}
data "huaweicloud_swr_repositories" "filter_by_name" {
  name = local.name
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_swr_repositories.filter_by_name.repositories) > 0 && alltrue(
	[for v in data.huaweicloud_swr_repositories.filter_by_name.repositories[*].name : v == local.name]
  )
}

data "huaweicloud_swr_repositories" "filter_by_name_not_found" {
  name = "%s"
}
output "name_filter_not_found" {
  value = length(data.huaweicloud_swr_repositories.filter_by_name_not_found.repositories) == 0
}

locals {
  organization = data.huaweicloud_swr_repositories.test.repositories[0].organization
}
data "huaweicloud_swr_repositories" "filter_by_organization" {
  organization = local.organization
}
output "organization_filter_is_useful" {
  value = length(data.huaweicloud_swr_repositories.filter_by_organization.repositories) > 0 && alltrue(
	[for v in data.huaweicloud_swr_repositories.filter_by_organization.repositories[*].organization : v == local.organization]
  )
}

locals {
  category = data.huaweicloud_swr_repositories.test.repositories[0].category
}
data "huaweicloud_swr_repositories" "filter_by_category" {
  category = local.category
}
output "category_filter_is_useful" {
  value = length(data.huaweicloud_swr_repositories.filter_by_category.repositories) > 0 && alltrue(
	[for v in data.huaweicloud_swr_repositories.filter_by_category.repositories[*].category : v == local.category]
  )
}

locals {
  is_public = data.huaweicloud_swr_repositories.test.repositories[0].is_public
}
data "huaweicloud_swr_repositories" "filter_by_is_public" {
  is_public = local.is_public
}
output "is_public_filter_is_useful" {
  value = length(data.huaweicloud_swr_repositories.filter_by_is_public.repositories) > 0 && alltrue(
	[for v in data.huaweicloud_swr_repositories.filter_by_is_public.repositories[*].is_public : v == local.is_public]
  )
}
`, rName)
}
