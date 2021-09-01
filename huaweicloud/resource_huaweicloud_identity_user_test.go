package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/users"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccIdentityV3User_basic(t *testing.T) {
	var user users.User
	var userName = fmt.Sprintf("acc-user-%s", acctest.RandString(5))
	resourceName := "huaweicloud_identity_user.user_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityV3UserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3User_basic(userName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3UserExists(resourceName, &user),
					resource.TestCheckResourceAttr(resourceName, "name", userName),
					resource.TestCheckResourceAttr(resourceName, "description", "tested by terraform"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "email", "user_1@abc.com"),
					resource.TestCheckResourceAttr(resourceName, "password_strength", "Strong"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
				},
			},
			{
				Config: testAccIdentityV3User_update(userName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3UserExists(resourceName, &user),
					resource.TestCheckResourceAttr(resourceName, "name", userName),
					resource.TestCheckResourceAttr(resourceName, "description", "updated by terraform"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "email", "user_1@abcd.com"),
				),
			},
		},
	})
}

func testAccCheckIdentityV3UserDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	iamClient, err := config.IAMV3Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_identity_user" {
			continue
		}

		_, err := users.Get(iamClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("User still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3UserExists(n string, user *users.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		iamClient, err := config.IAMV3Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud identity client: %s", err)
		}

		found, err := users.Get(iamClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("User not found")
		}

		*user = *found

		return nil
	}
}

func testAccIdentityV3User_basic(userName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user" "user_1" {
  name        = "%s"
  password    = "password123@!"
  enabled     = true
  email       = "user_1@abc.com"
  description = "tested by terraform"
}
`, userName)
}

func testAccIdentityV3User_update(userName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user" "user_1" {
  name        = "%s"
  password    = "password123@!"
  enabled     = false
  email       = "user_1@abcd.com"
  description = "updated by terraform"
}
`, userName)
}
