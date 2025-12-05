package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccElbRecycleBinLoadBalancerRecover_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckElbLoadbalancerID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccElbRecycleBinLoadBalancerRecover_basic(),
			},
		},
	})
}

func testAccElbRecycleBinLoadBalancerRecover_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_elb_recycle_bin_loadbalancer_recover" "test" {
  loadbalancer_id = "%s"
}
`, acceptance.HW_ELB_LOADBALANCER_ID)
}
