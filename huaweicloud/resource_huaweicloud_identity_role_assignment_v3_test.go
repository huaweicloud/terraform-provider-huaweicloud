package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/golangsdk/openstack/identity/v3/projects"

	"github.com/huaweicloud/golangsdk/pagination"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/identity/v3/groups"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3/roles"
)

func TestAccIdentityV3RoleAssignment_basic(t *testing.T) {
	var role roles.Role
	var group groups.Group
	var project projects.Project
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityV3RoleAssignmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccIdentityV3RoleAssignment_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3RoleAssignmentExists("huaweicloud_identity_role_assignment_v3.role_assignment_1", &role, &group, &project),
					resource.TestCheckResourceAttrPtr(
						"huaweicloud_identity_role_assignment_v3.role_assignment_1", "project_id", &project.ID),
					resource.TestCheckResourceAttrPtr(
						"huaweicloud_identity_role_assignment_v3.role_assignment_1", "group_id", &group.ID),
					resource.TestCheckResourceAttrPtr(
						"huaweicloud_identity_role_assignment_v3.role_assignment_1", "role_id", &role.ID),
				),
			},
		},
	})
}

func testAccCheckIdentityV3RoleAssignmentDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	identityClient, err := config.identityV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_identity_role_assignment_v3" {
			continue
		}

		_, err := roles.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Role assignment still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3RoleAssignmentExists(n string, role *roles.Role, group *groups.Group, project *projects.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		identityClient, err := config.identityV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud identity client: %s", err)
		}

		domainID, projectID, groupID, userID, roleID := extractRoleAssignmentID(rs.Primary.ID)

		var opts roles.ListAssignmentsOpts
		opts = roles.ListAssignmentsOpts{
			GroupID:        groupID,
			ScopeDomainID:  domainID,
			ScopeProjectID: projectID,
			UserID:         userID,
		}

		pager := roles.ListAssignments(identityClient, opts)
		var assignment roles.RoleAssignment

		err = pager.EachPage(func(page pagination.Page) (bool, error) {
			assignmentList, err := roles.ExtractRoleAssignments(page)
			if err != nil {
				return false, err
			}

			for _, a := range assignmentList {
				if a.ID == roleID {
					assignment = a
					return false, nil
				}
			}

			return true, nil
		})
		if err != nil {
			return err
		}

		p, err := projects.Get(identityClient, projectID).Extract()
		if err != nil {
			return fmt.Errorf("Project not found")
		}
		*project = *p
		g, err := groups.Get(identityClient, groupID).Extract()
		if err != nil {
			return fmt.Errorf("Group not found")
		}
		*group = *g
		r, err := roles.Get(identityClient, assignment.ID).Extract()
		if err != nil {
			return fmt.Errorf("Role not found")
		}
		*role = *r

		return nil
	}
}

const testAccIdentityV3RoleAssignment_basic = `
resource "huaweicloud_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "huaweicloud_identity_group_v3" "group_1" {
  name = "user_1"
}

data "huaweicloud_identity_role_v3" "role_1" {
  name = "secu_admin"
}

resource "huaweicloud_identity_role_assignment_v3" "role_assignment_1" {
  group_id = "${huaweicloud_identity_group_v3.group_1.id}"
  project_id = "${huaweicloud_identity_project_v3.project_1.id}"
  role_id = "${data.huaweicloud_identity_role_v3.role_1.id}"
}
`
