package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityGroup_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_identity_group.test"
		dc     = acceptance.InitDataSourceCheck(dcName)

		rName = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityGroup_by_name(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dcName, "name", rName),
					resource.TestCheckResourceAttr(dcName, "users.#", "0"),
				),
			},
		},
	})
}

func TestAccIdentityGroup_with_user(t *testing.T) {
	var (
		dcName = "data.huaweicloud_identity_group.test"
		dc     = acceptance.InitDataSourceCheck(dcName)

		rName    = acceptance.RandomAccResourceName()
		password = acceptance.RandomPassword()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityGroup_with_user(rName, password),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dcName, "name", rName),
					resource.TestCheckResourceAttr(dcName, "users.#", "2"),
					resource.TestCheckResourceAttrSet(dcName, "users.0.id"),
				),
			},
		},
	})
}

func testAccIdentityGroup_by_name(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "test" {
  name        = "%s"
  description = "Created by acc test"
}

data "huaweicloud_identity_group" "test" {
  name = huaweicloud_identity_group.test.name
  
  depends_on = [
    huaweicloud_identity_group.test
  ]
}
`, rName)
}

func testAccIdentityGroup_with_user(rName, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "test" {
  name        = "%[1]s"
  description = "Created by acc test"
}

resource "huaweicloud_identity_user" "test" {
  count    = 2
  name     = "%[1]s-${count.index}"
  password = "%[2]s"
  enabled  = true
}

resource "huaweicloud_identity_group_membership" "test" {
  group = huaweicloud_identity_group.test.id
  users = huaweicloud_identity_user.test.*.id
}

data "huaweicloud_identity_group" "test" {
  name = huaweicloud_identity_group.test.name

  depends_on = [
    huaweicloud_identity_group_membership.test
  ]
}
`, rName, password)
}
