package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccV3UserInfo_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_identity_user_info.test"

		rName = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccV3UserInfo_basic_step1(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "email"),
					resource.TestCheckResourceAttrSet(resourceName, "mobile"),
				),
			},
			{
				Config: testAccV3UserInfo_basic_step2(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "email"),
				),
			},
			{
				Config: testAccV3UserInfo_basic_step3(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "mobile"),
				),
			},
		},
	})
}

func testAccV3UserInfo_basic_step1(email string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user_info" "test" {
  email  = format("%[1]s@example.com")
  mobile = "0086-12345678"
}
`, email)
}

func testAccV3UserInfo_basic_step2(email string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user_info" "test" {
  email = format("%[1]s@example1.com")

  enable_force_new = true
}
`, email)
}

func testAccV3UserInfo_basic_step3() string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user_info" "test" {
  mobile = "0086-123456780"

  enable_force_new = true
}
`)
}
