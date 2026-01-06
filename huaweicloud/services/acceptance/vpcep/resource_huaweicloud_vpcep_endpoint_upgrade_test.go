package vpcep

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVpcepEndpointUpgrade_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testVpcepEndpointUpgrade_basic(rName),
			},
		},
	})
}

func testVpcepEndpointUpgrade_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpcep_endpoint_upgrade" "test" {
  endpoint_id = huaweicloud_vpcep_endpoint.test.id
}
`, testAccVPCEndpoint_Basic(rName))
}
