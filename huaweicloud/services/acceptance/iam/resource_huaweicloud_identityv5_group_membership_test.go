package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iam"
)

func getV5GroupMembershipFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	return iam.GetV5GroupassociateUsers(client, state.Primary.ID, nil)
}

// Please ensure that the user executing the acceptance test has 'admin' permission.
func TestAccV5GroupMembership_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		obj    interface{}
		rName  = "huaweicloud_identityv5_group_membership.test"
		rName2 = "huaweicloud_identityv5_group_membership.test2"
		rc     = acceptance.InitResourceCheck(rName, &obj, getV5GroupMembershipFunc)
		rc2    = acceptance.InitResourceCheck(rName2, &obj, getV5GroupMembershipFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rc.CheckResourceDestroy(),
			rc2.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccV5GroupMembership_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "group_id", "huaweicloud_identityv5_group.test", "id"),
					resource.TestMatchResourceAttr(rName, "user_id_list.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(rName, "users.#", "2"),
					resource.TestCheckResourceAttr(rName, "users_origin.#", "2"),
					resource.TestCheckResourceAttrSet(rName, "users.0.user_id"),
					resource.TestCheckResourceAttrSet(rName, "users.0.user_name"),
					resource.TestCheckResourceAttrSet(rName, "users.0.is_root_user"),
					resource.TestCheckResourceAttrSet(rName, "users.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "users.0.urn"),
					resource.TestCheckResourceAttrSet(rName, "users.0.description"),
					resource.TestCheckResourceAttrSet(rName, "users.0.enabled"),
					rc2.CheckResourceExists(),
					// Users number is the sum of itself and the number of huaweicloud_identityv5_user.test, so it is 3.
					resource.TestCheckResourceAttr(rName2, "users.#", "3"),
					resource.TestCheckResourceAttr(rName2, "users_origin.#", "1"),
				),
			},
			{
				Config: testAccV5GroupMembership_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					// After resources refreshed, the users will be overridden as all users under the same group.
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "users.#", "3"),
					resource.TestCheckResourceAttr(rName, "users_origin.#", "2"),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "users.#", "3"),
					resource.TestCheckResourceAttr(rName2, "users_origin.#", "1"),
				),
			},
			{
				Config: testAccV5GroupMembership_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					// When multiple resources are used to manage the same group, users will store the results
					// modified by other resources, resulting in users displaying all binding results except for the
					// first change.
					resource.TestMatchResourceAttr(rName, "users.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(rName, "users_origin.#", "1"),
					rc2.CheckResourceExists(),
					resource.TestMatchResourceAttr(rName2, "users.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(rName2, "users_origin.#", "1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"users_origin", // Only different in the acceptance test.
					"users",        // The order in the script differs from the order returned by the API.
				},
			},
			// After importing the resources, users will contain all the users under the same group.
			{
				Config: testAccV5GroupMembership_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "users.#", "2"),
					resource.TestCheckResourceAttr(rName, "users_origin.#", "1"),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "users.#", "2"),
					resource.TestCheckResourceAttr(rName2, "users_origin.#", "1"),
				),
			},
		},
	})
}

func testAccV5GroupMembership_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_group" "test" {
  group_name = "%[1]s"
}

resource "huaweicloud_identityv5_user" "test" {
  count = 4

  name        = "%[1]s${count.index}"
  description = "Created %[1]s${count.index} by terraform"
  enabled     = true
}

locals {
  user_ids = huaweicloud_identityv5_user.test[*].id
}
`, name)
}

func testAccV5GroupMembership_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identityv5_group_membership" "test" {
  group_id     = huaweicloud_identityv5_group.test.id

  dynamic "users" {
    for_each = slice(local.user_ids, 0, 2)

    content {
      user_id = users.value
    }
  }
}

resource "huaweicloud_identityv5_group_membership" "test2" {
  group_id     = huaweicloud_identityv5_group.test.id

  dynamic "users" {
    for_each = slice(local.user_ids, 3, 4)

    content {
      user_id = users.value
    }
  }

  # Wait for the first resource to bind users successfully, then bind new users.
  depends_on = [huaweicloud_identityv5_group_membership.test]
}
`, testAccV5GroupMembership_base(name))
}

func testAccV5GroupMembership_basic_step2(name string) string {
	// Refresh users for all resources.
	return testAccV5GroupMembership_basic_step1(name)
}

func testAccV5GroupMembership_basic_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identityv5_group_membership" "test" {
  group_id     = huaweicloud_identityv5_group.test.id

  dynamic "users" {
    for_each = slice(local.user_ids, 0, 1)

    content {
      user_id = users.value
    }
  }
}

resource "huaweicloud_identityv5_group_membership" "test2" {
  group_id     = huaweicloud_identityv5_group.test.id

  dynamic "users" {
    for_each = slice(local.user_ids, 2, 3)

    content {
      user_id = users.value
    }
  }
}
`, testAccV5GroupMembership_base(name))
}

// Please ensure that the user executing the acceptance test has 'admin' permission.
func TestAccV5GroupMembership_deprecated(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_identityv5_group_membership.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getV5GroupMembershipFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV5GroupMembership_deprecated_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "group_id", "huaweicloud_identityv5_group.test", "id"),
					resource.TestCheckResourceAttr(rName, "user_id_list.#", "2"),
				),
			},
			{
				Config: testAccV5GroupMembership_deprecated_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "user_id_list.#", "2"),
					resource.TestCheckResourceAttr(rName, "users.#", "2"),
					resource.TestCheckResourceAttrSet(rName, "users.0.user_id"),
					resource.TestCheckResourceAttrSet(rName, "users.0.user_name"),
					resource.TestCheckResourceAttr(rName, "users.0.is_root_user", "false"),
					resource.TestCheckResourceAttrSet(rName, "users.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "users.0.urn"),
					resource.TestCheckResourceAttrSet(rName, "users.0.description"),
					resource.TestCheckResourceAttr(rName, "users.0.enabled", "true"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccV5GroupMembership_deprecated_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_group" "test" {
  group_name = "%[1]s"
}

resource "huaweicloud_identityv5_user" "test" {
  count = 3

  name        = "%[1]s${count.index}"
  description = "Created %[1]s${count.index} by terraform"
  enabled     = true
}
`, name)
}

func testAccV5GroupMembership_deprecated_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identityv5_group_membership" "test" {
  group_id     = huaweicloud_identityv5_group.test.id
  user_id_list = slice(huaweicloud_identityv5_user.test[*].id, 0, 2)
}
`, testAccV5GroupMembership_deprecated_base(name))
}

func testAccV5GroupMembership_deprecated_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identityv5_group_membership" "test" {
  group_id     = huaweicloud_identityv5_group.test.id
  user_id_list = slice(huaweicloud_identityv5_user.test[*].id, 1, 3)
}
`, testAccV5GroupMembership_deprecated_base(name))
}
