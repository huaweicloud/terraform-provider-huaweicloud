package identitycenter

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

func getGroupMembershipResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getGroupMembership: query Identity Center group membership
	var (
		getGroupMembershipHttpUrl = "v1/identity-stores/{identity_store_id}/group-memberships/{membership_id}"
		getGroupMembershipProduct = "identitystore"
	)
	getGroupMembershipClient, err := cfg.NewServiceClient(getGroupMembershipProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Identity Center Client: %s", err)
	}

	getGroupMembershipPath := getGroupMembershipClient.Endpoint + getGroupMembershipHttpUrl
	getGroupMembershipPath = strings.ReplaceAll(getGroupMembershipPath, "{identity_store_id}",
		state.Primary.Attributes["identity_store_id"])
	getGroupMembershipPath = strings.ReplaceAll(getGroupMembershipPath, "{membership_id}", state.Primary.ID)

	getGroupMembershipOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getGroupMembershipResp, err := getGroupMembershipClient.Request("GET", getGroupMembershipPath,
		&getGroupMembershipOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Identity Center Group Membership: %s", err)
	}
	return utils.FlattenResponse(getGroupMembershipResp)
}

func TestAccGroupMembership_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_identitycenter_group_membership.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGroupMembershipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGroupMembership_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "identity_store_id",
						"data.huaweicloud_identitycenter_instance.test", "identity_store_id"),
					resource.TestCheckResourceAttrPair(rName, "group_id",
						"huaweicloud_identitycenter_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "member_id",
						"huaweicloud_identitycenter_user.test", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testGroupMembershipImportState(rName),
			},
		},
	})
}

func testGroupMembership_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identitycenter_group" "test" {
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
  name              = "%[2]s"
}

resource "huaweicloud_identitycenter_group_membership" "test" {
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
  group_id          = huaweicloud_identitycenter_group.test.id
  member_id         = huaweicloud_identitycenter_user.test.id
}
`, testIdentityCenterUser_basic_update(name), name)
}

// testGroupMembershipImportState use to return an id with format <identity_store_id>/<membership_id>
func testGroupMembershipImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		identityStoreId := rs.Primary.Attributes["identity_store_id"]
		if identityStoreId == "" {
			return "", fmt.Errorf("attribute (identity_store_id) of Resource (%s) not found: %s", name, rs)
		}
		return identityStoreId + "/" + rs.Primary.ID, nil
	}
}
