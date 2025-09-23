package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAppHdaConfigurations_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_workspace_app_hda_configurations.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byServerGroupId   = "data.huaweicloud_workspace_app_hda_configurations.filter_by_server_group_id"
		dcByServerGroupId = acceptance.InitDataSourceCheck(byServerGroupId)

		byServerName   = "data.huaweicloud_workspace_app_hda_configurations.filter_by_server_name"
		dcByServerName = acceptance.InitDataSourceCheck(byServerName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroupId(t)
			acceptance.TestAccPreCheckWorkspaceAppServerId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataAppHdaConfigurations_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "configurations.#", regexp.MustCompile(`^[0-9]+$`)),

					dcByServerGroupId.CheckResourceExists(),
					resource.TestCheckOutput("is_server_group_id_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(byServerGroupId, "configurations.#"),

					dcByServerName.CheckResourceExists(),
					resource.TestCheckOutput("is_server_name_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(byServerName, "configurations.#"),

					// Check item attributes if configurations exist
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.server_id"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.machine_name"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.maintain_status"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.server_name"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.server_group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.server_group_name"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.sid"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.session_count"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.current_version"),
				),
			},
		},
	})
}

func testDataAppHdaConfigurations_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_app_servers" "test" {
  server_id = "%[1]s"
}

locals {
  server_name = try(data.huaweicloud_workspace_app_servers.test.servers[0].name, "NOT_FOUND")
}

# Get all HDA configurations
data "huaweicloud_workspace_app_hda_configurations" "test" {}

# Filter by server group ID (using a test server group ID)
data "huaweicloud_workspace_app_hda_configurations" "filter_by_server_group_id" {
  server_group_id = "%[2]s"
}

locals {
  server_group_id_filter_result = [
    for v in data.huaweicloud_workspace_app_hda_configurations.filter_by_server_group_id.configurations[*].server_group_id :
      v == "%[2]s"
  ]
}

output "is_server_group_id_filter_useful" {
  value = length(local.server_group_id_filter_result) > 0 && alltrue(local.server_group_id_filter_result)
}

# Filter by server name (using a test server name)
data "huaweicloud_workspace_app_hda_configurations" "filter_by_server_name" {
  server_name = local.server_name
}

locals {
  server_name_filter_result = [
    for v in data.huaweicloud_workspace_app_hda_configurations.filter_by_server_name.configurations[*].server_name :
      v == local.server_name
  ]
}

output "is_server_name_filter_useful" {
  value = length(local.server_name_filter_result) > 0 && alltrue(local.server_name_filter_result)
}
`, acceptance.HW_WORKSPACE_APP_SERVER_ID, acceptance.HW_WORKSPACE_APP_SERVER_GROUP_ID)
}
