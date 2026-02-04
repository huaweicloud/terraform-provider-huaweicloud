package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iam"
)

func getV3GroupRoleAssignmentResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		err            error
		identityClient *golangsdk.ServiceClient
		iamClient      *golangsdk.ServiceClient

		groupId   = state.Primary.Attributes["group_id"]
		roleId    = state.Primary.Attributes["role_id"]
		domainId  = state.Primary.Attributes["domain_id"]
		projectId = state.Primary.Attributes["project_id"]
		epsId     = state.Primary.Attributes["enterprise_project_id"]
	)

	identityClient, err = c.IdentityV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM v3 client: %s", err)
	}
	iamClient, err = c.IAMV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM v3.0 client: %s", err)
	}

	if domainId != "" {
		err = iam.CheckV3GroupRoleAssignmentWithDomainId(identityClient, groupId, roleId, domainId)
	}
	if projectId != "" {
		err = iam.CheckV3GroupRoleAssignmentWithProjectId(identityClient, groupId, roleId, c.DomainID, projectId)
	}
	if epsId != "" {
		err = iam.CheckV3GroupRoleAssignmentWithEpsId(iamClient, groupId, roleId, epsId)
	}

	if err != nil {
		return nil, err
	}
	return roleId, nil
}

func TestAccV3GroupRoleAssignment_basic(t *testing.T) {
	var (
		obj interface{}

		applyForDomain   = "huaweicloud_identity_group_role_assignment.apply_for_domain"
		rcApplyForDomain = acceptance.InitResourceCheck(applyForDomain, &obj, getV3GroupRoleAssignmentResourceFunc)

		applyForProject   = "huaweicloud_identity_group_role_assignment.apply_for_project"
		rcApplyForProject = acceptance.InitResourceCheck(applyForProject, &obj, getV3GroupRoleAssignmentResourceFunc)

		applyForAllProjects   = "huaweicloud_identity_group_role_assignment.apply_for_all_projects"
		rcApplyForAllProjects = acceptance.InitResourceCheck(applyForAllProjects, &obj, getV3GroupRoleAssignmentResourceFunc)

		applyForEpsId   = "huaweicloud_identity_group_role_assignment.apply_for_eps_id"
		rcApplyForEpsId = acceptance.InitResourceCheck(applyForEpsId, &obj, getV3GroupRoleAssignmentResourceFunc)

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
			rcApplyForEpsId.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccV3GroupRoleAssignment_basic_step1(name),
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
					rcApplyForEpsId.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(applyForEpsId, "group_id", "huaweicloud_identity_group.test", "id"),
					resource.TestCheckResourceAttrPair(applyForEpsId, "role_id", "huaweicloud_identity_role.test.3", "id"),
					resource.TestCheckResourceAttr(applyForEpsId, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				ResourceName:      applyForDomain,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccV3GroupRoleAssignmentImportStateFunc(applyForDomain, "domain"),
			},
			{
				ResourceName:      applyForProject,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccV3GroupRoleAssignmentImportStateFunc(applyForProject, "project"),
			},
			{
				ResourceName:      applyForAllProjects,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccV3GroupRoleAssignmentImportStateFunc(applyForAllProjects, "project"),
			},
			{
				ResourceName:      applyForEpsId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccV3GroupRoleAssignmentImportStateFunc(applyForEpsId, "enterprise_project"),
			},
			// Legacy import
			{
				ResourceName:      applyForDomain,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccV3GroupRoleAssignmentDomainImportStateFunc(applyForDomain),
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
				ResourceName:      applyForEpsId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGroupRoleAssignmentEpsImportStateFunc(applyForEpsId),
			},
		},
	})
}

func testAccV3GroupRoleAssignmentImportStateFunc(resourceName, assignmentType string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		var err error
		switch assignmentType {
		case "domain":
			if rs.Primary.Attributes["group_id"] == "" ||
				rs.Primary.Attributes["role_id"] == "" || rs.Primary.Attributes["domain_id"] == "" {
				return "", fmt.Errorf("invalid format specified for import ID,"+
					" want '<group_id>/<role_id>/<domain_id>:domain', but got '%s/%s/%s:domain'",
					rs.Primary.Attributes["group_id"], rs.Primary.Attributes["role_id"],
					rs.Primary.Attributes["domain_id"])
			}
			return fmt.Sprintf("%s/%s/%s:domain", rs.Primary.Attributes["group_id"], rs.Primary.Attributes["role_id"],
				rs.Primary.Attributes["domain_id"]), nil
		case "project":
			if rs.Primary.Attributes["group_id"] == "" ||
				rs.Primary.Attributes["role_id"] == "" || rs.Primary.Attributes["project_id"] == "" {
				return "", fmt.Errorf("invalid format specified for import ID,"+
					" want '<group_id>/<role_id>/<project_id>:project', but got '%s/%s/%s:project'",
					rs.Primary.Attributes["group_id"], rs.Primary.Attributes["role_id"],
					rs.Primary.Attributes["project_id"])
			}
			return fmt.Sprintf("%s/%s/%s:project", rs.Primary.Attributes["group_id"], rs.Primary.Attributes["role_id"],
				rs.Primary.Attributes["project_id"]), nil
		case "enterprise_project":
			if rs.Primary.Attributes["group_id"] == "" ||
				rs.Primary.Attributes["role_id"] == "" || rs.Primary.Attributes["enterprise_project_id"] == "" {
				return "", fmt.Errorf("invalid format specified for import ID,"+
					" want '<group_id>/<role_id>/<enterprise_project_id>:enterprise_project', but got '%s/%s/%s:enterprise_project'",
					rs.Primary.Attributes["group_id"], rs.Primary.Attributes["role_id"],
					rs.Primary.Attributes["enterprise_project_id"])
			}
			return fmt.Sprintf("%s/%s/%s:enterprise_project", rs.Primary.Attributes["group_id"], rs.Primary.Attributes["role_id"],
				rs.Primary.Attributes["enterprise_project_id"]), nil
		default:
			err = fmt.Errorf(`invalid format specified for import ID, want these following format:
1. <group_id>/<role_id>/<domain_id>:domain
2. <group_id>/<role_id>/all:project
3. <group_id>/<role_id>/<project_id>:project
4. <group_id>/<role_id>/<domain_id>:enterprise_project
but got '%s:%s'`, rs.Primary.ID, assignmentType)
		}
		return "", err
	}
}

func testAccV3GroupRoleAssignmentDomainImportStateFunc(resourceName string) resource.ImportStateIdFunc {
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

func testAccV3GroupRoleAssignment_basic_base(name string) string {
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

func testAccV3GroupRoleAssignment_basic_step1(name string) string {
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
`, testAccV3GroupRoleAssignment_basic_base(name),
		acceptance.HW_DOMAIN_ID,
		acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
