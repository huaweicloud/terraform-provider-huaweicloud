package swrenterprise

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSwrEnterpriseRepositoryDelete_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrSignatureWithImageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccSwrEnterpriseRepositoryDelete_basic(),
			},
		},
	})
}

func testAccSwrEnterpriseRepositoryDelete_basic() string {
	return `
data "huaweicloud_swr_enterprise_instances" "test" {}

data "huaweicloud_swr_enterprise_repositories" "test" {
  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
}

resource "huaweicloud_swr_enterprise_repository_delete" "test" {
  instance_id     = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name  = try(data.huaweicloud_swr_enterprise_repositories.test.repositories[0].namespace_name, "")
  repository_name = try(data.huaweicloud_swr_enterprise_repositories.test.repositories[0].name, "")

  lifecycle {
    ignore_changes = [
      namespace_name, repository_name,
    ]
  }
}`
}
