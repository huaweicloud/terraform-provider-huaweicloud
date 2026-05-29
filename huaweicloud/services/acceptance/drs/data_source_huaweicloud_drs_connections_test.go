package drs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsConnections_basic(t *testing.T) {
	var (
		datasourceName = "data.huaweicloud_drs_connections.test"
		filterByName   = "data.huaweicloud_drs_connections.filter_by_name"
		filterByType   = "data.huaweicloud_drs_connections.filter_by_db_type"
		nonExistName   = "data.huaweicloud_drs_connections.non_exist"
		dc             = acceptance.InitDataSourceCheck(datasourceName)
		dcFilterName   = acceptance.InitDataSourceCheck(filterByName)
		dcFilter       = acceptance.InitDataSourceCheck(filterByType)
		dcNonExist     = acceptance.InitDataSourceCheck(nonExistName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Note: Please ensure that test data (DRS instances) is prepared in the test environment.
			acceptance.TestAccPreCheckDRSEnableFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDrsConnections_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					dcFilter.CheckResourceExists(),
					dcNonExist.CheckResourceExists(),
					dcFilterName.CheckResourceExists(),
					resource.TestCheckOutput("filter_by_type_is_valid", "true"),
					resource.TestCheckOutput("filter_by_name_is_valid", "true"),
					resource.TestCheckOutput("non_exist_is_empty", "true"),
					resource.TestCheckResourceAttr(filterByType, "connections.0.db_type", "oracle"),
					resource.TestCheckResourceAttrSet(datasourceName, "connections.#"),
					resource.TestCheckResourceAttrSet(filterByType, "connections.0.connection_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "connections.0.name"),
					resource.TestCheckResourceAttrSet(datasourceName, "connections.0.db_type"),
					resource.TestCheckResourceAttrSet(datasourceName, "connections.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "connections.0.create_time"),
					resource.TestCheckResourceAttrSet(datasourceName, "connections.0.connection_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "connections.0.endpoint.0.endpoint_name"),
					resource.TestCheckResourceAttrSet(datasourceName, "connections.0.endpoint.0.db_user"),
					resource.TestCheckResourceAttrSet(datasourceName, "connections.0.endpoint.0.db_port"),
					resource.TestCheckResourceAttrSet(datasourceName, "connections.0.ssl.0.ssl_link"),
				),
			},
		},
	})
}

const testAccDataSourceDrsConnections_basic = `
data "huaweicloud_drs_connections" "test" {}

data "huaweicloud_drs_connections" "filter_by_db_type" {
  db_type = data.huaweicloud_drs_connections.test.connections[0].db_type
}

data "huaweicloud_drs_connections" "filter_by_name" {
  name = data.huaweicloud_drs_connections.test.connections[0].name
}

data "huaweicloud_drs_connections" "non_exist" {
  name = "non_existent_conn_test"
}

output "filter_by_type_is_valid" {
  value = length(data.huaweicloud_drs_connections.filter_by_db_type.connections) >= 1
}

output "filter_by_name_is_valid" {
  value = length(data.huaweicloud_drs_connections.filter_by_name.connections) >= 1
}

output "non_exist_is_empty" {
  value = length(data.huaweicloud_drs_connections.non_exist.connections) == 0
}
`
