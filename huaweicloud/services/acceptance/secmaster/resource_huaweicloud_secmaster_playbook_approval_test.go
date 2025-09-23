package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccPlaybookApproval_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testPlaybookApproval_basic(name),
			},
		},
	})
}

func testPlaybookApproval_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_secmaster_playbook_version_action" "submit_version" {
  workspace_id = "%[2]s"
  version_id   = huaweicloud_secmaster_playbook_version.test.id
  status       = "APPROVING"

  depends_on = [huaweicloud_secmaster_playbook_action.test]
}

resource "huaweicloud_secmaster_playbook_approval" "test" {
  workspace_id = "%[2]s"
  version_id   = huaweicloud_secmaster_playbook_version.test.id
  result       = "PASS"
  content      = "ok"

  depends_on = [huaweicloud_secmaster_playbook_version_action.submit_version]
}
`, testPlaybookVersion_basic(name), acceptance.HW_SECMASTER_WORKSPACE_ID)
}
