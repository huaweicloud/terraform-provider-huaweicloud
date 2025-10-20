package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccContainerNetworkPolicySync_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testContainerNetworkPolicySync_basic,
			},
		},
	})
}

const testContainerNetworkPolicySync_basic = `
resource "huaweicloud_hss_container_network_policy_sync" "test" {
  cluster_id            = "non-exist-id"
  enterprise_project_id = "0"
}
`
