package iam

import (
	"fmt"
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

	return iam.GetV5GroupassociateUsers(client, state.Primary.ID)
}

// Please ensure that the user executing the acceptance test has 'admin' permission.
func TestAccV5GroupMembership_basic(t *testing.T) {
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
				Config: testAccV5GroupMembership_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "group_id", "huaweicloud_identityv5_group.test", "id"),
					resource.TestCheckResourceAttr(rName, "user_id_list.#", "2"),
				),
			},
			{
				Config: testAccV5GroupMembership_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "user_id_list.#", "2"),
					resource.TestCheckResourceAttr(rName, "users.#", "2"),
				),
			},
			{
				Config: testAccV5GroupMembership_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "user_id_list.#", "1"),
					resource.TestCheckResourceAttrPair(rName, "user_id_list.0", "huaweicloud_identityv5_user.test.1", "id"),
					resource.TestCheckResourceAttr(rName, "users.#", "1"),
					resource.TestCheckResourceAttrPair(rName, "users.0.user_id", "huaweicloud_identityv5_user.test.1", "id"),
					resource.TestCheckResourceAttrPair(rName, "users.0.user_name", "huaweicloud_identityv5_user.test.1", "name"),
					resource.TestCheckResourceAttrPair(rName, "users.0.is_root_user", "huaweicloud_identityv5_user.test.1", "is_root_user"),
					resource.TestCheckResourceAttrPair(rName, "users.0.created_at", "huaweicloud_identityv5_user.test.1", "created_at"),
					resource.TestCheckResourceAttrPair(rName, "users.0.urn", "huaweicloud_identityv5_user.test.1", "urn"),
					resource.TestCheckResourceAttrPair(rName, "users.0.description", "huaweicloud_identityv5_user.test.1", "description"),
					resource.TestCheckResourceAttrPair(rName, "users.0.enabled", "huaweicloud_identityv5_user.test.1", "enabled"),
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

func testAccV5GroupMembership_base(name string) string {
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

func testAccV5GroupMembership_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identityv5_group_membership" "test" {
  group_id     = huaweicloud_identityv5_group.test.id
  user_id_list = slice(huaweicloud_identityv5_user.test[*].id, 0, 2)
}
`, testAccV5GroupMembership_base(name))
}

func testAccV5GroupMembership_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identityv5_group_membership" "test" {
  group_id     = huaweicloud_identityv5_group.test.id
  user_id_list = slice(huaweicloud_identityv5_user.test[*].id, 1, 3)
}
`, testAccV5GroupMembership_base(name))
}

func testAccV5GroupMembership_basic_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identityv5_group_membership" "test" {
  group_id     = huaweicloud_identityv5_group.test.id
  user_id_list = slice(huaweicloud_identityv5_user.test[*].id, 1, 2)
}
`, testAccV5GroupMembership_base(name))
}
