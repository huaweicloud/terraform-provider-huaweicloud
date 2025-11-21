package iam

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3/users"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityV5GroupMembership_basic(t *testing.T) {
	var (
		groupName    = acceptance.RandomAccResourceName()
		userName1    = acceptance.RandomAccResourceName()
		userName2    = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_identityv5_group_membership.membership_1"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV5GroupMembershipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityGroupV5Membership_basic(groupName, userName1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "users.#"),
				),
			},
			{
				Config: testAccIdentityGroupV5Membership_update(groupName, userName1, userName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "users.#"),
				),
			},
			{
				Config: testAccIdentityGroupV5Membership_updatedown(groupName, userName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "users.#"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"group_id", "user_id_list"},
			},
		},
	})
}

func testAccCheckIdentityV5GroupMembershipDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := cfg.IAMNoVersionClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating IAM client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_identityv5_group_membership" {
			continue
		}

		_, err := users.ListInGroup(client, rs.Primary.Attributes["group_id"], nil).AllPages()
		if err == nil {
			return errors.New("user still exists in group")
		}
	}

	return nil
}

func testAccIdentityGroupV5Membership_basic(groupName, userName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_group" "group_1" {
  group_name = "%s"
}

resource "huaweicloud_identityv5_user" "user_1" {
  name        = "%s"
  description = "tested by terraform"
  enabled     = true
}
   
resource "huaweicloud_identityv5_group_membership" "membership_1" {
  group_id     = huaweicloud_identityv5_group.group_1.id
  user_id_list = [huaweicloud_identityv5_user.user_1.id]
}
`, groupName, userName)
}

func testAccIdentityGroupV5Membership_update(groupName, userName1, userName2 string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_group" "group_1" {
  group_name = "%s"
}

resource "huaweicloud_identityv5_user" "user_1" {
  name        = "%s"
  description = "tested by terraform"
  enabled     = true
}

resource "huaweicloud_identityv5_user" "user_2" {
  name        = "%s"
  description = "tested by terraform"
  enabled     = true
}

resource "huaweicloud_identityv5_group_membership" "membership_1" {
  group_id     = huaweicloud_identityv5_group.group_1.id
  user_id_list = [huaweicloud_identityv5_user.user_1.id, huaweicloud_identityv5_user.user_2.id]
}
`, groupName, userName1, userName2)
}

func testAccIdentityGroupV5Membership_updatedown(groupName, userName2 string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_group" "group_1" {
  group_name = "%s"
}

resource "huaweicloud_identityv5_user" "user_2" {
  name        = "%s"
  description = "tested by terraform"
  enabled     = true
}

resource "huaweicloud_identityv5_group_membership" "membership_1" {
  group_id        = huaweicloud_identityv5_group.group_1.id
  user_id_list = [huaweicloud_identityv5_user.user_2.id]
}
`, groupName, userName2)
}
