package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/eps_permissions"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iam"
)

func getIdentityUserRoleAssignmentResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.IAMV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM v3.0 client: %s", err)
	}

	userID := state.Primary.Attributes["user_id"]
	roleID := state.Primary.Attributes["role_id"]
	enterpriseProjectID := state.Primary.Attributes["enterprise_project_id"]

	return iam.GetUserRoleAssignmentWithEpsID(client, userID, roleID, enterpriseProjectID)
}

func TestAccIdentityUserRoleAssignment_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_identity_user_role_assignment.test"
	var role eps_permissions.Role

	rc := acceptance.InitResourceCheck(
		resourceName,
		&role,
		getIdentityUserRoleAssignmentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityUserRoleAssignment_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrPair(resourceName, "user_id",
						"huaweicloud_identity_user.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "role_id",
						"huaweicloud_identity_role.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccIdentityUserRoleAssignmentImportStateFunc(resourceName),
			},
		},
	})
}

func testAccIdentityUserRoleAssignmentImportStateFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}
		if rs.Primary.Attributes["user_id"] == "" ||
			rs.Primary.Attributes["role_id"] == "" || rs.Primary.Attributes["enterprise_project_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID,"+
				" want '<user_id>/<role_id>/<enterprise_project_id>', but got '%s/%s/%s'",
				rs.Primary.Attributes["user_id"], rs.Primary.Attributes["role_id"],
				rs.Primary.Attributes["enterprise_project_id"])
		}
		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["user_id"],
			rs.Primary.Attributes["role_id"], rs.Primary.Attributes["enterprise_project_id"]), nil
	}
}

func testAccIdentityUserRoleAssignment_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_role" test {
  name        = "%[1]s"
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

resource "huaweicloud_identity_user" "test" {
  name        = "%[1]s"
  description = "A user"
  password    = "Test@12345678"
}`, rName)
}

func testAccIdentityUserRoleAssignment_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_identity_user_role_assignment" "test" {
  user_id               = huaweicloud_identity_user.test.id
  role_id               = huaweicloud_identity_role.test.id
  enterprise_project_id = "%s"
}
`, testAccIdentityUserRoleAssignment_base(rName), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
