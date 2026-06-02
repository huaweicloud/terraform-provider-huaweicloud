package swrenterprise

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseAllRepositories_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_swr_enterprise_all_repositories.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSwrEnterpriseAllRepositories_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "repositories.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "repositories.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "repositories.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "repositories.0.namespace_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "repositories.0.tag_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "repositories.0.pull_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "repositories.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "repositories.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "repositories.0.namespace_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "repositories.0.namespace_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "repositories.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "repositories.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "repositories.0.resource_urn"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceSwrEnterpriseAllRepositories_basic() string {
	return `
data "huaweicloud_swr_enterprise_all_repositories" "test" {}

locals {
  filterName = data.huaweicloud_swr_enterprise_all_repositories.test.repositories[0].name
}

data "huaweicloud_swr_enterprise_all_repositories" "filter_by_name" {
  name = local.filterName
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_all_repositories.filter_by_name) > 0 && alltrue(
    [for v in data.huaweicloud_swr_enterprise_all_repositories.filter_by_name.repositories[*].name : v == local.filterName]
  )
}
`
}
