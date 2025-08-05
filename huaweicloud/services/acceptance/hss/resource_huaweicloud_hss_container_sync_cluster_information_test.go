package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccContainerSyncClusterInformation_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testContainerSyncClusterInformation_basic,
			},
		},
	})
}

const testContainerSyncClusterInformation_basic = `
resource "huaweicloud_hss_container_sync_cluster_information" "test" {
  enterprise_project_id = "0"
}
`
