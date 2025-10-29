package swrenterprise

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseRepositories_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_repositories.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrSignatureWithImageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseRepositories_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "repositories.#"),
					resource.TestCheckResourceAttrSet(dataSource, "repositories.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "repositories.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "repositories.0.namespace_id"),
					resource.TestCheckResourceAttrSet(dataSource, "repositories.0.namespace_name"),
					resource.TestCheckResourceAttrSet(dataSource, "repositories.0.pull_count"),
					resource.TestCheckResourceAttrSet(dataSource, "repositories.0.artifact_count"),
					resource.TestCheckResourceAttrSet(dataSource, "repositories.0.tag_count"),
					resource.TestCheckResourceAttrSet(dataSource, "repositories.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "repositories.0.updated_at"),

					resource.TestCheckOutput("namespace_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseRepositories_basic() string {
	return `
data "huaweicloud_swr_enterprise_instances" "test" {}

data "huaweicloud_swr_enterprise_repositories" "test" {
  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
}

data "huaweicloud_swr_enterprise_repositories" "filter_by_namespace_id" {
  instance_id  = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_id = data.huaweicloud_swr_enterprise_repositories.test.repositories[0].namespace_id
}

output "namespace_id_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_repositories.filter_by_namespace_id.repositories) > 0 && alltrue(
	[for v in data.huaweicloud_swr_enterprise_repositories.filter_by_namespace_id.repositories[*].namespace_id :
	  v == data.huaweicloud_swr_enterprise_repositories.test.repositories[0].namespace_id]
  )
}
`
}
