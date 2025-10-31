package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityUserInfo_basic(t *testing.T) {
	email := "IAMEmail11@huawei.com"
	mobile := "0086-123456789"
	resourceName := "huaweicloud_identity_user_info.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityUserInfo(email, mobile),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testAccIdentityUserInfoEmail(email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testAccIdentityUserInfoPhone(mobile),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccIdentityUserInfo(email, mobile string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user_info" "test" {
  email  = "%[1]s"
  mobile = "%[2]s"
}
`, email, mobile)
}

func testAccIdentityUserInfoEmail(email string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user_info" "test" {
  email  = "%[1]s"
}
`, "a"+email)
}

func testAccIdentityUserInfoPhone(mobile string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user_info" "test" {
  mobile = "%[1]s"
}
`, mobile)
}
