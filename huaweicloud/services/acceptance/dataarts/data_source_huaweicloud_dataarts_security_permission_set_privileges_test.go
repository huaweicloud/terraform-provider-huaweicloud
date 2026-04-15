package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSecurityPermissionSetPrivileges_basic(t *testing.T) {
	var (
		rName = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_security_permission_set_privileges.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byPrivilegeType   = "data.huaweicloud_dataarts_security_permission_set_privileges.filter_by_privilege_type"
		dcByPrivilegeType = acceptance.InitDataSourceCheck(byPrivilegeType)

		byPrivilegeAction   = "data.huaweicloud_dataarts_security_permission_set_privileges.filter_by_privilege_action"
		dcByPrivilegeAction = acceptance.InitDataSourceCheck(byPrivilegeAction)

		byClusterId   = "data.huaweicloud_dataarts_security_permission_set_privileges.filter_by_cluster_id"
		dcByClusterId = acceptance.InitDataSourceCheck(byClusterId)

		byClusterName   = "data.huaweicloud_dataarts_security_permission_set_privileges.filter_by_cluster_name"
		dcByClusterName = acceptance.InitDataSourceCheck(byClusterName)

		byDatabaseName   = "data.huaweicloud_dataarts_security_permission_set_privileges.filter_by_database_name"
		dcByDatabaseName = acceptance.InitDataSourceCheck(byDatabaseName)

		byTableName   = "data.huaweicloud_dataarts_security_permission_set_privileges.filter_by_table_name"
		dcByTableName = acceptance.InitDataSourceCheck(byTableName)

		byColumnName   = "data.huaweicloud_dataarts_security_permission_set_privileges.filter_by_column_name"
		dcByColumnName = acceptance.InitDataSourceCheck(byColumnName)

		bySyncStatus   = "data.huaweicloud_dataarts_security_permission_set_privileges.filter_by_sync_status"
		dcBySyncStatus = acceptance.InitDataSourceCheck(bySyncStatus)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsManagerID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSecurityPermissionSetPrivileges_nonExistentWorkspace(),
				ExpectError: regexp.MustCompile("error querying DataArts Security permission set privileges"),
			},
			{
				Config: testAccDataSecurityPermissionSetPrivileges_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "privileges.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByPrivilegeType.CheckResourceExists(),
					resource.TestCheckOutput("privilege_type_filter_is_useful", "true"),
					dcByPrivilegeAction.CheckResourceExists(),
					resource.TestCheckOutput("privilege_action_filter_is_useful", "true"),
					dcByClusterId.CheckResourceExists(),
					resource.TestCheckOutput("cluster_id_filter_is_useful", "true"),
					dcByClusterName.CheckResourceExists(),
					resource.TestCheckOutput("cluster_name_filter_is_useful", "true"),
					dcByDatabaseName.CheckResourceExists(),
					resource.TestCheckOutput("database_name_filter_is_useful", "true"),
					dcByTableName.CheckResourceExists(),
					resource.TestCheckOutput("table_name_filter_is_useful", "true"),
					dcByColumnName.CheckResourceExists(),
					resource.TestCheckOutput("column_name_filter_is_useful", "true"),
					dcBySyncStatus.CheckResourceExists(),
					resource.TestCheckOutput("sync_status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSecurityPermissionSetPrivileges_nonExistentWorkspace() string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
data "huaweicloud_dataarts_security_permission_set_privileges" "test" {
  workspace_id      = "%[1]s"
  permission_set_id = "%[1]s"
}
`, randUUID)
}

func testAccDataSecurityPermissionSetPrivileges_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_security_permission_set_privilege" "test" {
  workspace_id      = "%[2]s"
  permission_set_id = huaweicloud_dataarts_security_permission_set.test.id
  datasource_type   = "DLI"
  type              = "ALLOW"
  actions           = ["SELECT"]
  cluster_name      = "*"
  database_name     = huaweicloud_dli_database.test.name
  table_name        = huaweicloud_dli_table.test.name
  column_name       = huaweicloud_dli_table.test.columns[0].name
  connection_id     = var.data_connection_id != "" ? var.data_connection_id : huaweicloud_dataarts_studio_data_connection.test[0].id
}

# Query all privileges without any filters.
data "huaweicloud_dataarts_security_permission_set_privileges" "all" {
  depends_on = [
    huaweicloud_dataarts_security_permission_set_privilege.test,
  ]

  workspace_id      = "%[2]s"
  permission_set_id = huaweicloud_dataarts_security_permission_set.test.id
}

# Filter by permission type.
locals {
  privilege_type = huaweicloud_dataarts_security_permission_set_privilege.test.type
}

data "huaweicloud_dataarts_security_permission_set_privileges" "filter_by_privilege_type" {
  depends_on = [
    huaweicloud_dataarts_security_permission_set_privilege.test,
  ]

  workspace_id      = "%[2]s"
  permission_set_id = huaweicloud_dataarts_security_permission_set.test.id
  privilege_type    = local.privilege_type
}

locals {
  privilege_type_filter_result = [
    for v in data.huaweicloud_dataarts_security_permission_set_privileges.filter_by_privilege_type.privileges[*].type :
      v == local.privilege_type
  ]
}

output "privilege_type_filter_is_useful" {
  value = length(local.privilege_type_filter_result) > 0 && alltrue(local.privilege_type_filter_result)
}

# Filter by permission action.
locals {
  privilege_action = try(tolist(huaweicloud_dataarts_security_permission_set_privilege.test.actions)[0], "NOT_FOUND")
}

data "huaweicloud_dataarts_security_permission_set_privileges" "filter_by_privilege_action" {
  depends_on = [
    huaweicloud_dataarts_security_permission_set_privilege.test,
  ]

  workspace_id      = "%[2]s"
  permission_set_id = huaweicloud_dataarts_security_permission_set.test.id
  privilege_action  = local.privilege_action
}

locals {
  privilege_action_filter_result = [
    for v in data.huaweicloud_dataarts_security_permission_set_privileges.filter_by_privilege_action.privileges[*].actions :
      contains(v, local.privilege_action)
  ]
}

output "privilege_action_filter_is_useful" {
  value = length(local.privilege_action_filter_result) > 0 && alltrue(local.privilege_action_filter_result)
}

# Filter by cluster ID.
locals {
  cluster_id = huaweicloud_dataarts_security_permission_set_privilege.test.cluster_id
}

data "huaweicloud_dataarts_security_permission_set_privileges" "filter_by_cluster_id" {
  depends_on = [
    huaweicloud_dataarts_security_permission_set_privilege.test,
  ]

  workspace_id      = "%[2]s"
  permission_set_id = huaweicloud_dataarts_security_permission_set.test.id
  cluster_id        = local.cluster_id
}

locals {
  cluster_id_filter_result = [
    for v in data.huaweicloud_dataarts_security_permission_set_privileges.filter_by_cluster_id.privileges[*].cluster_id :
      v == local.cluster_id
  ]
}

output "cluster_id_filter_is_useful" {
  value = length(local.cluster_id_filter_result) > 0 && alltrue(local.cluster_id_filter_result)
}

# Filter by cluster name.
locals {
  cluster_name = huaweicloud_dataarts_security_permission_set_privilege.test.cluster_name
}

data "huaweicloud_dataarts_security_permission_set_privileges" "filter_by_cluster_name" {
  depends_on = [
    huaweicloud_dataarts_security_permission_set_privilege.test,
  ]

  workspace_id      = "%[2]s"
  permission_set_id = huaweicloud_dataarts_security_permission_set.test.id
  cluster_name      = local.cluster_name
}

locals {
  cluster_name_filter_result = [
    for v in data.huaweicloud_dataarts_security_permission_set_privileges.filter_by_cluster_name.privileges[*].cluster_name :
      v == local.cluster_name
  ]
}

output "cluster_name_filter_is_useful" {
  value = length(local.cluster_name_filter_result) > 0 && alltrue(local.cluster_name_filter_result)
}

# Filter by database name.
locals {
  database_name = huaweicloud_dataarts_security_permission_set_privilege.test.database_name
}

data "huaweicloud_dataarts_security_permission_set_privileges" "filter_by_database_name" {
  depends_on = [
    huaweicloud_dataarts_security_permission_set_privilege.test,
  ]

  workspace_id      = "%[2]s"
  permission_set_id = huaweicloud_dataarts_security_permission_set.test.id
  database_name     = local.database_name
}

locals {
  database_name_filter_result = [
    for v in data.huaweicloud_dataarts_security_permission_set_privileges.filter_by_database_name.privileges[*].database_name :
      v == local.database_name
  ]
}

output "database_name_filter_is_useful" {
  value = length(local.database_name_filter_result) > 0 && alltrue(local.database_name_filter_result)
}

# Filter by table name.
locals {
  table_name = huaweicloud_dataarts_security_permission_set_privilege.test.table_name
}

data "huaweicloud_dataarts_security_permission_set_privileges" "filter_by_table_name" {
  depends_on = [
    huaweicloud_dataarts_security_permission_set_privilege.test,
  ]

  workspace_id      = "%[2]s"
  permission_set_id = huaweicloud_dataarts_security_permission_set.test.id
  table_name        = local.table_name
}

locals {
  table_name_filter_result = [
    for v in data.huaweicloud_dataarts_security_permission_set_privileges.filter_by_table_name.privileges[*].table_name :
      v == local.table_name
  ]
}

output "table_name_filter_is_useful" {
  value = length(local.table_name_filter_result) > 0 && alltrue(local.table_name_filter_result)
}

# Filter by column name.
locals {
  column_name = huaweicloud_dataarts_security_permission_set_privilege.test.column_name
}

data "huaweicloud_dataarts_security_permission_set_privileges" "filter_by_column_name" {
  depends_on = [
    huaweicloud_dataarts_security_permission_set_privilege.test,
  ]

  workspace_id      = "%[2]s"
  permission_set_id = huaweicloud_dataarts_security_permission_set.test.id
  column_name       = local.column_name
}

locals {
  column_name_filter_result = [
    for v in data.huaweicloud_dataarts_security_permission_set_privileges.filter_by_column_name.privileges[*].column_name :
      v == local.column_name
  ]
}

output "column_name_filter_is_useful" {
  value = length(local.column_name_filter_result) > 0 && alltrue(local.column_name_filter_result)
}

# Filter by sync status.
locals {
  sync_status = huaweicloud_dataarts_security_permission_set_privilege.test.status
}

data "huaweicloud_dataarts_security_permission_set_privileges" "filter_by_sync_status" {
  depends_on = [
    huaweicloud_dataarts_security_permission_set_privilege.test,
  ]

  workspace_id      = "%[2]s"
  permission_set_id = huaweicloud_dataarts_security_permission_set.test.id
  sync_status       = local.sync_status
}

locals {
  sync_status_filter_result = [
    for v in data.huaweicloud_dataarts_security_permission_set_privileges.filter_by_sync_status.privileges[*].sync_status :
      v == local.sync_status
  ]
}

output "sync_status_filter_is_useful" {
  value = length(local.sync_status_filter_result) > 0 && alltrue(local.sync_status_filter_result)
}
`, testAccSecurityPermissionSetPrivilege_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}
