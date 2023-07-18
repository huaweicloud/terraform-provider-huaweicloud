package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/security"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccPasswordPolicy_basic(t *testing.T) {
	resourceName := "huaweicloud_identity_password_policy.enhanced"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckPasswordPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPasswordPolicy_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPasswordPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "password_char_combination", "4"),
					resource.TestCheckResourceAttr(resourceName, "minimum_password_length", "12"),
					resource.TestCheckResourceAttr(resourceName, "number_of_recent_passwords_disallowed", "2"),
					resource.TestCheckResourceAttr(resourceName, "password_validity_period", "180"),
					resource.TestCheckResourceAttr(resourceName, "minimum_password_age", "0"),
					resource.TestCheckResourceAttr(resourceName, "maximum_consecutive_identical_chars", "0"),
					resource.TestCheckResourceAttr(resourceName, "password_not_username_or_invert", "true"),
				),
			},
			{
				Config: testAccPasswordPolicy_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "password_char_combination", "2"),
					resource.TestCheckResourceAttr(resourceName, "minimum_password_length", "12"),
					resource.TestCheckResourceAttr(resourceName, "number_of_recent_passwords_disallowed", "1"),
					resource.TestCheckResourceAttr(resourceName, "password_validity_period", "90"),
					resource.TestCheckResourceAttr(resourceName, "minimum_password_age", "60"),
					resource.TestCheckResourceAttr(resourceName, "maximum_consecutive_identical_chars", "4"),
					resource.TestCheckResourceAttr(resourceName, "password_not_username_or_invert", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPasswordPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID is not set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := cfg.IAMV3Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating IAM client: %s", err)
		}

		_, err = security.GetPasswordPolicy(client, rs.Primary.ID)
		return err
	}
}

func testAccCheckPasswordPolicyDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := cfg.IAMV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating IAM client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_identity_password_policy" {
			continue
		}

		policy, err := security.GetPasswordPolicy(client, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching the IAM account password policy")
		}

		if policy.MinCharCombination != 2 || policy.MinPasswordLength != 8 || policy.RecentPasswordsDisallowedCount != 1 ||
			policy.PasswordValidityPeriod != 0 || policy.MinPasswordAge != 0 {
			return fmt.Errorf("the password policy failed to reset to defaults")
		}
	}

	return nil
}

func testAccPasswordPolicy_basic() string {
	return `
resource "huaweicloud_identity_password_policy" "enhanced" {
  password_char_combination             = 4
  minimum_password_length               = 12
  number_of_recent_passwords_disallowed = 2
  password_validity_period              = 180 
}
`
}

func testAccPasswordPolicy_update() string {
	return `
resource "huaweicloud_identity_password_policy" "enhanced" {
  password_char_combination             = 2
  minimum_password_length               = 12
  number_of_recent_passwords_disallowed = 1
  maximum_consecutive_identical_chars   = 4
  minimum_password_age                  = 60
  password_validity_period              = 90  
}
`
}
