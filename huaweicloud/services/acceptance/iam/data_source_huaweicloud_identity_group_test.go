package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataGroup_basic(t *testing.T) {
	var (
		name     = acceptance.RandomAccResourceName()
		password = acceptance.RandomPassword()

		all = "data.huaweicloud_identity_group.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataGroup_basic(name, password),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "name", name),
					resource.TestCheckResourceAttr(all, "users.#", "2"),
					resource.TestCheckResourceAttrSet(all, "users.0.id"),
					resource.TestCheckResourceAttrSet(all, "users.0.name"),
					resource.TestCheckResourceAttrSet(all, "users.0.enabled"),
				),
			},
		},
	})
}

func testAccDataGroup_base(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "test" {
  name        = "%[1]s"
  description = "An ACC test group"
}

resource "huaweicloud_identity_user" "test" {
  count    = 2
  name     = "%[1]s-${count.index}"
  password = "%[2]s"
  enabled  = true
}

resource "huaweicloud_identity_group_membership" "test" {
  group = huaweicloud_identity_group.test.id
  
  users = [
    for user in huaweicloud_identity_user.test : user.id
  ]
}
`, name, password)
}

func testAccDataGroup_basic(name, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_group" "all" {
  name = huaweicloud_identity_group.test.name

  # Waiting for the group membership to be created
  depends_on = [
    huaweicloud_identity_group_membership.test
  ]
}
`, testAccDataGroup_base(name, password))
}
