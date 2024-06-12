package dc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test case, you need to prepare a DC connection.
// DC connection only supports created by console
func TestAccDataSourceConnections_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_dc_connections.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byConnectionId   = "data.huaweicloud_dc_connections.filter_by_connection_id"
		dcByConnectionId = acceptance.InitDataSourceCheck(byConnectionId)

		byName   = "data.huaweicloud_dc_connections.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byType   = "data.huaweicloud_dc_connections.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byStatus = "data.huaweicloud_dc_connections.filter_by_status"
		dcStatus = acceptance.InitDataSourceCheck(byStatus)

		byPortType   = "data.huaweicloud_dc_connections.filter_by_port_type"
		dcByPortType = acceptance.InitDataSourceCheck(byPortType)

		byEnterpriseProjectId   = "data.huaweicloud_dc_connections.filter_by_eps"
		dcByEnterpriseProjectId = acceptance.InitDataSourceCheck(byEnterpriseProjectId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDcFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceConnections_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "direct_connects.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "direct_connects.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "direct_connects.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "direct_connects.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "direct_connects.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "direct_connects.0.port_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "direct_connects.0.bandwidth"),
					resource.TestCheckResourceAttrSet(dataSourceName, "direct_connects.0.location"),

					dcByConnectionId.CheckResourceExists(),
					resource.TestCheckOutput("connection_id_filter_useful", "true"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_useful", "true"),

					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("type_filter_useful", "true"),

					dcStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_useful", "true"),

					dcByPortType.CheckResourceExists(),
					resource.TestCheckOutput("port_type_filter_useful", "true"),

					dcByEnterpriseProjectId.CheckResourceExists(),
					resource.TestCheckOutput("eps_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceConnections_basic() string {
	return (`
data "huaweicloud_dc_connections" "test" {}

locals {
  connection_id = data.huaweicloud_dc_connections.test.direct_connects[0].id
}

data "huaweicloud_dc_connections" "filter_by_connection_id" {
  connection_id = local.connection_id
}

output "connection_id_filter_useful" {
  value = length(data.huaweicloud_dc_connections.filter_by_connection_id.direct_connects) == 1 && alltrue(
    [for v in data.huaweicloud_dc_connections.filter_by_connection_id.direct_connects[*].id : v == local.connection_id]
  )
}

locals {
  name = data.huaweicloud_dc_connections.test.direct_connects[0].name
}

data "huaweicloud_dc_connections" "filter_by_name" {
  name = local.name
}

output "name_filter_useful" {
  value = length(data.huaweicloud_dc_connections.filter_by_name.direct_connects) > 0 && alltrue(
    [for v in data.huaweicloud_dc_connections.filter_by_name.direct_connects[*].name : v == local.name]
  )
}

locals {
  type = data.huaweicloud_dc_connections.test.direct_connects[0].type
}


data "huaweicloud_dc_connections" "filter_by_type" {
  type = local.type
}

output "type_filter_useful" {
  value = length(data.huaweicloud_dc_connections.filter_by_type.direct_connects) > 0 && alltrue(
    [for v in data.huaweicloud_dc_connections.filter_by_type.direct_connects[*].type : v == local.type]
  )
}

locals {
  status = data.huaweicloud_dc_connections.test.direct_connects[0].status
}

data "huaweicloud_dc_connections" "filter_by_status" {
  status = local.status
}

output "status_filter_useful" {
  value = length(data.huaweicloud_dc_connections.filter_by_status.direct_connects) > 0 && alltrue(
    [for v in data.huaweicloud_dc_connections.filter_by_status.direct_connects[*].status : v == local.status]
  )
}

locals {
  port_type = data.huaweicloud_dc_connections.test.direct_connects[0].port_type
}

data "huaweicloud_dc_connections" "filter_by_port_type" {
  port_type = local.port_type
}

output "port_type_filter_useful" {
  value = length(data.huaweicloud_dc_connections.filter_by_port_type.direct_connects) > 0 && alltrue(
    [for v in data.huaweicloud_dc_connections.filter_by_port_type.direct_connects[*].port_type : v == local.port_type]
  )
}

locals {
  enterprise_project_id = data.huaweicloud_dc_connections.test.direct_connects[0].enterprise_project_id
}

data "huaweicloud_dc_connections" "filter_by_eps" {
  enterprise_project_id = local.enterprise_project_id
}

output "eps_filter_useful" {
  value = length(data.huaweicloud_dc_connections.filter_by_eps.direct_connects) > 0 && alltrue(
    [for v in data.huaweicloud_dc_connections.filter_by_eps.direct_connects[*].enterprise_project_id : 
      v == local.enterprise_project_id]
  )
}
`)
}
