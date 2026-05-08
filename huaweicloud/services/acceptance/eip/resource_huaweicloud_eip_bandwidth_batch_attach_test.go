package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEipBatchAttachShareBandwidth_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckEipIDS(t)
			acceptance.TestAccPreCheckVpcEipBandwidthId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccEipBatchAttachShareBandwidth_basic(),
			},
		},
	})
}

func testAccEipBatchAttachShareBandwidth_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_eip_bandwidth_batch_attach" "test" {
  dynamic "publicips" {
    for_each = split(",", "%s")
    content {
      bandwidth_id = "%s"
      publicip_id  = publicips.value
    }
  }
}
`, acceptance.HW_EIP_IDS, acceptance.HW_VPC_BANDWIDTH_ID)
}
