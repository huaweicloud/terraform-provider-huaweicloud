package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityUserPassword_basic(t *testing.T) {
	password := "TestPassword#$%111"
	originalPassword := "TestPassword#$%"
	resourceName := "huaweicloud_identity_user_password.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityUserPassword(password, originalPassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccIdentityUserPassword(password, originalPassword string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user_password" "test" {
  password          = "%[1]s"
  original_password = "%[2]s"
}
`, password, originalPassword)
}
