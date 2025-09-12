package vpn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccP2CGatewayJobDelete_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckVPNP2cGatewayJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testP2CGatewayJobDelete_basic(),
			},
		},
	})
}

func testP2CGatewayJobDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_vpn_p2c_gateway_job_delete" "test" {
  job_id = "%s"
}
`, acceptance.HW_VPN_GATEWAY_JOB_ID)
}
