package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3/users"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getGroupMembershipResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.IdentityV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	groupId := state.Primary.ID

	allPages, err := users.ListInGroup(client, groupId, nil).AllPages()
	if err != nil {
		return nil, err
	}

	allUsers, err := users.ExtractUsers(allPages)
	if err != nil {
		return nil, err
	} else if len(allUsers) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}

	return allUsers, nil
}

func TestAccGroupMembership_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_identity_group_membership.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getGroupMembershipResourceFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "hashicorp/random",
				// The version of the random provider must be greater than 3.3.0 to support the 'numeric' parameter.
				VersionConstraint: "3.3.0",
			},
		},
		CheckDestroy: rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGroupMembership_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "group", "huaweicloud_identity_group.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "2"),
				),
			},
			{
				Config: testAccGroupMembership_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "group", "huaweicloud_identity_group.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "2"),
				),
			},
			{
				Config: testAccGroupMembership_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "group", "huaweicloud_identity_group.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "1"),
				),
			},
		},
	})
}

func testAccGroupMembership_basic_base(name string) string {
	return fmt.Sprintf(`
resource "random_string" "test" {
  length           = 10
  min_numeric      = 1
  min_special      = 1
  min_lower        = 1
  override_special = "@!"
}

resource "huaweicloud_identity_group" "test" {
  name = "%s"
}

resource "huaweicloud_identity_user" "test" {
  count = 3

  name     = format("%[1]s_%%d", count.index)
  password = random_string.test.result
  enabled  = true
}
`, name)
}

func testAccGroupMembership_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s
   
resource "huaweicloud_identity_group_membership" "test" {
  group = huaweicloud_identity_group.test.id
  users = slice(huaweicloud_identity_user.test[*].id, 0, 2)
}
`, testAccGroupMembership_basic_base(name))
}

func testAccGroupMembership_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s
   
resource "huaweicloud_identity_group_membership" "test" {
  group = huaweicloud_identity_group.test.id
  users = slice(huaweicloud_identity_user.test[*].id, 1, 3)
}
`, testAccGroupMembership_basic_base(name))
}

func testAccGroupMembership_basic_step3(name string) string {
	return fmt.Sprintf(`
%[1]s
   
resource "huaweicloud_identity_group_membership" "test" {
  group = huaweicloud_identity_group.test.id
  users = slice(huaweicloud_identity_user.test[*].id, 1, 2)
}
`, testAccGroupMembership_basic_base(name))
}
