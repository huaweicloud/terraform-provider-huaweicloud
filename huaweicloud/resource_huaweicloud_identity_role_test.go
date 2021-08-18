package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/golangsdk/openstack/identity/v3.0/policies"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccIdentityRole_basic(t *testing.T) {
	var role policies.Role
	var roleName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	var roleNameUpdate = roleName + "update"
	resourceName := "huaweicloud_identity_role.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityRole_basic(roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityRoleExists(resourceName, &role),
					resource.TestCheckResourceAttrPtr(resourceName, "name", &role.Name),
					resource.TestCheckResourceAttrPtr(resourceName, "description", &role.Description),
					resource.TestCheckResourceAttr(resourceName, "type", "AX"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIdentityRole_update(roleNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityRoleExists(resourceName, &role),
					resource.TestCheckResourceAttrPtr(resourceName, "name", &role.Name),
					resource.TestCheckResourceAttrPtr(resourceName, "description", &role.Description),
					resource.TestCheckResourceAttr(resourceName, "type", "AX"),
				),
			},
		},
	})
}

func TestAccIdentityRole_agency(t *testing.T) {
	var role policies.Role
	var roleName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_identity_role.agency"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityRole_agency(roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityRoleExists(resourceName, &role),
					resource.TestCheckResourceAttrPtr(resourceName, "name", &role.Name),
					resource.TestCheckResourceAttrPtr(resourceName, "description", &role.Description),
					resource.TestCheckResourceAttr(resourceName, "type", "AX"),
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

func testAccCheckIdentityRoleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	identityClient, err := config.IAMV3Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_identity_role" {
			continue
		}

		_, err := policies.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Role still exists")
		}
	}

	return nil
}

func testAccCheckIdentityRoleExists(n string, role *policies.Role) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		identityClient, err := config.IAMV3Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud identity client: %s", err)
		}

		found, err := policies.Get(identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Role not found")
		}

		*role = *found

		return nil
	}
}

func testAccIdentityRole_basic(roleName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_role" test {
  name        = "%s"
  description = "created by terraform"
  type        = "AX"
  policy      = <<EOF
{
  "Version": "1.1",
  "Statement": [
    {
      "Action": [
        "obs:bucket:GetBucketAcl"
      ],
      "Effect": "Allow",
      "Resource": [
        "obs:*:*:bucket:*"
      ]
    }
  ]
}
EOF
}
`, roleName)
}

func testAccIdentityRole_update(roleName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_role" test {
  name        = "%s"
  description = "created by terraform"
  type        = "AX"
  policy      = <<EOF
{
  "Version": "1.1",
  "Statement": [
    {
      "Action": [
        "obs:bucket:GetBucketAcl"
      ],
      "Effect": "Allow",
      "Resource": [
        "obs:*:*:bucket:*"
      ],
      "Condition": {
        "StringStartWith": {
          "g:ProjectName": [
            "%s"
          ]
        }
      }
    }
  ]
}
EOF
}
`, roleName, HW_REGION_NAME)
}

func testAccIdentityRole_agency(roleName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_role" agency {
  name        = "%s"
  description = "created by terraform"
  type        = "AX"
  policy      = <<EOF
{
  "Version": "1.1",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "iam:agencies:assume"
      ],
      "Resource": {
        "uri": [
          "/iam/agencies/07805aca-ba80-0fdd-4fbd-c00b8f888c7c",
          "/iam/agencies/16d4d672-8665-496e-a0b5-71a8ad7f2fe8"
        ]
      }
    }
  ]
}
EOF
}
`, roleName)
}
