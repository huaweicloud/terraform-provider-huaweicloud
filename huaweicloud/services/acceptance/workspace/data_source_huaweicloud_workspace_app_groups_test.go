package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test, please create a workspace APP server group with SESSION_DESKTOP_APP type.
func TestAccDataSourceAppGroups_basic(t *testing.T) {
	var (
		rName      = acceptance.RandomAccResourceName()
		dataSource = "data.huaweicloud_workspace_app_groups.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byServerGroupId   = "data.huaweicloud_workspace_app_groups.filter_by_server_group_id"
		dcByServerGroupId = acceptance.InitDataSourceCheck(byServerGroupId)

		byId   = "data.huaweicloud_workspace_app_groups.filter_by_id"
		dcById = acceptance.InitDataSourceCheck(byId)

		byName   = "data.huaweicloud_workspace_app_groups.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byType   = "data.huaweicloud_workspace_app_groups.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceWorkspaceAppGroups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "groups.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByServerGroupId.CheckResourceExists(),
					resource.TestCheckOutput("is_server_group_id_filter_useful", "true"),
					dcById.CheckResourceExists(),
					resource.TestCheckOutput("is_group_id_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.description"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.server_group_id"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.server_group_name"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.server_group_name"),
					resource.TestCheckResourceAttrSet(byId, "groups.0.app_count"),
					resource.TestMatchResourceAttr(byId, "groups.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceAppGroups_base(name string, appType ...string) string {
	actAppType := "SESSION_DESKTOP_APP"
	if len(appType) > 0 {
		actAppType = appType[0]
	}

	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_group" "test" {
  server_group_id = "%[1]s"
  name            = "%[2]s"
  type            = "%[3]s"
  description     = "Created by terraform script"
}
`, acceptance.HW_WORKSPACE_APP_SERVER_GROUP_ID, name, actAppType)
}

func testDataSourceDataSourceWorkspaceAppGroups_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_workspace_app_groups" "test" {
  depends_on = [huaweicloud_workspace_app_group.test]
}

locals {
  server_group_id = huaweicloud_workspace_app_group.test.server_group_id
  group_id        = huaweicloud_workspace_app_group.test.id
  group_name      = huaweicloud_workspace_app_group.test.name
  group_type      = huaweicloud_workspace_app_group.test.type
}

data "huaweicloud_workspace_app_groups" "filter_by_server_group_id" {
  depends_on = [huaweicloud_workspace_app_group.test]

  server_group_id = local.server_group_id
}

locals {
  server_group_id_result = [for v in data.huaweicloud_workspace_app_groups.filter_by_server_group_id.groups[*].server_group_id :
  v == local.server_group_id]
}

output "is_server_group_id_filter_useful" {
  value = length(local.server_group_id_result) > 0 && alltrue(local.server_group_id_result)
}

data "huaweicloud_workspace_app_groups" "filter_by_id" {
  depends_on = [huaweicloud_workspace_app_group.test]

  group_id = local.group_id
}

locals {
  group_id_result = [for v in data.huaweicloud_workspace_app_groups.filter_by_id.groups[*].id : v == local.group_id]
}

output "is_group_id_filter_useful" {
  value = length(local.group_id_result) > 0 && alltrue(local.group_id_result)
}

# Fuzzy search is supported.
data "huaweicloud_workspace_app_groups" "filter_by_name" {
  depends_on = [huaweicloud_workspace_app_group.test]

  name = local.group_name
}

locals {
  name_result = [for v in data.huaweicloud_workspace_app_groups.filter_by_name.groups[*].name : strcontains(v, local.group_name)]
}

output "is_name_filter_useful" {
  value = length(local.name_result) > 0 && alltrue(local.name_result)
}

data "huaweicloud_workspace_app_groups" "filter_by_type" {
  depends_on = [huaweicloud_workspace_app_group.test]

  type = local.group_type
}

locals {
  type_result = [for v in data.huaweicloud_workspace_app_groups.filter_by_type.groups[*].type : v == local.group_type]
}

output "is_type_filter_useful" {
  value = length(local.type_result) > 0 && alltrue(local.type_result)
}
`, testDataSourceAppGroups_base(name))
}
