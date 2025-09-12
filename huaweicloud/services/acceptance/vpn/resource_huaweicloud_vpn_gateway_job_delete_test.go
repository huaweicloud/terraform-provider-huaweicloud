package vpn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGatewayJobDelete_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckVPNGatewayJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testGatewayJobDelete_basic(),
			},
		},
	})
}

func testGatewayJobDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_vpn_gateway_job_delete" "test" {
  job_id = "%s"
}
`, acceptance.HW_VPN_GATEWAY_JOB_ID)
}
