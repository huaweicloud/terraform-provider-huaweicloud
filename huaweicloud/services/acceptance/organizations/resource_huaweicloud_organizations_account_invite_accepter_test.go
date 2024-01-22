package organizations

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getAccountInviteAccepterResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	// getAccountInviteAccepter: Query Organizations account invite accepter
	var (
		region                          = acceptance.HW_REGION_NAME
		getAccountInviteAccepterHttpUrl = "v1/organizations/handshakes/{handshake_id}"
		getAccountInviteAccepterProduct = "organizations"
	)
	getAccountInviteAccepterClient, err := cfg.NewServiceClient(getAccountInviteAccepterProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Organizations Client: %s", err)
	}

	getAccountInviteAccepterPath := getAccountInviteAccepterClient.Endpoint + getAccountInviteAccepterHttpUrl
	getAccountInviteAccepterPath = strings.ReplaceAll(getAccountInviteAccepterPath, "{handshake_id}",
		state.Primary.ID)

	getAccountInviteAccepterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getAccountInviteAccepterResp, err := getAccountInviteAccepterClient.Request("GET",
		getAccountInviteAccepterPath, &getAccountInviteAccepterOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving AccountInviteAccepter: %s", err)
	}
	return utils.FlattenResponse(getAccountInviteAccepterResp)
}

func TestAccAccountInviteAccepter_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_organizations_account_invite_accepter.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAccountInviteAccepterResourceFunc,
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
				Config: testAccountInviteAccepter_basic(),
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

func testAccountInviteAccepter_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_organizations_account_invite_accepter" "test" {
  invitation_id                 = "%s"
  leave_organization_on_destroy = true
}
`, acceptance.HW_ORGANIZATIONS_INVITATION_ID)
}
