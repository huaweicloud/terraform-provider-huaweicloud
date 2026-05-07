package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBatchAttachInternetBandwidths_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGlobalEipIds(t)
			acceptance.TestAccPreCheckGlobalInternetBandwidthId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testBatchAttachInternetBandwidthsBasic(),
			},
		},
	})
}

func testBatchAttachInternetBandwidthsBasic() string {
	return fmt.Sprintf(`
locals {
  global_eip_ids = split(",", "%[1]s")
}

resource "huaweicloud_global_batch_attach_internet_bandwidths" "test" {
  dynamic "global_eips" {
    for_each = local.global_eip_ids
    content {
      global_eip_id         = global_eips.value
      internet_bandwidth_id = "%[2]s"
    }
  }
}
`, acceptance.HW_GLOBAL_EIP_IDS, acceptance.HW_GLOBAL_INTERNET_BANDWIDTH_ID)
}
