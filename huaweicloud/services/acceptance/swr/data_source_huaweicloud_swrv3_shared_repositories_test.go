package swr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSwrv3SourceSharedRepositories_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swrv3_shared_repositories.test"
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
				Config: testAccDataSourceSwrv3SharedRepositories_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "repos.#"),
					resource.TestCheckResourceAttrSet(dataSource, "repos.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "repos.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "repos.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "repos.0.organization"),
					resource.TestCheckResourceAttrSet(dataSource, "repos.0.num_download"),
					resource.TestCheckResourceAttrSet(dataSource, "repos.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "repos.0.category"),
					resource.TestCheckResourceAttrSet(dataSource, "repos.0.is_public"),
					resource.TestCheckResourceAttrSet(dataSource, "repos.0.num_images"),
					resource.TestCheckResourceAttrSet(dataSource, "repos.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "repos.0.updated_at"),

					resource.TestCheckOutput("organization_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceSwrv3SharedRepositories_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_swrv3_shared_repositories" "test" {
  depends_on = [huaweicloud_swr_repository_sharing.test]
  
  shared_by = "self"
}

locals {
  repositories = data.huaweicloud_swrv3_shared_repositories.test.repos
  organization = local.repositories[0].organization
  name         = local.repositories[0].name
}

data "huaweicloud_swrv3_shared_repositories" "filter_by_organization" {
  shared_by    = "self"
  organization = local.organization
}

data "huaweicloud_swrv3_shared_repositories" "filter_by_name" {
  shared_by = "self"
  name      = local.name
}

locals {
  list_by_organization = data.huaweicloud_swrv3_shared_repositories.filter_by_organization.repos
  list_by_name         = data.huaweicloud_swrv3_shared_repositories.filter_by_name.repos
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
`, testAccSWRRepositorySharing_basic(rName))
}
