package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityV5AccessKey_basic(t *testing.T) {
	resourceName := "data.huaweicloud_identityv5_access_key.test"
	userName := acceptance.RandomAccResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIdentityV5AccessKey_basic(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceName, "access_key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
		},
	})
}

func testAccDataSourceIdentityV5AccessKey_basic(username string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "user_1" {
  name = "%[1]s"
}

resource "huaweicloud_identityv5_access_key" "key_1" {
  user_id = huaweicloud_identityv5_user.user_1.id
  status  = "inactive"
}

data "huaweicloud_identityv5_access_key" "test" {
  user_id = huaweicloud_identityv5_access_key.key_1.user_id
}
`, username)
}
