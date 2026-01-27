package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataGroup_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_identity_group.all"
		dc  = acceptance.InitDataSourceCheck(all)
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
				Config: testAccDataGroup_basic(name),
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

func testAccDataGroup_base(name string) string {
	return fmt.Sprintf(`
resource "random_string" "test" {
  length           = 10
  min_numeric      = 1
  min_special      = 1
  min_lower        = 1
  override_special = "@!"
}

resource "huaweicloud_identity_group" "test" {
  name        = "%[1]s"
}

resource "huaweicloud_identity_user" "test" {
  count    = 2
  name     = "%[1]s-${count.index}"
  password = random_string.test.result
  enabled  = true
}

resource "huaweicloud_identity_group_membership" "test" {
  group = huaweicloud_identity_group.test.id
  users = huaweicloud_identity_user.test[*].id
}
`, name)
}

func testAccDataGroup_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_group" "all" {
  name = huaweicloud_identity_group.test.name

  # Waiting for the group membership to be created
  depends_on = [
    huaweicloud_identity_group_membership.test
  ]
}
`, testAccDataGroup_base(name))
}
