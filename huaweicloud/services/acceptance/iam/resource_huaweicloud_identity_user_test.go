package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
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
	var (
		user         users.User
		resourceName = "huaweicloud_identity_user.test"

		rName        = acceptance.RandomAccResourceName()
		initPassword = acceptance.RandomPassword()
		newPassword  = acceptance.RandomPassword()
	)

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
				Config: testAccIdentityUser_basic(rName, initPassword),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by acc test"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "pwd_reset", "true"),
					resource.TestCheckResourceAttr(resourceName, "email", rName+"@abc.com"),
					resource.TestCheckResourceAttr(resourceName, "password_strength", "Strong"),
					resource.TestCheckResourceAttr(resourceName, "login_protect_verification_method", "email"),
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
				Config: testAccIdentityUser_update(rName, newPassword),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by acc test"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "pwd_reset", "false"),
					resource.TestCheckResourceAttr(resourceName, "email", rName+"@abcd.com"),
					resource.TestCheckResourceAttr(resourceName, "login_protect_verification_method", ""),
				),
			},
			{
				Config: testAccIdentityUser_no_desc(rName, newPassword),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
		},
	})
}

func TestAccIdentityUser_external(t *testing.T) {
	var (
		user         users.User
		resourceName = "huaweicloud_identity_user.test"

		rName       = acceptance.RandomAccResourceName()
		password    = acceptance.RandomPassword()
		initXUserID = rName + acctest.RandString(5)
		newXUserID  = rName + acctest.RandString(5)
	)

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
				Config: testAccIdentityUser_external(rName, password, initXUserID),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "IAM user with external identity id"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "pwd_reset", "true"),
					resource.TestCheckResourceAttr(resourceName, "password_strength", "Strong"),
					resource.TestCheckResourceAttr(resourceName, "external_identity_id", initXUserID),
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
				Config: testAccIdentityUser_external(rName, password, newXUserID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "external_identity_id", newXUserID),
				),
			},
		},
	})
}

func testAccIdentityUser_basic(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user" "test" {
  name        = "%[1]s"
  password    = "%[2]s"
  enabled     = true
  email       = "%[1]s@abc.com"
  description = "Created by acc test"
  
  login_protect_verification_method = "email"
}
`, name, password)
}

func testAccIdentityUser_update(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user" "test" {
  name        = "%[1]s"
  password    = "%[2]s"
  pwd_reset   = false
  enabled     = false
  email       = "%[1]s@abcd.com"
  description = "Updated by acc test"
}
`, name, password)
}

func testAccIdentityUser_no_desc(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user" "test" {
  name      = "%[1]s"
  password  = "%[2]s"
  pwd_reset = false
  enabled   = false
  email     = "%[1]s@abcd.com"
}
`, name, password)
}

func testAccIdentityUser_external(name, password, xUserID string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user" "test" {
  name                 = "%[1]s"
  password             = "%[2]s"
  description          = "IAM user with external identity id"
  external_identity_id = "%[3]s"
}
`, name, password, xUserID)
}
