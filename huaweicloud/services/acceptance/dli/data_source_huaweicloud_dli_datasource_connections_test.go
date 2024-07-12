package dli

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDatasourceConnections_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceName()
		dataSourceName = "data.huaweicloud_dli_datasource_connections.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byName   = "data.huaweicloud_dli_datasource_connections.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNameNotFound   = "data.huaweicloud_dli_datasource_connections.filter_by_name_not_found"
		dcbyNameNotFound = acceptance.InitDataSourceCheck(byNameNotFound)

		byTags   = "data.huaweicloud_dli_datasource_connections.filter_by_tags"
		dcByTags = acceptance.InitDataSourceCheck(byTags)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDatasourceConnections_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "connections.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "connections.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "connections.0.routes.#"),
					resource.TestCheckResourceAttr(dataSourceName, "connections.0.routes.0.name", name),
					resource.TestCheckResourceAttr(dataSourceName, "connections.0.routes.0.cidr", "10.169.0.0/16"),
					resource.TestCheckResourceAttrSet(dataSourceName, "connections.0.queues.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "connections.0.hosts.#"),
					resource.TestCheckResourceAttr(dataSourceName, "connections.0.hosts.0.ip", "172.0.0.2"),
					resource.TestCheckResourceAttr(dataSourceName, "connections.0.hosts.0.name", name),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),

					dcbyNameNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful_not_found", "true"),

					dcByTags.CheckResourceExists(),
					resource.TestCheckResourceAttr(byTags, "connections.#", "1"),
				),
			},
		},
	})
}

func testDatasourceConnections_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  tags = {
    foo  = "bar"
    foo1 = "bar1"
  }
}

resource "huaweicloud_dli_datasource_connection" "test" {
  name      = "%[2]s"
  vpc_id    = huaweicloud_vpc.test.id
  subnet_id = huaweicloud_vpc_subnet.test.id

  routes {
    cidr = "10.169.0.0/16"
    name = "%[2]s"
  }

  hosts {
    ip   = "172.0.0.2"
    name = "%[2]s"
  }

  tags = local.tags
}

data "huaweicloud_dli_datasource_connections" "test" {
  depends_on = [
    huaweicloud_dli_datasource_connection.test
  ] 
}

locals {
  name = huaweicloud_dli_datasource_connection.test.name
}

data "huaweicloud_dli_datasource_connections" "filter_by_name" {
  depends_on = [
    huaweicloud_dli_datasource_connection.test
  ] 
  name = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dli_datasource_connections.filter_by_name.connections[*].name : v == local.name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) ==1 && alltrue(local.name_filter_result)
}

data "huaweicloud_dli_datasource_connections" "filter_by_name_not_found" {
  depends_on = [
    huaweicloud_dli_datasource_connection.test
  ] 
  name = "tf_test"
}

output "is_name_filter_useful_not_found" {
  value = length(data.huaweicloud_dli_datasource_connections.filter_by_name_not_found.connections) == 0
}

data "huaweicloud_dli_datasource_connections" "filter_by_tags" {
  depends_on = [
    huaweicloud_dli_datasource_connection.test
  ]
  tags = local.tags
}
`, common.TestVpc(name), name)
}
