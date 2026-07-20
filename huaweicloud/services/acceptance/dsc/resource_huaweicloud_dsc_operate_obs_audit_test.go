package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceOperateObsAudit_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscBucketId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOperateObsAudit_basic(),
			},
		},
	})
}

func testAccResourceOperateObsAudit_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_operate_obs_audit" "test" {
  bucket_id        = "%s"
  operation_status = true
}
`, acceptance.HW_DSC_BUCKET_ID)
}
