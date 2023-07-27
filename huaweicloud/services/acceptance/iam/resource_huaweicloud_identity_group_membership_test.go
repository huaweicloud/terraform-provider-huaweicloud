package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3/users"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityGroupMembership_basic(t *testing.T) {
	var (
		groupName    = acceptance.RandomAccResourceName()
		userName     = acceptance.RandomAccResourceName()
		userName2    = acceptance.RandomAccResourceName()
		password     = acceptance.RandomPassword()
		resourceName = "huaweicloud_identity_group_membership.membership_1"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityGroupMembershipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityGroupMembership_basic(groupName, userName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityGroupMembershipExists(resourceName, []string{userName}),
				),
			},
			{
				Config: testAccIdentityGroupMembership_update(groupName, userName, userName2, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityGroupMembershipExists(resourceName, []string{userName, userName2}),
				),
			},
			{
				Config: testAccIdentityGroupMembership_updatedown(groupName, userName2, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityGroupMembershipExists(resourceName, []string{userName2}),
				),
			},
		},
	})
}

func testAccCheckIdentityGroupMembershipDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	identityClient, err := cfg.IdentityV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating IAM client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_identity_group_membership" {
			continue
		}

		_, err := users.ListInGroup(identityClient, rs.Primary.Attributes["group"], nil).AllPages()
		if err == nil {
			return fmt.Errorf("user still exists in group")
		}
	}

	return nil
}

func testAccCheckIdentityGroupMembershipExists(n string, us []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		identityClient, err := cfg.IdentityV3Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating IAM client: %s", err)
		}
		group := rs.Primary.Attributes["group"]
		if group == "" {
			return fmt.Errorf("No group is set")
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
			return fmt.Errorf("bad group membership compare, excepted(%d), but(%d)", len(us), len(founds))
		}

		return nil
	}
}

func testAccIdentityGroupMembership_basic(groupName, userName, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "group_1" {
  name = "%s"
}

resource "huaweicloud_identity_user" "user_1" {
  name     = "%s"
  password = "%s"
  enabled  = true
}
   
resource "huaweicloud_identity_group_membership" "membership_1" {
  group = huaweicloud_identity_group.group_1.id
  users = [huaweicloud_identity_user.user_1.id]
}
`, groupName, userName, password)
}

func testAccIdentityGroupMembership_update(groupName, userName1, userName2, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "group_1" {
  name = "%s"
}

resource "huaweicloud_identity_user" "user_1" {
  name     = "%s"
  password = "%s"
  enabled  = true
}

resource "huaweicloud_identity_user" "user_2" {
  name     = "%s"
  password = "%s"
  enabled  = true
}

   
resource "huaweicloud_identity_group_membership" "membership_1" {
  group = huaweicloud_identity_group.group_1.id
  users = [huaweicloud_identity_user.user_1.id, huaweicloud_identity_user.user_2.id]
}
`, groupName, userName1, password, userName2, password)
}

func testAccIdentityGroupMembership_updatedown(groupName, userName, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "group_1" {
  name = "%s"
}

resource "huaweicloud_identity_user" "user_2" {
  name     = "%s"
  password = "%s"
  enabled  = true
}

resource "huaweicloud_identity_group_membership" "membership_1" {
  group = huaweicloud_identity_group.group_1.id
  users = [huaweicloud_identity_user.user_2.id]
}
`, groupName, userName, password)
}
