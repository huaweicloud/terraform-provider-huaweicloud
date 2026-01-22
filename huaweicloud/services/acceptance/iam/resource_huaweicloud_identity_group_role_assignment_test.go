package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3/roles"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iam"
)

func getGroupRoleAssignmentResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	identityClient, err := c.IdentityV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM v3 client: %s", err)
	}

	iamClient, err := c.IAMV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM v3.0 client: %s", err)
	}

	groupID := state.Primary.Attributes["group_id"]
	roleID := state.Primary.Attributes["role_id"]
	domainID := state.Primary.Attributes["domain_id"]
	projectID := state.Primary.Attributes["project_id"]
	enterpriseProjectID := state.Primary.Attributes["enterprise_project_id"]

	if domainID != "" {
		return iam.GetGroupRoleAssignmentWithDomainID(identityClient, groupID, roleID, domainID)
	}

	if projectID != "" {
		if projectID == "all" {
			specifiedRole := roles.Role{
				ID: roleID,
			}
			err = roles.CheckAllResourcesPermission(identityClient, c.DomainID, groupID, roleID).ExtractErr()
			return specifiedRole, err
		}

		return iam.GetGroupRoleAssignmentWithProjectID(identityClient, groupID, roleID, projectID)
	}

	if enterpriseProjectID != "" {
		return iam.GetGroupRoleAssignmentWithEpsID(iamClient, groupID, roleID, enterpriseProjectID)
	}

	return nil, golangsdk.ErrDefault404{}
}

func TestAccGroupRoleAssignment_basic(t *testing.T) {
	var (
		obj interface{}

		applyForDomain   = "huaweicloud_identity_group_role_assignment.apply_for_domain"
		rcApplyForDomain = acceptance.InitResourceCheck(applyForDomain, &obj, getGroupRoleAssignmentResourceFunc)

		applyForProject   = "huaweicloud_identity_group_role_assignment.apply_for_project"
		rcApplyForProject = acceptance.InitResourceCheck(applyForProject, &obj, getGroupRoleAssignmentResourceFunc)

		applyForAllProjects   = "huaweicloud_identity_group_role_assignment.apply_for_all_projects"
		rcApplyForAllProjects = acceptance.InitResourceCheck(applyForAllProjects, &obj, getGroupRoleAssignmentResourceFunc)

		applyForEpsID   = "huaweicloud_identity_group_role_assignment.apply_for_eps_id"
		rcApplyForEpsID = acceptance.InitResourceCheck(applyForEpsID, &obj, getGroupRoleAssignmentResourceFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcApplyForDomain.CheckResourceDestroy(),
			rcApplyForProject.CheckResourceDestroy(),
			rcApplyForAllProjects.CheckResourceDestroy(),
			rcApplyForEpsID.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccGroupRoleAssignment_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					// Assign role to a specified domain
					rcApplyForDomain.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(applyForDomain, "group_id", "huaweicloud_identity_group.test", "id"),
					resource.TestCheckResourceAttrPair(applyForDomain, "role_id", "huaweicloud_identity_role.test.0", "id"),
					resource.TestCheckResourceAttr(applyForDomain, "domain_id", acceptance.HW_DOMAIN_ID),
					// Assign role to a specified project
					rcApplyForProject.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(applyForProject, "group_id", "huaweicloud_identity_group.test", "id"),
					resource.TestCheckResourceAttrPair(applyForProject, "role_id", "huaweicloud_identity_role.test.1", "id"),
					resource.TestCheckResourceAttrPair(applyForProject, "project_id", "data.huaweicloud_identity_projects.test", "projects.0.id"),
					// Assign role to all projects
					rcApplyForAllProjects.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(applyForAllProjects, "group_id", "huaweicloud_identity_group.test", "id"),
					resource.TestCheckResourceAttrPair(applyForAllProjects, "role_id", "huaweicloud_identity_role.test.2", "id"),
					resource.TestCheckResourceAttr(applyForAllProjects, "project_id", "all"),
					// Assign role to a specified enterprise project
					rcApplyForEpsID.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(applyForEpsID, "group_id", "huaweicloud_identity_group.test", "id"),
					resource.TestCheckResourceAttrPair(applyForEpsID, "role_id", "huaweicloud_identity_role.test.3", "id"),
					resource.TestCheckResourceAttr(applyForEpsID, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				ResourceName:      applyForDomain,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGroupRoleAssignmentDomainImportStateFunc(applyForDomain),
			},
			{
				ResourceName:      applyForProject,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGroupRoleAssignmentProjectImportStateFunc(applyForProject),
			},
			{
				ResourceName:      applyForAllProjects,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGroupRoleAssignmentProjectImportStateFunc(applyForAllProjects),
			},
			{
				ResourceName:      applyForEpsID,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGroupRoleAssignmentEpsImportStateFunc(applyForEpsID),
			},
		},
	})
}

