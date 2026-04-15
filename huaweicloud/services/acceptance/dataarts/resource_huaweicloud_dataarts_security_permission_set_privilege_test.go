package dataarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dataarts"
)

func getSecurityPermissionSetPrivilegeFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dataarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	workspaceId := state.Primary.Attributes["workspace_id"]
	permissionSetId := state.Primary.Attributes["permission_set_id"]
	return dataarts.GetPrivilegeById(client, workspaceId, permissionSetId, state.Primary.ID)
}

// Please setting HW_DATAARTS_CONNECTION_ID if a existing connection is prepared.
func TestAccSecurityPermissionSetPrivilege_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_dataarts_security_permission_set_privilege.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getSecurityPermissionSetPrivilegeFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsManagerID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityPermissionSetPrivilege_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttrPair(resourceName, "permission_set_id", "huaweicloud_dataarts_security_permission_set.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "datasource_type", "DLI"),
					resource.TestCheckResourceAttr(resourceName, "type", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", "*"),
					resource.TestCheckResourceAttr(resourceName, "cluster_id", "DLI"),
					resource.TestCheckResourceAttr(resourceName, "database_name", name),
					resource.TestCheckResourceAttr(resourceName, "table_name", name),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				Config: testAccSecurityPermissionSetPrivilege_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccSecurityPermissionSetPrivilegeImportFunc(resourceName),
				ImportStateVerifyIgnore: []string{"connection_id"},
			},
		},
	})
}

func testAccSecurityPermissionSetPrivilegeImportFunc(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", n, rs)
		}
		workspaceId := rs.Primary.Attributes["workspace_id"]
		permissionSetId := rs.Primary.Attributes["permission_set_id"]
		privilegeId := rs.Primary.ID
		if workspaceId == "" || permissionSetId == "" || privilegeId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<workspace_id>/<permission_set_id>/<id>', but got '%s/%s/%s'",
				workspaceId, permissionSetId, privilegeId)
		}
		return fmt.Sprintf("%s/%s/%s", workspaceId, permissionSetId, privilegeId), nil
	}
}

func testAccSecurityPermissionSetPrivilege_base(name string) string {
	return fmt.Sprintf(`
variable "data_connection_id" {
  type    = string
  default = "%[1]s"
}

variable "enterprise_project_id" {
  type    = string
  default = "%[2]s"
}

resource "huaweicloud_dataarts_studio_data_connection" "test" {
  count = var.data_connection_id == "" ? 1 : 0

  workspace_id = "%[3]s"
  type         = "DLI"
  name         = "%[4]s"
  env_type     = 0
  config       = jsonencode({
    "cdm_property_enable": "false"
  })

  lifecycle {
    ignore_changes = [
      config,
    ]
  }
}

resource "huaweicloud_dataarts_security_permission_set" "test" {
  workspace_id = "%[3]s"
  name         = "%[4]s"
  parent_id    = "0"
  manager_id   = "%[5]s"
}

resource "huaweicloud_dli_database" "test" {
  name                  = "%[4]s"
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_dli_table" "test" {
  depends_on = [
    huaweicloud_dli_database.test,
  ]

  database_name = huaweicloud_dli_database.test.name
  name          = "%[4]s"
  data_location = "DLI"

  columns {
    name        = "name"
    type        = "string"
    description = "person name"
  }
}
`, acceptance.HW_DATAARTS_CONNECTION_ID, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, acceptance.HW_DATAARTS_WORKSPACE_ID,
		name, acceptance.HW_DATAARTS_MANAGER_ID)
}

func testAccSecurityPermissionSetPrivilege_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_security_permission_set_privilege" "test" {
  depends_on = [
    huaweicloud_dli_table.test,
  ]

  workspace_id      = "%[2]s"
  permission_set_id = huaweicloud_dataarts_security_permission_set.test.id
  datasource_type   = "DLI"
  type              = "ALLOW"
  actions           = ["SELECT", "INSERT", "ALTER"]
  cluster_name      = "*"
  database_name     = huaweicloud_dli_database.test.name
  table_name        = huaweicloud_dli_table.test.name
  connection_id     = var.data_connection_id != "" ? var.data_connection_id : huaweicloud_dataarts_studio_data_connection.test[0].id
}
`, testAccSecurityPermissionSetPrivilege_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}

func testAccSecurityPermissionSetPrivilege_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_security_permission_set_privilege" "test" {
  depends_on = [
    huaweicloud_dli_table.test,
  ]

  workspace_id      = "%[2]s"
  permission_set_id = huaweicloud_dataarts_security_permission_set.test.id
  datasource_type   = "DLI"
  type              = "ALLOW"
  actions           = ["SELECT", "DROP"]
  cluster_name      = "*"
  database_name     = huaweicloud_dli_database.test.name
  table_name        = huaweicloud_dli_table.test.name
  connection_id     = var.data_connection_id != "" ? var.data_connection_id : huaweicloud_dataarts_studio_data_connection.test[0].id
}
`, testAccSecurityPermissionSetPrivilege_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}
