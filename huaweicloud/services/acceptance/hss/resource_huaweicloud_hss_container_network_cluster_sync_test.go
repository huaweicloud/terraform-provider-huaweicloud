package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccContainerNetworkClusterSync_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testContainerNetworkClusterSync_basic(),
			},
		},
	})
}

func testContainerNetworkClusterSync_basic() string {
	return `
resource "huaweicloud_hss_container_network_cluster_sync" "test" {
  enterprise_project_id = "0"
}
`
}
