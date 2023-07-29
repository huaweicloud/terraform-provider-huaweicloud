package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/users"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getIdentityUserResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.IAMV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}
	return users.Get(client, state.Primary.ID).Extract()
}

func TestAccIdentityUser_basic(t *testing.T) {
	var user users.User
	userName := acceptance.RandomAccResourceName()
	initPassword := acceptance.RandomPassword()
	newPassword := acceptance.RandomPassword()
	resourceName := "huaweicloud_identity_user.user_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&user,
		getIdentityUserResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityUser_basic(userName, initPassword),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", userName),
					resource.TestCheckResourceAttr(resourceName, "description", "tested by terraform"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "pwd_reset", "true"),
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
				Config: testAccIdentityUser_update(userName, newPassword),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", userName),
					resource.TestCheckResourceAttr(resourceName, "description", "updated by terraform"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "pwd_reset", "false"),
					resource.TestCheckResourceAttr(resourceName, "email", "user_1@abcd.com"),
				),
			},
		},
	})
}

func testAccIdentityUser_basic(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user" "user_1" {
  name        = "%s"
  password    = "%s"
  enabled     = true
  email       = "user_1@abc.com"
  description = "tested by terraform"
}
`, name, password)
}

func testAccIdentityUser_update(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user" "user_1" {
  name        = "%s"
  password    = "%s"
  pwd_reset   = false
  enabled     = false
  email       = "user_1@abcd.com"
  description = "updated by terraform"
}
`, name, password)
}
