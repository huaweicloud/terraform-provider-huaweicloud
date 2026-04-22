package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	da "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dataarts"
)

func getSecurityResourcePermissionPolicyFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dataarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	return da.GetSecurityResourcePermissionPolicyById(client, state.Primary.Attributes["workspace_id"], state.Primary.ID)
}

// HW_USER_NAME cannot be the creator of the Workspace.
func TestAccSecurityResourcePermissionPolicy_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_dataarts_security_resource_permission_policy.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getSecurityResourcePermissionPolicyFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckUserName(t)
			acceptance.TestAccPreCheckDataArtsConnectionID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityResourcePermissionPolicy_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "resources.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(rName, "resources.*", map[string]string{
						"resource_id":   acceptance.HW_DATAARTS_CONNECTION_ID,
						"resource_type": "DATA_CONNECTION",
					}),
					resource.TestCheckTypeSetElemAttrPair(rName, "resources.*.resource_name",
						"data.huaweicloud_dataarts_studio_data_connections.test", "connections.0.name"),
					resource.TestCheckTypeSetElemAttrPair(rName, "resources.*.resource_id", "huaweicloud_identity_agency.test.0", "id"),
					resource.TestCheckTypeSetElemAttrPair(rName, "resources.*.resource_name", "huaweicloud_identity_agency.test.0", "name"),
					resource.TestCheckResourceAttr(rName, "members.#", "2"),
					resource.TestCheckTypeSetElemAttrPair(rName, "members.*.member_id", "huaweicloud_dataarts_studio_workspace_user.test", "id"),
					resource.TestCheckTypeSetElemNestedAttrs(rName, "members.*", map[string]string{
						"member_name": acceptance.HW_USER_NAME,
						"member_type": "USER",
					}),
					resource.TestCheckTypeSetElemAttrPair(rName, "members.*.member_id",
						"data.huaweicloud_dataarts_studio_workspace_user_roles.test", "roles.0.id"),
					resource.TestCheckTypeSetElemAttrPair(rName, "members.*.member_name",
						"data.huaweicloud_dataarts_studio_workspace_user_roles.test", "roles.0.name"),
					resource.TestMatchResourceAttr(rName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(rName, "create_user"),
				),
			},
			{
				Config: testAccSecurityResourcePermissionPolicy_basic_step2(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "resources.#", "2"),
					resource.TestCheckTypeSetElemAttrPair(rName, "resources.*.resource_id", "huaweicloud_identity_agency.test.0", "id"),
					resource.TestCheckTypeSetElemAttrPair(rName, "resources.*.resource_name", "huaweicloud_identity_agency.test.0", "name"),
					resource.TestCheckTypeSetElemAttrPair(rName, "resources.*.resource_id", "huaweicloud_identity_agency.test.1", "id"),
					resource.TestCheckTypeSetElemAttrPair(rName, "resources.*.resource_name", "huaweicloud_identity_agency.test.1", "name"),
					resource.TestCheckResourceAttr(rName, "resources.0.resource_type", "AGENCY"),
					resource.TestCheckResourceAttr(rName, "resources.1.resource_type", "AGENCY"),
					resource.TestCheckResourceAttr(rName, "members.#", "2"),
					resource.TestCheckTypeSetElemAttrPair(rName, "members.*.member_id",
						"data.huaweicloud_dataarts_studio_workspace_user_roles.test", "roles.0.id"),
					resource.TestCheckTypeSetElemAttrPair(rName, "members.*.member_name",
						"data.huaweicloud_dataarts_studio_workspace_user_roles.test", "roles.0.name"),
					resource.TestCheckTypeSetElemAttrPair(rName, "members.*.member_name",
						"data.huaweicloud_dataarts_studio_workspace_user_roles.test", "roles.1.name"),
					resource.TestCheckTypeSetElemAttrPair(rName, "members.*.member_id",
						"data.huaweicloud_dataarts_studio_workspace_user_roles.test", "roles.1.id"),
					resource.TestCheckResourceAttr(rName, "members.0.member_type", "WORKSPACE_ROLE"),
					resource.TestCheckResourceAttr(rName, "members.1.member_type", "WORKSPACE_ROLE"),
					resource.TestMatchResourceAttr(rName, "update_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSecurityResourcePermissionPolicyImportStateIDFunc(rName),
			},
		},
	})
}

func testAccSecurityResourcePermissionPolicyImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found in state", resourceName)
		}

		workspaceId := rs.Primary.Attributes["workspace_id"]
		policyId := rs.Primary.ID
		if workspaceId == "" || policyId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<workspace_id>/<id>', but got '%s/%s'", workspaceId, policyId)
		}
		return fmt.Sprintf("%s/%s", workspaceId, policyId), nil
	}
}

func testAccSecurityResourcePermissionPolicy_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_agency" "test" {
  count = 2

  name                   = "%[1]s${count.index}"
  delegated_service_name = "op_svc_dlg"
}

data "huaweicloud_dataarts_studio_workspace_user_roles" "test" {
  workspace_id = "%[2]s"
}

data "huaweicloud_identity_users" "test" {
  name = "%[3]s"
}

resource "huaweicloud_dataarts_studio_workspace_user" "test" {
  workspace_id = "%[2]s"
  user_id      = try(data.huaweicloud_identity_users.test.users[0].id, "NOT_FOUND")

  roles {
    id = try(data.huaweicloud_dataarts_studio_workspace_user_roles.test.roles[0].id, "NOT_FOUND")
  }
}

data "huaweicloud_dataarts_studio_data_connections" "test" {
  workspace_id  = "%[2]s"
  connection_id = "%[4]s"
}
`, name, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_USER_NAME, acceptance.HW_DATAARTS_CONNECTION_ID)
}

func testAccSecurityResourcePermissionPolicy_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_security_resource_permission_policy" "test" {
  workspace_id = "%[2]s"
  name         = "%[3]s"

  resources {
    resource_id   = try(data.huaweicloud_dataarts_studio_data_connections.test.connections[0].id, "NOT_FOUND")
    resource_name = try(data.huaweicloud_dataarts_studio_data_connections.test.connections[0].name, "NOT_FOUND")
    resource_type = "DATA_CONNECTION"
  }
  resources {
    resource_id   = try(huaweicloud_identity_agency.test[0].id, "NOT_FOUND")
    resource_name = try(huaweicloud_identity_agency.test[0].name, "NOT_FOUND")
    resource_type = "AGENCY"
  }

  members {
    member_id   = huaweicloud_dataarts_studio_workspace_user.test.id
    member_name = try(data.huaweicloud_identity_users.test.users[0].name, "NOT_FOUND")
    member_type = "USER"
  }
  members {
    member_id   = try(data.huaweicloud_dataarts_studio_workspace_user_roles.test.roles[0].id, "NOT_FOUND")
    member_name = try(data.huaweicloud_dataarts_studio_workspace_user_roles.test.roles[0].name, "NOT_FOUND")
    member_type = "WORKSPACE_ROLE"
  }
}
`, testAccSecurityResourcePermissionPolicy_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testAccSecurityResourcePermissionPolicy_basic_step2(name, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_security_resource_permission_policy" "test" {
  workspace_id = "%[2]s"
  name         = "%[3]s"

  dynamic "resources" {
    for_each = huaweicloud_identity_agency.test[*]

    content {
      resource_id   = resources.value.id
      resource_name = resources.value.name
      resource_type = "AGENCY"
    }
  }

  dynamic "members" {
    for_each = slice(data.huaweicloud_dataarts_studio_workspace_user_roles.test.roles, 0, 2)

    content {
      member_id   = members.value.id
      member_name = members.value.name
      member_type = "WORKSPACE_ROLE"
    }
  }
}
`, testAccSecurityResourcePermissionPolicy_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID, updateName)
}
