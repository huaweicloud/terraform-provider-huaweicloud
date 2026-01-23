package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Please ensure that the user executing the acceptance test has 'admin' permission.
func TestAccDataV5Groups_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_identityv5_groups.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byUserId   = "data.huaweicloud_identityv5_groups.filter_by_user_id"
		dcByUserId = acceptance.InitDataSourceCheck(byUserId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataV5Groups_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "groups.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'user_id' parameter.
					dcByUserId.CheckResourceExists(),
					resource.TestCheckOutput("is_user_id_filter_useful", "true"),
					resource.TestCheckResourceAttrPair(byUserId, "groups.0.group_id", "huaweicloud_identityv5_group.test", "id"),
					resource.TestCheckResourceAttrPair(byUserId, "groups.0.group_name", "huaweicloud_identityv5_group.test", "group_name"),
					resource.TestCheckResourceAttrPair(byUserId, "groups.0.urn", "huaweicloud_identityv5_group.test", "urn"),
					resource.TestCheckResourceAttrPair(byUserId, "groups.0.description", "huaweicloud_identityv5_group.test", "description"),
					resource.TestCheckResourceAttrSet(byUserId, "groups.0.created_at"),
				),
			},
		},
	})
}

func testAccDataV5Groups_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "test" {
  name = "%[1]s"
}

resource "huaweicloud_identityv5_group" "test" {
  group_name  = "%[1]s"
  description = "created by terraform script"
}

resource "huaweicloud_identityv5_group_membership" "test" {
  group_id     = huaweicloud_identityv5_group.test.id
  user_id_list = [huaweicloud_identityv5_user.test.id]
}
`, name)
}

func testAccDataV5Groups_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameters.
data "huaweicloud_identityv5_groups" "test" {
  depends_on = [huaweicloud_identityv5_group.test]
}

# Filter by 'user_id' parameter.
locals {
  user_id = huaweicloud_identityv5_user.test.id
}

data "huaweicloud_identityv5_groups" "filter_by_user_id" {
  user_id = local.user_id

  depends_on = [huaweicloud_identityv5_group_membership.test]
}

# Created a new user and only associated with one group, so here only need to check if the first group ID is equal to the created group ID.
output "is_user_id_filter_useful" {
  value = data.huaweicloud_identityv5_groups.filter_by_user_id.groups[0].group_id == huaweicloud_identityv5_group.test.id
}
`, testAccDataV5Groups_base(name))
}
