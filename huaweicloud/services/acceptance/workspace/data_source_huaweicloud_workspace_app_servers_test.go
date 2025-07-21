package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAppServers_basic(t *testing.T) {
	var (
		all  = "data.huaweicloud_workspace_app_servers.test"
		dc   = acceptance.InitDataSourceCheck(all)
		name = acceptance.RandomAccResourceName()

		byServerName   = "data.huaweicloud_workspace_app_servers.filter_by_server_name"
		dcByServerName = acceptance.InitDataSourceCheck(byServerName)

		byMachineName   = "data.huaweicloud_workspace_app_servers.filter_by_machine_name"
		dcByMachineName = acceptance.InitDataSourceCheck(byMachineName)

		byServerGroupId   = "data.huaweicloud_workspace_app_servers.filter_by_server_group_id"
		dcByServerGroupId = acceptance.InitDataSourceCheck(byServerGroupId)

		byMaintainStatus   = "data.huaweicloud_workspace_app_servers.filter_by_maintain_status"
		dcByMaintainStatus = acceptance.InitDataSourceCheck(byMaintainStatus)

		byScalingAutoCreate   = "data.huaweicloud_workspace_app_servers.filter_by_scaling_auto_create"
		dcByScalingAutoCreate = acceptance.InitDataSourceCheck(byScalingAutoCreate)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAppServers_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "servers.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckResourceAttrSet(all, "servers.0.id"),
					resource.TestCheckResourceAttrSet(all, "servers.0.name"),
					resource.TestCheckResourceAttrSet(all, "servers.0.machine_name"),
					resource.TestCheckResourceAttrSet(all, "servers.0.server_group_id"),
					resource.TestCheckResourceAttrSet(all, "servers.0.status"),
					resource.TestCheckResourceAttrSet(all, "servers.0.flavor.#"),
					resource.TestCheckResourceAttrSet(all, "servers.0.product_info.#"),
					resource.TestCheckResourceAttrSet(all, "servers.0.host_address.#"),
					resource.TestCheckResourceAttrSet(all, "servers.0.tags.#"),
					dcByServerName.CheckResourceExists(),
					resource.TestCheckOutput("is_server_name_filter_useful", "true"),
					dcByMachineName.CheckResourceExists(),
					resource.TestCheckOutput("is_machine_name_filter_useful", "true"),
					dcByServerGroupId.CheckResourceExists(),
					resource.TestCheckOutput("is_server_group_id_filter_useful", "true"),
					dcByMaintainStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_maintain_status_filter_useful", "true"),
					dcByScalingAutoCreate.CheckResourceExists(),
					resource.TestCheckOutput("is_scaling_auto_create_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceAppServers_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_app_server_groups" "test" {
  server_group_id = "%[1]s"
}

data "huaweicloud_vpc_subnets" "test" {
  id = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].subnet_id, null)
}

resource "huaweicloud_workspace_app_server" "test" {
  name            = "%[2]s"
  server_group_id = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].id, null)
  type            = "createApps"
  flavor_id       = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].product_id, null)
  vpc_id          = try(data.huaweicloud_vpc_subnets.test.subnets[0].vpc_id, null)
  subnet_id       = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].subnet_id, null)
  maintain_status = false

  root_volume {
    type = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].system_disk_type, null)
    size = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].system_disk_size, null)
  }
}
`, acceptance.HW_WORKSPACE_APP_SERVER_GROUP_ID, name)
}

func testDataSourceAppServers_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_workspace_app_servers" "test" {
  server_id = huaweicloud_workspace_app_server.test.id
}

# Filter by server name
locals {
  server_name = data.huaweicloud_workspace_app_servers.test.servers[0].name
}

data "huaweicloud_workspace_app_servers" "filter_by_server_name" {
  depends_on = [
    data.huaweicloud_workspace_app_servers.test
  ]

  server_name = local.server_name
}

locals {
  server_name_filter_result = [
    for v in data.huaweicloud_workspace_app_servers.filter_by_server_name.servers[*].name : v == local.server_name
  ]
}

output "is_server_name_filter_useful" {
  value = length(local.server_name_filter_result) > 0 && alltrue(local.server_name_filter_result)
}

# Filter by machine name
locals {
  machine_name = data.huaweicloud_workspace_app_servers.test.servers[0].machine_name
}

data "huaweicloud_workspace_app_servers" "filter_by_machine_name" {
  depends_on = [
    data.huaweicloud_workspace_app_servers.test
  ]

  machine_name = local.machine_name
}

locals {
  machine_name_filter_result = [
    for v in data.huaweicloud_workspace_app_servers.filter_by_machine_name.servers[*].machine_name : v == local.machine_name
  ]
}

output "is_machine_name_filter_useful" {
  value = length(local.machine_name_filter_result) > 0 && alltrue(local.machine_name_filter_result)
}

# Filter by server group ID
locals {
  server_group_id = data.huaweicloud_workspace_app_servers.test.servers[0].server_group_id
}

data "huaweicloud_workspace_app_servers" "filter_by_server_group_id" {
  depends_on = [
    data.huaweicloud_workspace_app_servers.test
  ]

  server_group_id = local.server_group_id
}

locals {
  server_group_id_filter_result = [
    for v in data.huaweicloud_workspace_app_servers.filter_by_server_group_id.servers[*].server_group_id : v == local.server_group_id
  ]
}

output "is_server_group_id_filter_useful" {
  value = length(local.server_group_id_filter_result) > 0 && alltrue(local.server_group_id_filter_result)
}

# Filter by maintain status
data "huaweicloud_workspace_app_servers" "filter_by_maintain_status" {
  depends_on = [
    data.huaweicloud_workspace_app_servers.test
  ]

  maintain_status = false
}

locals {
  maintain_status = data.huaweicloud_workspace_app_servers.test.servers[0].maintain_status

  maintain_status_filter_result = [
    for v in data.huaweicloud_workspace_app_servers.filter_by_maintain_status.servers[*].maintain_status : v == false
  ]
}

output "is_maintain_status_filter_useful" {
  value = length(local.maintain_status_filter_result) > 0 && alltrue(local.maintain_status_filter_result)
}

# Filter by scaling auto create
data "huaweicloud_workspace_app_servers" "filter_by_scaling_auto_create" {
  depends_on = [
    data.huaweicloud_workspace_app_servers.test
  ]

  scaling_auto_create = false
}

locals {
  scaling_auto_create = data.huaweicloud_workspace_app_servers.test.servers[0].scaling_auto_create

  scaling_auto_create_filter_result = [
    for v in data.huaweicloud_workspace_app_servers.filter_by_scaling_auto_create.servers[*].scaling_auto_create : v == false
  ]
}

output "is_scaling_auto_create_filter_useful" {
  value = length(local.scaling_auto_create_filter_result) > 0 && alltrue(local.scaling_auto_create_filter_result)
}
`, testDataSourceAppServers_base(name))
}
