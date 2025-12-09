package vpcep

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVpcepServiceUpgrade_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testVpcepServiceUpgrade_basic(rName),
			},
		},
	})
}

func testVpcepServiceUpgrade_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpcep_service_upgrade" "test" {
  service_id = huaweicloud_vpcep_service.test.id
}
`, testAccVPCEPService_Basic(rName))
}
