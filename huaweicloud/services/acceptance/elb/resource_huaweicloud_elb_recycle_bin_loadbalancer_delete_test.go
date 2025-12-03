package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccElbRecycleBinLoadBalancerDelete_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckElbLoadbalancerID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccElbRecycleBinLoadBalancerDelete_basic(),
			},
		},
	})
}

func testAccElbRecycleBinLoadBalancerDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_elb_recycle_bin_loadbalancer_delete" "test" {
  loadbalancer_id = "%s"
}
`, acceptance.HW_ELB_LOADBALANCER_ID)
}
