package organizations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOrganizationsReceivedInvitations_basic(t *testing.T) {
	dataSource := "data.huaweicloud_organizations_received_invitations.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

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
				Config: testDataSourceDataSourceOrganizationsReceivedInvitations_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "handshakes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "handshakes.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "handshakes.0.urn"),
					resource.TestCheckResourceAttrSet(dataSource, "handshakes.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "handshakes.0.organization_id"),
					resource.TestCheckResourceAttrSet(dataSource, "handshakes.0.management_account_id"),
					resource.TestCheckResourceAttrSet(dataSource, "handshakes.0.management_account_name"),
					resource.TestCheckResourceAttrSet(dataSource, "handshakes.0.target.#"),
					resource.TestCheckResourceAttrSet(dataSource, "handshakes.0.target.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "handshakes.0.target.0.entity"),
					resource.TestCheckResourceAttrSet(dataSource, "handshakes.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "handshakes.0.updated_at"),
				),
			},
		},
	})
}

func testDataSourceDataSourceOrganizationsReceivedInvitations_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_organizations_account_invite_accepter" "test" {
  invitation_id = "%s"
}

data "huaweicloud_organizations_received_invitations" "test" {
  depends_on = [huaweicloud_organizations_account_invite_accepter.test]
}
`, acceptance.HW_ORGANIZATIONS_INVITATION_ID)
}
