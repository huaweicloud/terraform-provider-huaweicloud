package swrenterprise

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseInstanceArtifactAccessories_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_instance_artifact_accessories.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrSignatureWithImageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseInstanceArtifactAccessories_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "accessories.#"),
					resource.TestCheckResourceAttrSet(dataSource, "accessories.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "accessories.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "accessories.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "accessories.0.digest"),
					resource.TestCheckResourceAttrSet(dataSource, "accessories.0.artifact_id"),
					resource.TestCheckResourceAttrSet(dataSource, "accessories.0.subject_artifact_id"),
					resource.TestCheckResourceAttrSet(dataSource, "accessories.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "total"),

					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseInstanceArtifactAccessories_basic() string {
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

data "huaweicloud_swr_enterprise_instance_artifact_accessories" "test" {
  instance_id     = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name  = "library"
  repository_name = data.huaweicloud_swr_enterprise_repositories.test.repositories[0].name
  reference       = data.huaweicloud_swr_enterprise_instance_artifacts.test.artifacts[0].digest
}

data "huaweicloud_swr_enterprise_instance_artifact_accessories" "filter_by_type" {
  instance_id     = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name  = "library"
  repository_name = data.huaweicloud_swr_enterprise_repositories.test.repositories[0].name
  reference       = data.huaweicloud_swr_enterprise_instance_artifacts.test.artifacts[0].digest
  type            = data.huaweicloud_swr_enterprise_instance_artifact_accessories.test.accessories[0].type
}

output "type_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_instance_artifact_accessories.filter_by_type.accessories) > 0 && alltrue(
	[for v in data.huaweicloud_swr_enterprise_instance_artifact_accessories.filter_by_type.accessories[*].type :
	  strcontains(v, data.huaweicloud_swr_enterprise_instance_artifact_accessories.test.accessories[0].type)]
  )
}`
}
