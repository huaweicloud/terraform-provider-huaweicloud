package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccClonePlaybookAndVersion_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterVersionId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccClonePlaybookAndVersion_basic(name),
			},
		},
	})
}

func testAccClonePlaybookAndVersion_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_clone_playbook_version" "test" {
  workspace_id = "%[1]s"
  version_id   = "%[2]s"
  name         = "%[3]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_VERSION_ID, name)
}
