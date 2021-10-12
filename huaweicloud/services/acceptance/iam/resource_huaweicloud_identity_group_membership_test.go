package iam

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3/users"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccIdentityV3GroupMembership_basic(t *testing.T) {
	var groupName = acceptance.RandomAccResourceName()
	var userName = acceptance.RandomAccResourceName()
	var userName2 = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_identity_group_membership.membership_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3GroupMembershipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3GroupMembership_basic(groupName, userName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3GroupMembershipExists(resourceName, []string{userName}),
				),
			},
			{
				Config: testAccIdentityV3GroupMembership_update(groupName, userName, userName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3GroupMembershipExists(resourceName, []string{userName, userName2}),
				),
			},
			{
				Config: testAccIdentityV3GroupMembership_updatedown(groupName, userName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3GroupMembershipExists(resourceName, []string{userName2}),
				),
			},
		},
	})
}

func testAccCheckIdentityV3GroupMembershipDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	identityClient, err := config.IdentityV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_identity_group_membership" {
			continue
		}

		_, err := users.ListInGroup(identityClient, rs.Primary.Attributes["group"], nil).AllPages()

		if err == nil {
			return fmtp.Errorf("User still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3GroupMembershipExists(n string, us []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		identityClient, err := config.IdentityV3Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud identity client: %s", err)
		}
		group := rs.Primary.Attributes["group"]
		if group == "" {
			return fmtp.Errorf("No group is set")
		}

		pages, err := users.ListInGroup(identityClient, group, nil).AllPages()
		if err != nil {
			return err
		}

		founds, err := users.ExtractUsers(pages)
		if err != nil {
			return err
		}

		uc := len(us)
		for _, u := range us {
			for _, f := range founds {
				if f.Name == u {
					uc--
				}
			}
		}

		if uc > 0 {
			return fmtp.Errorf("Bad group membership compare, excepted(%d), but(%d)", len(us), len(founds))
		}

		return nil
	}
}

func testAccIdentityV3GroupMembership_basic(groupName, userName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "group_1" {
  name = "%s"
}

resource "huaweicloud_identity_user" "user_1" {
  name     = "%s"
  password = "password123@#"
  enabled  = true
}
   
resource "huaweicloud_identity_group_membership" "membership_1" {
  group = huaweicloud_identity_group.group_1.id
  users = [huaweicloud_identity_user.user_1.id]
}
`, groupName, userName)
}

func testAccIdentityV3GroupMembership_update(groupName, userName string, userName2 string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "group_1" {
  name = "%s"
}

resource "huaweicloud_identity_user" "user_1" {
  name     = "%s"
  password = "password123@#"
  enabled  = true
}

resource "huaweicloud_identity_user" "user_2" {
  name     = "%s"
  password = "password123@#"
  enabled  = true
}

   
resource "huaweicloud_identity_group_membership" "membership_1" {
  group = huaweicloud_identity_group.group_1.id
  users = [huaweicloud_identity_user.user_1.id, huaweicloud_identity_user.user_2.id]
}
`, groupName, userName, userName2)
}

func testAccIdentityV3GroupMembership_updatedown(groupName, userName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "group_1" {
  name = "%s"
}

resource "huaweicloud_identity_user" "user_2" {
  name     = "%s"
  password = "password123@#"
  enabled  = true
}

resource "huaweicloud_identity_group_membership" "membership_1" {
  group = huaweicloud_identity_group.group_1.id
  users = [huaweicloud_identity_user.user_2.id]
}
`, groupName, userName)
}
