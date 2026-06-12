package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDeleteNodes_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterNodeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDeleteNodes_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("is_delete_result_not_empty", "true"),
				),
			},
		},
	})
}

func testDeleteNodes_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_delete_nodes" "test" {
  workspace_id = "%[1]s"
  delete_ids   = ["%[2]s"]
}

output "is_delete_result_not_empty" {
  value = (length(huaweicloud_secmaster_delete_nodes.test.delete_fail_list) > 0 ||
    length(huaweicloud_secmaster_delete_nodes.test.delete_success_list) > 0)
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_NODE_ID)
}
