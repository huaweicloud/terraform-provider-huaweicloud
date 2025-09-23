package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// The probability of successfully creating and enabling CCE cluster protection using the API is extremely low, so the
// test case only tests error protection scenarios.
func TestAccHSSContainerKubernetesClusterProtectionEnable_basic(t *testing.T) {
	rName := "huaweicloud_hss_container_kubernetes_cluster_protection_enable.test"

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testContainerKubernetesClusterProtectionEnable_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "id"),
					resource.TestCheckResourceAttrSet(rName, "protect_status"),
				),
			},
		},
	})
}

const testContainerKubernetesClusterProtectionEnable_basic = `
resource "huaweicloud_hss_container_kubernetes_cluster_protection_enable" "test" {
  cluster_name                = "test-cluster-name"
  cluster_id                  = "test-cluster-id"
  cluster_type                = "adding"
  charging_mode               = "on_demand"
  cce_protection_type         = "cluster_level"
  enterprise_project_id       = "0"
  prefer_packet_cycle         = true
  monitor_protection_statuses = ["error_protect"]
}
`
