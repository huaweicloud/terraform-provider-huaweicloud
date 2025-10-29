package swrenterprise

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseNamespaceRepositories_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_namespace_repositories.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrSignatureWithImageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseNamespaceRepositories_basic(),
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
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseNamespaceRepositories_basic() string {
	return `
data "huaweicloud_swr_enterprise_instances" "test" {}

data "huaweicloud_swr_enterprise_namespace_repositories" "test" {
  instance_id    = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name = "library"
}
`
}
