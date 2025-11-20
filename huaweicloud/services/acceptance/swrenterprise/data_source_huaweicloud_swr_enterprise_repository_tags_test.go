package swrenterprise

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseRepositoryTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_repository_tags.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrSignatureWithImageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseRepositoryTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.namespace_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.repository_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.digest"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.media_type"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.manifest_media_type"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.artifact_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.pull_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.push_time"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseRepositoryTags_basic() string {
	return `
data "huaweicloud_swr_enterprise_instances" "test" {}

data "huaweicloud_swr_enterprise_repositories" "test" {
  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
}

data "huaweicloud_swr_enterprise_repository_tags" "test" {
  instance_id     = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name  = "library"
  repository_name = data.huaweicloud_swr_enterprise_repositories.test.repositories[0].name
}`
}
