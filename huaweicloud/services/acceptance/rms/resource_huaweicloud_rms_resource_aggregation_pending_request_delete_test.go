package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceAggregationPendingRequestDelete_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRMSRequesterAccountId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAggregationPendingRequestDelete_basic(),
			},
		},
	})
}

func testAccResourceAggregationPendingRequestDelete_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rms_resource_aggregation_pending_request_delete" "test" {
  requester_account_id = "%[1]s"
}
`, acceptance.HW_RMS_REQUESTER_ACCOUNT_ID)
}
