package esw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEswAllConnections_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_esw_all_connections.test"
	dc := acceptance.InitDataSourceCheck(rName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceEswAllConnections_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "connections.#"),
					resource.TestCheckResourceAttrSet(rName, "connections.0.id"),
					resource.TestCheckResourceAttrSet(rName, "connections.0.instance_id"),
					resource.TestCheckResourceAttrSet(rName, "connections.0.name"),
					resource.TestCheckResourceAttrSet(rName, "connections.0.project_id"),
					resource.TestCheckResourceAttrSet(rName, "connections.0.fixed_ips.#"),
					resource.TestCheckResourceAttrSet(rName, "connections.0.remote_infos.#"),
					resource.TestCheckResourceAttrSet(rName, "connections.0.remote_infos.0.segmentation_id"),
					resource.TestCheckResourceAttrSet(rName, "connections.0.remote_infos.0.tunnel_ip"),
					resource.TestCheckResourceAttrSet(rName, "connections.0.remote_infos.0.tunnel_port"),
					resource.TestCheckResourceAttrSet(rName, "connections.0.remote_infos.0.tunnel_type"),
					resource.TestCheckResourceAttrSet(rName, "connections.0.vpc_id"),
					resource.TestCheckResourceAttrSet(rName, "connections.0.virsubnet_id"),
					resource.TestCheckResourceAttrSet(rName, "connections.0.status"),
					resource.TestCheckResourceAttrSet(rName, "connections.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "connections.0.updated_at"),
					resource.TestCheckOutput("instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("connection_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("vpc_id_filter_is_useful", "true"),
					resource.TestCheckOutput("virsubnet_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceEswAllConnections_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_esw_all_connections" "test" {
  depends_on = [huaweicloud_esw_connection.test]
}

locals{
  instance_id = huaweicloud_esw_instance.test.id
}
data "huaweicloud_esw_all_connections" "instance_id_filter" {
  depends_on = [huaweicloud_esw_connection.test]

  instance_id = huaweicloud_esw_instance.test.id
}
output "instance_id_filter_is_useful" {
  value = length(data.huaweicloud_esw_all_connections.instance_id_filter.connections) > 0 && alltrue(
    [for v in data.huaweicloud_esw_all_connections.instance_id_filter.connections[*].instance_id : v == local.instance_id]
  )
}

locals{
  connection_id = huaweicloud_esw_connection.test.id
}
data "huaweicloud_esw_all_connections" "connection_id_filter" {
  depends_on = [huaweicloud_esw_connection.test]

  connection_id = huaweicloud_esw_connection.test.id
}
output "connection_id_filter_is_useful" {
  value = length(data.huaweicloud_esw_all_connections.connection_id_filter.connections) > 0 && alltrue(
    [for v in data.huaweicloud_esw_all_connections.connection_id_filter.connections[*].id : v == local.connection_id]
  )
}

locals{
  name = huaweicloud_esw_connection.test.name
}
data "huaweicloud_esw_all_connections" "name_filter" {
  depends_on = [huaweicloud_esw_connection.test]

  name = huaweicloud_esw_connection.test.name
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_esw_all_connections.name_filter.connections) > 0 && alltrue(
    [for v in data.huaweicloud_esw_all_connections.name_filter.connections[*].name : v == local.name]
  )
}

locals{
  vpc_id = huaweicloud_esw_connection.test.vpc_id
}
data "huaweicloud_esw_all_connections" "vpc_id_filter" {
  depends_on = [huaweicloud_esw_connection.test]

  vpc_id = huaweicloud_esw_connection.test.vpc_id
}
output "vpc_id_filter_is_useful" {
  value = length(data.huaweicloud_esw_all_connections.vpc_id_filter.connections) > 0 && alltrue(
    [for v in data.huaweicloud_esw_all_connections.vpc_id_filter.connections[*].vpc_id : v == local.vpc_id]
  )
}

locals{
  virsubnet_id = huaweicloud_esw_connection.test.virsubnet_id
}
data "huaweicloud_esw_all_connections" "virsubnet_id_filter" {
  depends_on = [huaweicloud_esw_connection.test]

  virsubnet_id = huaweicloud_esw_connection.test.virsubnet_id
}
output "virsubnet_id_filter_is_useful" {
  value = length(data.huaweicloud_esw_all_connections.virsubnet_id_filter.connections) > 0 && alltrue(
    [for v in data.huaweicloud_esw_all_connections.virsubnet_id_filter.connections[*].virsubnet_id : v == local.virsubnet_id]
  )
}
`, testAccEswConnection_basic(name))
}
