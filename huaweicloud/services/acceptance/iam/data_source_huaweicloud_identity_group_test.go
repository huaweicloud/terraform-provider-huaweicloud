package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityGroupDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identity_group.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityGroupDataSource_by_name(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "users.#", "0"),
				),
			},
		},
	})
}

func TestAccIdentityGroupDataSource_with_user(t *testing.T) {
	dataSourceName := "data.huaweicloud_identity_group.test"
	rName := acceptance.RandomAccResourceName()
	password := acceptance.RandomPassword()
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityGroupDataSource_with_user(rName, password),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "users.#", "2"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.id"),
				),
			},
		},
	})
}

func testAccIdentityGroupDataSource_by_name(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "test" {
  name        = "%s"
  description = "An ACC test group"
}

data "huaweicloud_identity_group" "test" {
  name = huaweicloud_identity_group.test.name
  
  depends_on = [
    huaweicloud_identity_group.test
  ]
}
`, rName)
}

func testAccIdentityGroupDataSource_with_user(rName, password string) string {
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
