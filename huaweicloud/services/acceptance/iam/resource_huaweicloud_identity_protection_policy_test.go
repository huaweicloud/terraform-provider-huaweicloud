package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3.0/security"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getProtectionPolicyResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.IAMV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	policy, err := security.GetProtectPolicy(client, state.Primary.ID)
	if err != nil {
		return nil, err
	}

	if policy.Protection || !policy.AllowUser.ManageAccesskey || !policy.AllowUser.ManagePassword ||
		!policy.AllowUser.ManageEmail || !policy.AllowUser.ManageMobile {
		return policy, nil
	}

	return nil, golangsdk.ErrDefault404{}
}

func TestAccProtectionPolicy_basic(t *testing.T) {
	var (
		object interface{}

		resourceName = "huaweicloud_identity_protection_policy.test"
		rc           = acceptance.InitResourceCheck(resourceName, &object, getProtectionPolicyResourceFunc)
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
				Config: testAccProtectPolicy_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "self_verification", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "self_management.#"),
				),
			},
			{
				Config: testAccProtectPolicy_basic_step2(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
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
				Config: testAccProtectPolicy_basic_step3(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
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

func testAccProtectPolicy_basic_step1() string {
	return `
resource "huaweicloud_identity_protection_policy" "test" {
  protection_enabled = true
}
`
}

// verification by email
func testAccProtectPolicy_basic_step2() string {
	return `
resource "huaweicloud_identity_protection_policy" "test" {
  protection_enabled = true
  verification_email = "example@email.com"
}
`
}

// disable the operation protection and update self_management
func testAccProtectPolicy_basic_step3() string {
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
