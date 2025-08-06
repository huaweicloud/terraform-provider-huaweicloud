package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test, please enable a service that connects to LocalAD and the corresponding OU is created.
func TestAccDataSourceAppImageServers_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_workspace_app_image_servers.test"
		rName      = acceptance.RandomAccResourceName()
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byServerId   = "data.huaweicloud_workspace_app_image_servers.filter_by_server_id"
		dcByServerId = acceptance.InitDataSourceCheck(byServerId)

		byName   = "data.huaweicloud_workspace_app_image_servers.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byEpsId   = "data.huaweicloud_workspace_app_image_servers.filter_by_eps_id"
		dcByEpsId = acceptance.InitDataSourceCheck(byEpsId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
			acceptance.TestAccPreCheckWorkspaceAppImageSpecCode(t)
			acceptance.TestAccPreCheckWorkspaceOUName(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAppImageServers_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "servers.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByServerId.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(byServerId, "servers.0.image_id"),
					resource.TestCheckResourceAttrSet(byServerId, "servers.0.image_type"),
					resource.TestCheckResourceAttrSet(byServerId, "servers.0.spce_code"),
					resource.TestCheckResourceAttrSet(byServerId, "servers.0.aps_server_group_id"),
					resource.TestCheckResourceAttrSet(byServerId, "servers.0.aps_server_id"),
					resource.TestCheckResourceAttrSet(byServerId, "servers.0.app_group_id"),
					resource.TestCheckResourceAttrSet(byServerId, "servers.0.status"),
					resource.TestCheckResourceAttr(byServerId, "servers.0.authorize_accounts.#", "1"),
					resource.TestCheckResourceAttr(byServerId, "servers.0.authorize_accounts.0.type", "USER"),
					resource.TestCheckResourceAttrSet(byServerId, "servers.0.authorize_accounts.0.domain"),
					resource.TestCheckOutput("is_server_id_filter_useful", "true"),
					resource.TestMatchResourceAttr(byServerId, "servers.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byServerId, "servers.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByEpsId.CheckResourceExists(),
					resource.TestCheckOutput("is_eps_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceAppImageServers_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_workspace_app_image_servers" "test" {
  depends_on = [huaweicloud_workspace_app_image_server.test]
}

locals {
  server_id   = huaweicloud_workspace_app_image_server.test.id
  server_name = huaweicloud_workspace_app_image_server.test.name
  eps_id      = huaweicloud_workspace_app_image_server.test.enterprise_project_id
}

data "huaweicloud_workspace_app_image_servers" "filter_by_server_id" {
  depends_on = [huaweicloud_workspace_app_image_server.test]

  server_id = local.server_id
}

locals {
  server_id_filter_result = [for v in data.huaweicloud_workspace_app_image_servers.filter_by_server_id.servers[*].id :
  v == local.server_id]
}

output "is_server_id_filter_useful" {
  value = length(local.server_id_filter_result) == 1 && alltrue(local.server_id_filter_result)
}

# Fuzzy search is supported.
data "huaweicloud_workspace_app_image_servers" "filter_by_name" {
  depends_on = [huaweicloud_workspace_app_image_server.test]

  name = local.server_name
}

locals {
  name_filter_result = [for v in data.huaweicloud_workspace_app_image_servers.filter_by_name.servers : strcontains(v.name, local.server_name)]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

data "huaweicloud_workspace_app_image_servers" "filter_by_eps_id" {
  depends_on = [huaweicloud_workspace_app_image_server.test]

  enterprise_project_id = local.eps_id
}

locals {
  eps_id_filter_result = [for v in data.huaweicloud_workspace_app_image_servers.filter_by_eps_id.servers[*].enterprise_project_id :
  v == local.eps_id]
}

output "is_eps_id_filter_useful" {
  value = length(local.eps_id_filter_result) > 0 && alltrue(local.eps_id_filter_result)
}
`, testResourceAppImageServer_withAD(name, "Data source test"))
}
