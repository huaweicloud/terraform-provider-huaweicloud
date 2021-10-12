package iam

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/policies"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getIdentityRoleResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.IAMV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("Error creating HuaweiCloud IAM client: %s", err)
	}
	return policies.Get(client, state.Primary.ID).Extract()
}

func TestAccIdentityRole_basic(t *testing.T) {
	var role policies.Role
	var roleName = acceptance.RandomAccResourceName()
	var roleNameUpdate = roleName + "update"
	resourceName := "huaweicloud_identity_role.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&role,
		getIdentityRoleResourceFunc,
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
				Config: testAccIdentityRole_basic(roleName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", roleName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform"),
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
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", roleNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform"),
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

	rc := acceptance.InitResourceCheck(
		resourceName,
		&role,
		getIdentityRoleResourceFunc,
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
				Config: testAccIdentityRole_agency(roleName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", roleName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform"),
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
`, roleName, acceptance.HW_REGION_NAME)
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
