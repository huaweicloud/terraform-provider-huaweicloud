package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBatchDetachInternetBandwidths_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGlobalEipId(t)
			acceptance.TestAccPreCheckGlobalEipId2(t)
			acceptance.TestAccPreCheckGlobalInternetBandwidthId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testBatchDetachInternetBandwidths_basic(),
			},
		},
	})
}

func testBatchDetachInternetBandwidths_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_global_batch_detach_internet_bandwidths" "test" {
  global_eips {
    global_eip_id         = "%s"
  }
  global_eips {
    global_eip_id         = "%s"
  }
}
`, acceptance.HW_GLOBAL_EIP_ID, acceptance.HW_GLOBAL_EIP_ID2)
}
