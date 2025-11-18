package swrenterprise

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseInstanceArtifacts_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_instance_artifacts.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrSignatureWithImageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseInstanceArtifacts_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.#"),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.0.namespace_id"),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.0.repository_id"),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.0.repository_name"),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.0.digest"),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.0.media_type"),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.0.manifest_media_type"),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.0.pull_time"),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.0.push_time"),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.0.tags.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.0.tags.0.repository_id"),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.0.tags.0.artifact_id"),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.0.tags.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.0.tags.0.push_time"),
					resource.TestCheckResourceAttrSet(dataSource, "artifacts.0.tags.0.pull_time"),
					resource.TestCheckResourceAttrSet(dataSource, "total"),

					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("tags_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseInstanceArtifacts_basic() string {
	return `
data "huaweicloud_swr_enterprise_instances" "test" {}

data "huaweicloud_swr_enterprise_repositories" "test" {
  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
}

data "huaweicloud_swr_enterprise_instance_artifacts" "test" {
  instance_id     = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name  = "library"
  repository_name = data.huaweicloud_swr_enterprise_repositories.test.repositories[0].name
}

data "huaweicloud_swr_enterprise_instance_artifacts" "filter_by_type" {
  instance_id     = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name  = "library"
  repository_name = data.huaweicloud_swr_enterprise_repositories.test.repositories[0].name
  type            = data.huaweicloud_swr_enterprise_instance_artifacts.test.artifacts[0].type
}

output "type_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_instance_artifacts.filter_by_type.artifacts) > 0 && alltrue(
	[for v in data.huaweicloud_swr_enterprise_instance_artifacts.filter_by_type.artifacts[*].type :
	  strcontains(v, data.huaweicloud_swr_enterprise_instance_artifacts.test.artifacts[0].type)]
  )
}

data "huaweicloud_swr_enterprise_instance_artifacts" "filter_by_tags" {
  instance_id     = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name  = "library"
  repository_name = data.huaweicloud_swr_enterprise_repositories.test.repositories[0].name
  tags            = data.huaweicloud_swr_enterprise_instance_artifacts.test.artifacts[0].tags[0].name
}

output "tags_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_instance_artifacts.filter_by_tags.artifacts) > 0
}
`
}
