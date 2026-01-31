package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Please ensure that the user executing the acceptance test has 'admin' permission.
func TestAccDataV5Users_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_identityv5_users.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byGroupId   = "data.huaweicloud_identityv5_users.filter_by_group_id"
		dcByGroupId = acceptance.InitDataSourceCheck(byGroupId)

		byUserId   = "data.huaweicloud_identityv5_users.filter_by_user_id"
		dcByUserId = acceptance.InitDataSourceCheck(byUserId)

		byWithoutLoginProfile   = "data.huaweicloud_identityv5_users.filter_by_without_login_profile"
		dcByWithoutLoginProfile = acceptance.InitDataSourceCheck(byWithoutLoginProfile)
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
		Steps: []resource.TestStep{
			{
				Config: testAccDataV5Users_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "users.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'group_id' parameter.
					dcByGroupId.CheckResourceExists(),
					resource.TestCheckOutput("is_group_id_filter_useful", "true"),
					// Filter by 'user_id' parameter.
					dcByUserId.CheckResourceExists(),
					resource.TestCheckOutput("is_user_id_filter_useful", "true"),
					// Check other attributes.
					// Current user is new created, not logged in yet, so 'last_login_at' attribute is empty and not checked.
					resource.TestCheckResourceAttrPair(byUserId, "users.0.user_id", "huaweicloud_identityv5_user.test.0", "id"),
					resource.TestCheckResourceAttrPair(byUserId, "users.0.user_name", "huaweicloud_identityv5_user.test.0", "name"),
					resource.TestCheckResourceAttrPair(byUserId, "users.0.enabled", "huaweicloud_identityv5_user.test.0", "enabled"),
					resource.TestCheckResourceAttrPair(byUserId, "users.0.description", "huaweicloud_identityv5_user.test.0", "description"),
					resource.TestCheckResourceAttrPair(byUserId, "users.0.is_root_user", "huaweicloud_identityv5_user.test.0", "is_root_user"),
					resource.TestCheckResourceAttrPair(byUserId, "users.0.created_at", "huaweicloud_identityv5_user.test.0", "created_at"),
					resource.TestCheckResourceAttrPair(byUserId, "users.0.urn", "huaweicloud_identityv5_user.test.0", "urn"),
					resource.TestCheckResourceAttr(byUserId, "users.0.tags.0.tag_key", "owner"),
					resource.TestCheckResourceAttr(byUserId, "users.0.tags.0.tag_value", "terraform"),
					resource.TestCheckResourceAttrPair(byUserId, "users.0.password_reset_required",
						"huaweicloud_identityv5_login_profile.test", "password_reset_required"),
					resource.TestCheckResourceAttrPair(byUserId, "users.0.password_expires_at",
						"huaweicloud_identityv5_login_profile.test", "password_expires_at"),
					// To test the case of new users not setting login password.
					dcByWithoutLoginProfile.CheckResourceExists(),
					resource.TestCheckOutput("without_login_psw_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataV5Users_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "test" {
  count = 2

  name        = "%[1]s${count.index}"
  description = "created by terraform script"
  enabled     = true
}

resource "huaweicloud_identityv5_group" "test" {
  group_name = "%[1]s"
}

locals {
  user_id = try(huaweicloud_identityv5_user.test[0].id, null)
}

resource "huaweicloud_identityv5_group_membership" "test" {
  group_id     = huaweicloud_identityv5_group.test.id
  user_id_list = [local.user_id]
}

resource "huaweicloud_identityv5_resource_tag" "test" {
  resource_type = "user"
  resource_id   = local.user_id

  tags = {
    owner = "terraform"
  }
}

resource "random_string" "test" {
  length           = 10
  min_numeric      = 1
  min_special      = 1
  min_lower        = 1
  override_special = "@!"
}

resource "huaweicloud_identityv5_login_profile" "test" {
  user_id                 = local.user_id
  password                = random_string.test.result
  password_reset_required = true
}
`, name)
}

func testAccDataV5Users_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameters.
data "huaweicloud_identityv5_users" "test" {
  depends_on = [huaweicloud_identityv5_user.test]
}

# Filter by 'group_id' parameter.
locals {
  group_id = huaweicloud_identityv5_group.test.id
}

data "huaweicloud_identityv5_users" "filter_by_group_id" {
  group_id = local.group_id

  depends_on = [huaweicloud_identityv5_group_membership.test]
}

output "is_group_id_filter_useful" {
  value = contains(data.huaweicloud_identityv5_users.filter_by_group_id.users[*].user_id, local.user_id)
}

# Filter by 'user_id' parameter.
data "huaweicloud_identityv5_users" "filter_by_user_id" {
  user_id = local.user_id

  depends_on = [
    huaweicloud_identityv5_resource_tag.test,
    huaweicloud_identityv5_login_profile.test,
  ]
}

locals {
  user_id_filter_result = [for v in data.huaweicloud_identityv5_users.filter_by_user_id.users[*].user_id :
  v == local.user_id]
}

output "is_user_id_filter_useful" {
  value = length(local.user_id_filter_result) > 0 && alltrue(local.user_id_filter_result)
}

# To test the case of new users not setting login password.
data "huaweicloud_identityv5_users" "filter_by_without_login_profile" {
  user_id = try(huaweicloud_identityv5_user.test[1].id, null)
}

output "without_login_psw_validation_pass" {
  value = length(data.huaweicloud_identityv5_users.filter_by_without_login_profile.users) > 0
}
`, testAccDataV5Users_base(name))
}
