package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccV3UserPassword_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_identity_user_password.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckV3UserPassword(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccV3UserPassword_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testAccV3UserPassword_basic_step2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccV3UserPassword_basic_step1() string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user_password" "test" {
  original_password = "%[1]s"
  password          = "%[2]s"
}
`, acceptance.HW_IDENTITY_ORIGINAL_PASSWORD, acceptance.HW_IDENTITY_NEW_PASSWORD)
}

func testAccV3UserPassword_basic_step2() string {
	return fmt.Sprintf(`
// change the password to the original password

resource "huaweicloud_identity_user_password" "test" {
  original_password = "%[1]s"  
  password          = "%[2]s"

  enable_force_new = true
}
`, acceptance.HW_IDENTITY_NEW_PASSWORD, acceptance.HW_IDENTITY_ORIGINAL_PASSWORD)
}
