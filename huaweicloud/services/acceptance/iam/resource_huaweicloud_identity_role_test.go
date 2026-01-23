package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/policies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getRoleResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.IAMV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}
	return policies.Get(client, state.Primary.ID).Extract()
}

func TestAccRole_basic(t *testing.T) {
	var (
		obj interface{}

		normalAssign       = "huaweicloud_identity_role.normal"
		rcNormalAssign     = acceptance.InitResourceCheck(normalAssign, &obj, getRoleResourceFunc)
		assumeAgencyAssign = "huaweicloud_identity_role.assume_agency"
		rcAssumeRoleAssign = acceptance.InitResourceCheck(assumeAgencyAssign, &obj, getRoleResourceFunc)

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: "3.3.0",
			},
		},
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcNormalAssign.CheckResourceDestroy(),
			rcAssumeRoleAssign.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccRole_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcNormalAssign.CheckResourceExists(),
					resource.TestCheckResourceAttr(normalAssign, "name", name+"_normal"),
					resource.TestCheckResourceAttr(normalAssign, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(normalAssign, "type", "AX"),
					rcAssumeRoleAssign.CheckResourceExists(),
					resource.TestCheckResourceAttr(assumeAgencyAssign, "name", name+"_assume_agency"),
					resource.TestCheckResourceAttr(assumeAgencyAssign, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(assumeAgencyAssign, "type", "AX"),
				),
			},
			{
				Config: testAccRole_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rcNormalAssign.CheckResourceExists(),
					resource.TestCheckResourceAttr(normalAssign, "name", updateName+"_normal"),
					resource.TestCheckResourceAttr(normalAssign, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(normalAssign, "type", "AX"),
					rcAssumeRoleAssign.CheckResourceExists(),
					resource.TestCheckResourceAttr(assumeAgencyAssign, "name", updateName+"_assume_agency"),
					resource.TestCheckResourceAttr(assumeAgencyAssign, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(assumeAgencyAssign, "type", "AX"),
				),
			},
			{
				ResourceName:      normalAssign,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      assumeAgencyAssign,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRole_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "random_uuid" "test" {
  count = 3
}

resource "huaweicloud_identity_role" "normal" {
  name        = "%[1]s_normal"
  description = "Created by terraform script"
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

resource "huaweicloud_identity_role" "assume_agency" {
  name        = "%[1]s_assume_agency"
  type        = "AX"
  description = "Created by terraform script"
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
          "/iam/agencies/${random_uuid.test[0].result}",
          "/iam/agencies/${random_uuid.test[1].result}"
        ]
      }
    }
  ]
}
EOF
}
`, name)
}

func testAccRole_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "random_uuid" "test" {
  count = 3
}

resource "huaweicloud_identity_role" "normal" {
  name        = "%[1]s_normal"
  type        = "AX"
  description = "Updated by terraform script"
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
            "%[2]s"
          ]
        }
      }
    }
  ]
}
EOF
}

resource "huaweicloud_identity_role" "assume_agency" {
  name        = "%[1]s_assume_agency"
  type        = "AX"
  description = "Updated by terraform script"
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
          "/iam/agencies/${random_uuid.test[1].result}",
          "/iam/agencies/${random_uuid.test[2].result}"
        ]
      }
    }
  ]
}
EOF
}
`, name, acceptance.HW_REGION_NAME)
}
