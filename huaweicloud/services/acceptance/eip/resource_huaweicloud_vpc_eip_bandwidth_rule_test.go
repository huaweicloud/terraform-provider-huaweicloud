package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEipBandwidthRule_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_vpc_eip_bandwidth_rule.test"
		rName        = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckVpcEipBandwidthName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccEipBandwidthRule_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccEipBandwidthRule_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_bandwidth" "test" {
  name = "%s"
}

resource "huaweicloud_vpc_eip_bandwidth_rule" "test" {
  bandwidth_id           = data.huaweicloud_vpc_bandwidth.test.id
  name                   = "%s"
  egress_size            = 10
  egress_guarented_size  = 10
  description            = "test bandwidth rule with public IP"
  
  publicip_info {
    publicip_id = data.huaweicloud_vpc_bandwidth.test.publicips[0].id
  }
}
`, acceptance.HW_VPC_BANDWIDTH_NAME, rName)
}
