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
		user interface{}

		resourceName = "huaweicloud_identity_user.test"
		rc           = acceptance.InitResourceCheck(resourceName, &user, getIdentityUserResourceFunc)

		rName = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "hashicorp/random",
				// The version of the random provider must be greater than 3.3.0 to support the 'numeric' parameter.
				VersionConstraint: "3.3.0",
			},
		},
		CheckDestroy: rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityUser_basic(rName),
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
				Config: testAccIdentityUser_update(rName),
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
				Config: testAccIdentityUser_no_desc(rName),
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
		user interface{}

		resourceName = "huaweicloud_identity_user.test"
		rc           = acceptance.InitResourceCheck(resourceName, &user, getIdentityUserResourceFunc)

		rName       = acceptance.RandomAccResourceName()
		initXUserID = rName + acctest.RandString(5)
		newXUserID  = rName + acctest.RandString(5)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "hashicorp/random",
				// The version of the random provider must be greater than 3.3.0 to support the 'numeric' parameter.
				VersionConstraint: "3.3.0",
			},
		},
		CheckDestroy: rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityUser_external(rName, initXUserID),
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
				Config: testAccIdentityUser_external(rName, newXUserID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "external_identity_id", newXUserID),
				),
			},
		},
	})
}

func testAccIdentityUser_base() string {
	return fmt.Sprintf(`
resource "random_string" "test" {
  length           = 10
  min_numeric      = 1
  min_special      = 1
  min_lower        = 1
  override_special = "@!"
}
`)
}

func testAccIdentityUser_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_user" "test" {
  name        = "%[2]s"
  password    = random_string.test.result
  enabled     = true
  email       = "%[2]s@abc.com"
  description = "Created by acc test"
  
  login_protect_verification_method = "email"
}
`, testAccIdentityUser_base(), name)
}

func testAccIdentityUser_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_user" "test" {
  name        = "%[2]s"
  password    = random_string.test.result
  pwd_reset   = false
  enabled     = false
  email       = "%[2]s@abcd.com"
  description = "Updated by acc test"
}
`, testAccIdentityUser_base(), name)
}

func testAccIdentityUser_no_desc(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_user" "test" {
  name      = "%[2]s"
  password  = random_string.test.result
  pwd_reset = false
  enabled   = false
  email     = "%[2]s@abcd.com"
}
`, testAccIdentityUser_base(), name)
}

func testAccIdentityUser_external(name, xUserID string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_user" "test" {
  name                 = "%[2]s"
  password             = random_string.test.result
  description          = "IAM user with external identity id"
  external_identity_id = "%[3]s"
}
`, testAccIdentityUser_base(), name, xUserID)
}
