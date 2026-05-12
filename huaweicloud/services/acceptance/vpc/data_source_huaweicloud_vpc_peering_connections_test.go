package vpc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataPeeringConnections_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_vpc_peering_connections.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byConnectionId   = "data.huaweicloud_vpc_peering_connections.filter_by_connection_id"
		dcByConnectionId = acceptance.InitDataSourceCheck(byConnectionId)

		byName   = "data.huaweicloud_vpc_peering_connections.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byStatus   = "data.huaweicloud_vpc_peering_connections.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byProjectId   = "data.huaweicloud_vpc_peering_connections.filter_by_project_id"
		dcByProjectId = acceptance.InitDataSourceCheck(byProjectId)

		byVpcId   = "data.huaweicloud_vpc_peering_connections.filter_by_vpc_id"
		dcByVpcId = acceptance.InitDataSourceCheck(byVpcId)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPeeringConnections_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "connections.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'connection_id' parameter.
					dcByConnectionId.CheckResourceExists(),
					resource.TestCheckOutput("is_connection_id_filter_useful", "true"),
					resource.TestCheckResourceAttr(byConnectionId, "connections.#", "1"),
					resource.TestCheckResourceAttrSet(byConnectionId, "connections.0.id"),
					resource.TestCheckResourceAttrSet(byConnectionId, "connections.0.name"),
					resource.TestCheckResourceAttrSet(byConnectionId, "connections.0.status"),
					resource.TestCheckResourceAttr(byConnectionId, "connections.0.request_vpc_info.#", "1"),
					resource.TestCheckResourceAttrPair(byConnectionId, "connections.0.request_vpc_info.0.vpc_id",
						"huaweicloud_vpc.test.0", "id"),
					resource.TestCheckResourceAttrPair(byConnectionId, "connections.0.request_vpc_info.0.project_id",
						"data.huaweicloud_identity_projects.test", "projects.0.id"),
					resource.TestCheckResourceAttr(byConnectionId, "connections.0.accept_vpc_info.#", "1"),
					resource.TestCheckResourceAttrPair(byConnectionId, "connections.0.accept_vpc_info.0.vpc_id",
						"huaweicloud_vpc.test.1", "id"),
					resource.TestCheckResourceAttrPair(byConnectionId, "connections.0.accept_vpc_info.0.project_id",
						"data.huaweicloud_identity_projects.test", "projects.0.id"),
					resource.TestCheckResourceAttrSet(byConnectionId, "connections.0.description"),
					resource.TestMatchResourceAttr(byConnectionId, "connections.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byConnectionId, "connections.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(Z|([+-]\d{2}:\d{2}))$`)),
					// Filter by 'name' parameter.
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					// Filter by 'status' parameter.
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					// Filter by 'project_id' parameter.
					dcByProjectId.CheckResourceExists(),
					resource.TestCheckOutput("is_project_id_filter_useful", "true"),
					// Filter by 'vpc_id' parameter.
					dcByVpcId.CheckResourceExists(),
					resource.TestCheckOutput("is_vpc_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataPeeringConnections_basic_base(name string) string {
	return fmt.Sprintf(`
locals {
  basic_cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc" "test" {
  count = 2

  name = format("%[1]s_%%d", count.index)
  cidr = cidrsubnet(local.basic_cidr, 4, count.index)
}

resource "huaweicloud_vpc_peering_connection" "test" {
  name        = "%[1]s"
  vpc_id      = huaweicloud_vpc.test[0].id
  peer_vpc_id = huaweicloud_vpc.test[1].id
  description = "Created by acceptance test"
}
`, name)
}

func testAccDataPeeringConnections_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameter.
data "huaweicloud_vpc_peering_connections" "all" {
  depends_on = [huaweicloud_vpc_peering_connection.test]
}

# Filter by 'connection_id' parameter.
locals {
  connection_id = huaweicloud_vpc_peering_connection.test.id
}

data "huaweicloud_vpc_peering_connections" "filter_by_connection_id" {
  depends_on = [huaweicloud_vpc_peering_connection.test]

  connection_id = huaweicloud_vpc_peering_connection.test.id
}

locals {
  connection_id_filter_result = [
    for v in data.huaweicloud_vpc_peering_connections.filter_by_connection_id.connections[*].id : v == local.connection_id
  ]
}

output "is_connection_id_filter_useful" {
  value = length(local.connection_id_filter_result) > 0 && alltrue(local.connection_id_filter_result)
}

# Filter by 'name' parameter.
locals {
  connection_name = huaweicloud_vpc_peering_connection.test.name
}

data "huaweicloud_vpc_peering_connections" "filter_by_name" {
  depends_on = [huaweicloud_vpc_peering_connection.test]

  name = local.connection_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_vpc_peering_connections.filter_by_name.connections[*].name : v == local.connection_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by 'status' parameter.
locals {
  connection_status = huaweicloud_vpc_peering_connection.test.status
}

data "huaweicloud_vpc_peering_connections" "filter_by_status" {
  depends_on = [huaweicloud_vpc_peering_connection.test]

  status = local.connection_status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_vpc_peering_connections.filter_by_status.connections[*].status : v == local.connection_status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

# Filter by 'project_id' parameter.
data "huaweicloud_identity_projects" "test" {
  name = "%[2]s"
}

locals {
  project_id = data.huaweicloud_identity_projects.test.projects[0].id
}

data "huaweicloud_vpc_peering_connections" "filter_by_project_id" {
  depends_on = [huaweicloud_vpc_peering_connection.test]

  project_id = local.project_id
}

locals {
  project_id_filter_result = [
    for v in data.huaweicloud_vpc_peering_connections.filter_by_project_id.connections[*].accept_vpc_info[0].project_id : v == local.project_id
  ]
}

output "is_project_id_filter_useful" {
  value = length(local.project_id_filter_result) > 0 && alltrue(local.project_id_filter_result)
}

# Filter by 'vpc_id' parameter.
locals {
  vpc_id = huaweicloud_vpc_peering_connection.test.vpc_id
}

data "huaweicloud_vpc_peering_connections" "filter_by_vpc_id" {
  depends_on = [huaweicloud_vpc_peering_connection.test]

  vpc_id = local.vpc_id
}

locals {
  vpc_id_filter_result = [
    for v in data.huaweicloud_vpc_peering_connections.filter_by_vpc_id.connections[*].request_vpc_info[0].vpc_id : v == local.vpc_id
  ]
}

output "is_vpc_id_filter_useful" {
  value = length(local.vpc_id_filter_result) > 0 && alltrue(local.vpc_id_filter_result)
}
`, testAccDataPeeringConnections_basic_base(name), acceptance.HW_REGION_NAME)
}