func testAccGroupRoleAssignmentDomainImportStateFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}
		if rs.Primary.Attributes["group_id"] == "" ||
			rs.Primary.Attributes["role_id"] == "" || rs.Primary.Attributes["domain_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID,"+
				" want '<group_id>/<role_id>/<domain_id>', but got '%s/%s/%s'",
				rs.Primary.Attributes["group_id"], rs.Primary.Attributes["role_id"],
				rs.Primary.Attributes["domain_id"])
		}
		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["group_id"],
			rs.Primary.Attributes["role_id"], rs.Primary.Attributes["domain_id"]), nil
	}
}

func testAccGroupRoleAssignmentProjectImportStateFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}
		if rs.Primary.Attributes["group_id"] == "" ||
			rs.Primary.Attributes["role_id"] == "" || rs.Primary.Attributes["project_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID,"+
				" want '<group_id>/<role_id>/<project_id>', but got '%s/%s/%s'",
				rs.Primary.Attributes["group_id"], rs.Primary.Attributes["role_id"],
				rs.Primary.Attributes["project_id"])
		}
		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["group_id"],
			rs.Primary.Attributes["role_id"], rs.Primary.Attributes["project_id"]), nil
	}
}

func testAccGroupRoleAssignmentEpsImportStateFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}
		if rs.Primary.Attributes["group_id"] == "" ||
			rs.Primary.Attributes["role_id"] == "" || rs.Primary.Attributes["enterprise_project_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID,"+
				" want '<group_id>/<role_id>/<enterprise_project_id>', but got '%s/%s/%s'",
				rs.Primary.Attributes["group_id"], rs.Primary.Attributes["role_id"],
				rs.Primary.Attributes["enterprise_project_id"])
		}
		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["group_id"],
			rs.Primary.Attributes["role_id"], rs.Primary.Attributes["enterprise_project_id"]), nil
	}
}

func testAccGroupRoleAssignment_basic_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_identity_projects" "test" {
  name = "%[1]s"
}

variable "obs_role_privileges" {
  type    = list(string)
  default = [
	"obs:bucket:CreateBucket",
	"obs:bucket:GetBucketAcl",
	"obs:bucket:ListBucket",
	"obs:bucket:HeadBucket",
  ]
}

resource "huaweicloud_identity_role" "test" {
  count = length(var.obs_role_privileges)

  name        = format("%[2]s_%%s", count.index)
  description = "Created by terraform script"
  type        = "AX"
  policy      = <<EOT
{
  "Version": "1.1",
  "Statement": [
    {
      "Action": [
         "${var.obs_role_privileges[count.index]}"
      ],
      "Effect": "Allow",
      "Resource": [
        "obs:*:*:bucket:*"
      ]
    }
  ]
}
EOT
}

resource "huaweicloud_identity_group" "test" {
  name        = "%[2]s"
  description = "Created by terraform script"
}
`, acceptance.HW_REGION_NAME, name)
}

func testAccGroupRoleAssignment_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_group_role_assignment" "apply_for_domain" {
  group_id  = huaweicloud_identity_group.test.id
  role_id   = huaweicloud_identity_role.test[0].id
  domain_id = "%[2]s"
}

resource "huaweicloud_identity_group_role_assignment" "apply_for_project" {
  group_id   = huaweicloud_identity_group.test.id
  role_id    = huaweicloud_identity_role.test[1].id
  project_id = try(data.huaweicloud_identity_projects.test.projects[0].id, "NOT_FOUND")
}

resource "huaweicloud_identity_group_role_assignment" "apply_for_all_projects" {
  group_id   = huaweicloud_identity_group.test.id
  role_id    = huaweicloud_identity_role.test[2].id
  project_id = "all"
}

resource "huaweicloud_identity_group_role_assignment" "apply_for_eps_id" {
  group_id              = huaweicloud_identity_group.test.id
  role_id               = huaweicloud_identity_role.test[3].id
  enterprise_project_id = "%[3]s"
}
`, testAccGroupRoleAssignment_basic_base(name),
		acceptance.HW_DOMAIN_ID,
		acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
