package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBatchDeleteAddressGroupMember_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfwAddressGroupId(t)
			acceptance.TestAccPreCheckCfwAddressGroupMemberId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testBatchDeleteAddressGroupMember_basic(),
			},
		},
	})
}

func testBatchDeleteAddressGroupMember_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cfw_batch_delete_address_group_member" "test" {
  set_id           = "%[1]s"
  address_item_ids = ["%[2]s"]
}
`, acceptance.HW_CFW_ADDRESS_GROUP_ID, acceptance.HW_CFW_ADDRESS_GROUP_MEMBER_ID)
}
