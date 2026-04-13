package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBatchDeleteServiceGroupMembers_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfwServiceGroupId(t)
			acceptance.TestAccPreCheckCfwServiceGroupMemberId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testBatchDeleteServiceGroupMembers_basic(),
			},
		},
	})
}

func testBatchDeleteServiceGroupMembers_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cfw_batch_delete_service_group_members" "test" {
  set_id           = "%[1]s"
  service_item_ids = ["%[2]s"]
}
`, acceptance.HW_CFW_SERVICE_GROUP_ID, acceptance.HW_CFW_SERVICE_GROUP_MEMBER_ID)
}
