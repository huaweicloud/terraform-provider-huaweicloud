package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceOperateInstanceMetadata_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDscInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOperateInstanceMetadata_basic(),
			},
		},
	})
}

func testAccResourceOperateInstanceMetadata_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_operate_instance_metadata" "test" {
  instance_id = "%s"
  action      = "REFRESH"
}
`, acceptance.HW_DSC_INSTANCE_ID)
}
