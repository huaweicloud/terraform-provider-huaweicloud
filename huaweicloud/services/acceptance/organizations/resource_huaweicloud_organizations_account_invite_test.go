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

func getAccountInviteResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	// getAccountInvite: Query Organizations account invite
	var (
		region                  = acceptance.HW_REGION_NAME
		getAccountHttpUrl       = "v1/organizations/accounts/{account_id}"
		getAccountInviteHttpUrl = "v1/organizations/handshakes/{handshake_id}"
		getAccountInviteProduct = "organizations"
	)
	getAccountInviteClient, err := cfg.NewServiceClient(getAccountInviteProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Organizations Client: %s", err)
	}

	getAccountInvitePath := getAccountInviteClient.Endpoint + getAccountInviteHttpUrl
	getAccountInvitePath = strings.ReplaceAll(getAccountInvitePath, "{handshake_id}", state.Primary.ID)

	getAccountInviteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getAccountInviteResp, err := getAccountInviteClient.Request("GET", getAccountInvitePath,
		&getAccountInviteOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving AccountInvite: %s", err)
	}

	getAccountInviteRespBody, err := utils.FlattenResponse(getAccountInviteResp)
	if err != nil {
		return nil, err
	}

	status := utils.PathSearch("handshake.status", getAccountInviteRespBody, "")
	if status == "cancelled" || status == "expired" {
		return nil, fmt.Errorf("failed to get Organizations AccountInvite")
	}

	accountID := utils.PathSearch("handshake.target.entity", getAccountInviteRespBody, "")

	// the handshake will always exist, so it is necessary to check whether the account can be obtained if the
	// status is accepted
	if status == "accepted" {
		getAccountPath := getAccountInviteClient.Endpoint + getAccountHttpUrl
		getAccountPath = strings.ReplaceAll(getAccountPath, "{account_id}", accountID.(string))

		getAccountOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		_, err = getAccountInviteClient.Request("GET", getAccountPath, &getAccountOpt)

		if err != nil {
			return nil, fmt.Errorf("failed to get Organizations AccountInvite")
		}
	}

	return getAccountInviteRespBody, nil
}

func TestAccAccountInvite_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_organizations_account_invite.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAccountInviteResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
			acceptance.TestAccPreCheckOrganizationsInviteAccountId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccountInvite_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "account_id",
						acceptance.HW_ORGANIZATIONS_INVITE_ACCOUNT_ID),
					resource.TestCheckResourceAttrSet(rName, "urn"),
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
				ImportStateVerifyIgnore: []string{"remove_account_on_destroy"},
			},
		},
	})
}

func testAccountInvite_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_organizations_account_invite" "test" {
  account_id                = "%[1]s"
  remove_account_on_destroy = true
}
`, acceptance.HW_ORGANIZATIONS_INVITE_ACCOUNT_ID)
}
