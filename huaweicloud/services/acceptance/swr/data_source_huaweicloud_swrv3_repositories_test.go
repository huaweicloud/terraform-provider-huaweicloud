package swr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrv3Repositories_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swrv3_repositories.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrv3Repositories_basic(rName),
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

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("organization_filter_is_useful", "true"),
					resource.TestCheckOutput("category_filter_is_useful", "true"),
					resource.TestCheckOutput("is_public_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSwrv3Repositories_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_swrv3_repositories" "test" {
  depends_on = [huaweicloud_swr_repository.test]
}

data "huaweicloud_swrv3_repositories" "filter_by_name" {
  name = huaweicloud_swr_repository.test.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_swrv3_repositories.filter_by_name.repos) > 0 && alltrue(
	[for v in data.huaweicloud_swrv3_repositories.filter_by_name.repos[*].name : v == huaweicloud_swr_repository.test.name]
  )
}

data "huaweicloud_swrv3_repositories" "filter_by_organization" {
  organization = huaweicloud_swr_repository.test.organization
}

output "organization_filter_is_useful" {
  value = length(data.huaweicloud_swrv3_repositories.filter_by_organization.repos) > 0 && alltrue(
	[for v in data.huaweicloud_swrv3_repositories.filter_by_organization.repos[*].organization : v == huaweicloud_swr_repository.test.organization]
  )
}

data "huaweicloud_swrv3_repositories" "filter_by_category" {
  category = huaweicloud_swr_repository.test.category
}

output "category_filter_is_useful" {
  value = length(data.huaweicloud_swrv3_repositories.filter_by_category.repos) > 0 && alltrue(
	[for v in data.huaweicloud_swrv3_repositories.filter_by_category.repos[*].category : v == huaweicloud_swr_repository.test.category]
  )
}

data "huaweicloud_swrv3_repositories" "filter_by_is_public" {
  is_public = huaweicloud_swr_repository.test.is_public
}

output "is_public_filter_is_useful" {
  value = length(data.huaweicloud_swrv3_repositories.filter_by_is_public.repos) > 0 && alltrue(
	[for v in data.huaweicloud_swrv3_repositories.filter_by_is_public.repos[*].is_public : v == huaweicloud_swr_repository.test.is_public]
  )
}
`, testAccSWRRepository_basic(name))
}
