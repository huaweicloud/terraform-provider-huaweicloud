package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccComputeInstanceRedeploy_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckECSID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceRedeploy_basic(),
			},
		},
	})
}

func testAccComputeInstanceRedeploy_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_compute_instance_redeploy" "test" {
  server_id = "%s"
}
`, acceptance.HW_ECS_ID)
}
