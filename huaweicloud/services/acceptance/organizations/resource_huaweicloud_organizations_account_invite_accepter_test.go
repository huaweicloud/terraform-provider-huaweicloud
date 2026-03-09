package organizations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/organizations"
)

func getAccountInviteAccepterResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	getAccountInviteAccepterClient, err := cfg.NewServiceClient("organizations", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Organizations client: %s", err)
	}

	return organizations.GetAccountInviteById(getAccountInviteAccepterClient, state.Primary.ID)
}

func TestAccAccountInviteAccepter_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_organizations_account_invite_accepter.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getAccountInviteAccepterResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
			acceptance.TestAccPreCheckOrganizationsInvitationId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAccountInviteAccepter_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "invitation_id", acceptance.HW_ORGANIZATIONS_INVITATION_ID),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttrSet(rName, "account_id"),
					resource.TestCheckResourceAttrSet(rName, "master_account_id"),
					resource.TestCheckResourceAttrSet(rName, "master_account_name"),
					resource.TestCheckResourceAttrSet(rName, "organization_id"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"leave_organization_on_destroy"},
			},
		},
	})
}

func testAccAccountInviteAccepter_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_organizations_account_invite_accepter" "test" {
  invitation_id                 = "%s"
  leave_organization_on_destroy = true
}
`, acceptance.HW_ORGANIZATIONS_INVITATION_ID)
}
