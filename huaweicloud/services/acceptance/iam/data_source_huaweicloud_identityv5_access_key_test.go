package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Please ensure that the user executing the acceptance test has 'admin' permission.
func TestAccDataV5AccessKey_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		dcName = "data.huaweicloud_identityv5_access_key.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataV5AccessKey_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(dcName, "user_id", "huaweicloud_identityv5_user.test", "id"),
					resource.TestCheckResourceAttrPair(dcName, "access_key_id", "huaweicloud_identityv5_access_key.test", "access_key_id"),
					resource.TestCheckResourceAttr(dcName, "status", "inactive"),
					resource.TestCheckResourceAttrPair(dcName, "created_at", "huaweicloud_identityv5_access_key.test", "created_at"),
				),
			},
		},
	})
}

func testAccDataV5AccessKey_basic(username string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "test" {
  name = "%[1]s"
}

resource "huaweicloud_identityv5_access_key" "test" {
  user_id = huaweicloud_identityv5_user.test.id
  status  = "inactive"
}

data "huaweicloud_identityv5_access_key" "test" {
  user_id = huaweicloud_identityv5_access_key.test.user_id
}
`, username)
}
