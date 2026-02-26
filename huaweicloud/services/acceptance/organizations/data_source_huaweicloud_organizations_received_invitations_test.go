package organizations

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataReceivedInvitations_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_organizations_received_invitations.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
			acceptance.TestAccPreCheckOrganizationsInvitationId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataReceivedInvitations_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "handshakes.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "handshakes.0.id"),
					resource.TestCheckResourceAttrSet(all, "handshakes.0.urn"),
					resource.TestCheckResourceAttrSet(all, "handshakes.0.status"),
					resource.TestCheckResourceAttrSet(all, "handshakes.0.organization_id"),
					resource.TestCheckResourceAttrSet(all, "handshakes.0.management_account_id"),
					resource.TestCheckResourceAttrSet(all, "handshakes.0.management_account_name"),
					resource.TestCheckResourceAttrSet(all, "handshakes.0.target.#"),
					resource.TestCheckResourceAttrSet(all, "handshakes.0.target.0.type"),
					resource.TestCheckResourceAttrSet(all, "handshakes.0.target.0.entity"),
					resource.TestCheckResourceAttrSet(all, "handshakes.0.created_at"),
					resource.TestCheckResourceAttrSet(all, "handshakes.0.updated_at"),
				),
			},
		},
	})
}

func testAccDataReceivedInvitations_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_organizations_account_invite_decliner" "test" {
  invitation_id = "%s"
}

# Without any filter parameters.
data "huaweicloud_organizations_received_invitations" "test" {
  depends_on = [huaweicloud_organizations_account_invite_decliner.test]
}
`, acceptance.HW_ORGANIZATIONS_INVITATION_ID)
}
