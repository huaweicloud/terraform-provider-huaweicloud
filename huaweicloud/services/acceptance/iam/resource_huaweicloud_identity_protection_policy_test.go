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

func TestAccProtectionPolicy_basic(t *testing.T) {
	resourceName := "huaweicloud_identity_protection_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckProtectionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProtectPolicy_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProtectionPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "self_verification", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "self_management.#"),
				),
			},
			{
				Config: testAccProtectPolicy_email(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "self_verification", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "verification_email"),
					resource.TestCheckResourceAttrSet(resourceName, "self_management.#"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccProtectPolicy_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "protection_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "self_management.0.access_key", "true"),
					resource.TestCheckResourceAttr(resourceName, "self_management.0.password", "true"),
					resource.TestCheckResourceAttr(resourceName, "self_management.0.email", "false"),
					resource.TestCheckResourceAttr(resourceName, "self_management.0.mobile", "false"),
				),
			},
		},
	})
}

func testAccCheckProtectionPolicyExists(n string) resource.TestCheckFunc {
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

		_, err = security.GetProtectPolicy(client, rs.Primary.ID)
		return err
	}
}

func testAccCheckProtectionPolicyDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := cfg.IAMV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating IAM client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_identity_protection_policy" {
			continue
		}

		policy, err := security.GetProtectPolicy(client, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching the IAM account password policy")
		}

		if policy.Protection {
			return fmt.Errorf("the operation protection failed to reset to defaults")
		}

		if !policy.AllowUser.ManageAccesskey || !policy.AllowUser.ManagePassword ||
			!policy.AllowUser.ManageEmail || !policy.AllowUser.ManageMobile {
			return fmt.Errorf("the self-management failed to reset to defaults")
		}
	}

	return nil
}

func testAccProtectPolicy_basic() string {
	return `
resource "huaweicloud_identity_protection_policy" "test" {
  protection_enabled = true
}
`
}

// verification by email
func testAccProtectPolicy_email() string {
	return `
resource "huaweicloud_identity_protection_policy" "test" {
  protection_enabled = true
  verification_email = "example@email.com"
}
`
}

// disable the operation protection and update self_management
func testAccProtectPolicy_update() string {
	return `
resource "huaweicloud_identity_protection_policy" "test" {
  protection_enabled = false
  self_management {
    access_key = true
    password   = true
    email      = false
    mobile     = false
  }
}
`
}
