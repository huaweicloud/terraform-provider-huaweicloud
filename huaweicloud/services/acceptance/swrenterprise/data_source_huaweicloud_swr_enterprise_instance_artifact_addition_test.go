package swrenterprise

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseInstanceArtifactAddition_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_instance_artifact_addition.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseInstanceArtifactAddition_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "build_histories.#"),
					resource.TestCheckResourceAttrSet(dataSource, "build_histories.0.media_type"),
					resource.TestCheckResourceAttrSet(dataSource, "build_histories.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "build_histories.0.digest"),
					resource.TestCheckResourceAttrSet(dataSource, "build_histories.0.created_by"),
					resource.TestCheckResourceAttrSet(dataSource, "build_histories.0.empty_layer"),
					resource.TestCheckResourceAttrSet(dataSource, "build_histories.0.created"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseInstanceArtifactAddition_basic() string {
	return `
data "huaweicloud_swr_enterprise_instances" "test" {}

data "huaweicloud_swr_enterprise_repositories" "test" {
  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
}

data "huaweicloud_swr_enterprise_instance_artifacts" "test" {
  instance_id     = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name  = data.huaweicloud_swr_enterprise_repositories.test.repositories[0].namespace_name
  repository_name = data.huaweicloud_swr_enterprise_repositories.test.repositories[0].name
}

data "huaweicloud_swr_enterprise_instance_artifact_addition" "test" {
  instance_id     = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name  = data.huaweicloud_swr_enterprise_repositories.test.repositories[0].namespace_name
  repository_name = data.huaweicloud_swr_enterprise_repositories.test.repositories[0].name
  reference       = data.huaweicloud_swr_enterprise_instance_artifacts.test.artifacts[0].digest
  addition        = "build_history"
}
`
}
