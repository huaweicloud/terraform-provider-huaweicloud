package swr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSharedRepositories_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_shared_repositories.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSWRDomian(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSharedRepositories_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "repositories.#"),
					resource.TestCheckResourceAttrSet(dataSource, "repositories.0.organization"),
					resource.TestCheckResourceAttrSet(dataSource, "repositories.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "repositories.0.domain_name"),

					resource.TestCheckOutput("organization_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("domain_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceSharedRepositories_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_swr_shared_repositories" "test" {
  depends_on = [huaweicloud_swr_repository_sharing.test]
  
  center = "self"
}

locals {
  repositories = data.huaweicloud_swr_shared_repositories.test.repositories
  organization = local.repositories[0].organization
  name         = local.repositories[0].name
  domain_name  = local.repositories[0].domain_name
}

data "huaweicloud_swr_shared_repositories" "filter_by_organization" {
  center       = "self"
  organization = local.organization
}

data "huaweicloud_swr_shared_repositories" "filter_by_name" {
  center = "self"
  name   = local.name
}

data "huaweicloud_swr_shared_repositories" "filter_by_domain_name" {
  center      = "self"
  domain_name = local.domain_name
}

locals {
  list_by_organization = data.huaweicloud_swr_shared_repositories.filter_by_organization.repositories
  list_by_name         = data.huaweicloud_swr_shared_repositories.filter_by_name.repositories
  list_by_domain_name  = data.huaweicloud_swr_shared_repositories.filter_by_domain_name.repositories
}

output "organization_filter_is_useful" {
  value = length(local.list_by_organization) > 0 && alltrue(
    [for v in local.list_by_organization[*].organization : v == local.organization]
  )
}

output "name_filter_is_useful" {
  value = length(local.list_by_name) > 0 && alltrue(
    [for v in local.list_by_name[*].name : v == local.name]
  )
}

output "domain_name_filter_is_useful" {
  value = length(local.list_by_domain_name) > 0 && alltrue(
    [for v in local.list_by_domain_name[*].domain_name : v == local.domain_name]
  )
}
`, testAccSWRRepositorySharing_basic(rName))
}
