package swrenterprise

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSwrEnterpriseArtifactManualScan_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrSignatureWithImageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccSwrEnterpriseArtifactManualScan_basic(),
			},
		},
	})
}

func testAccSwrEnterpriseArtifactManualScan_basic() string {
	return `
data "huaweicloud_swr_enterprise_instances" "test" {}

data "huaweicloud_swr_enterprise_repositories" "test" {
  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
}

data "huaweicloud_swr_enterprise_instance_artifacts" "test" {
  instance_id     = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name  = data.huaweicloud_swr_enterprise_repositories.test.repositories[0].namespace_name
  repository_name = data.huaweicloud_swr_enterprise_repositories.test.repositories[1].name
}

resource "huaweicloud_swr_enterprise_instance_artifact_manual_scan" "test" {
  instance_id     = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name  = data.huaweicloud_swr_enterprise_repositories.test.repositories[0].namespace_name
  repository_name = data.huaweicloud_swr_enterprise_repositories.test.repositories[1].name
  reference       = data.huaweicloud_swr_enterprise_instance_artifacts.test.artifacts[0].digest
}`
}
