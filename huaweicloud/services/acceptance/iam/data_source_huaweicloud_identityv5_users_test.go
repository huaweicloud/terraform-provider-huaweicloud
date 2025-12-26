package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityV5Users_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identityv5_users.test"
	name := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIdentityV5Users_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "users.#"),
				),
			},
			{
				Config: testAccDataSourceIdentityV5Users_showUser(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "users.#"),
				),
			},
		},
	})
}

func testAccDataSourceIdentityV5Users_basic() string {
	return `
data "huaweicloud_identityv5_users" "test" {}
`
}

func testAccDataSourceIdentityV5Users_showUser(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "user_1" {
 name        = "%[1]s"
 description = "tested by terraform"
}

data "huaweicloud_identityv5_users" "test" {
	user_id = huaweicloud_identityv5_user.user_1.id
}
`, name)
}
