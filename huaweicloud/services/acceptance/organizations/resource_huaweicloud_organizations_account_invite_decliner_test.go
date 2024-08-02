package organizations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAccountInviteDecliner_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
			acceptance.TestAccPreCheckOrganizationsInvitationId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccountInviteDecliner_basic(),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccountInviteDecliner_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_organizations_account_invite_decliner" "test" {
  invitation_id = "%s"
}
`, acceptance.HW_ORGANIZATIONS_INVITATION_ID)
}
