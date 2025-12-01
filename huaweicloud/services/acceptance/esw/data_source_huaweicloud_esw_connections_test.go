package esw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEswConnections_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_esw_connections.test"
	dc := acceptance.InitDataSourceCheck(rName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceEswConnections_basic(name),
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
					resource.TestCheckOutput("connection_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceEswConnections_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_esw_connections" "test" {
  depends_on = [huaweicloud_esw_connection.test]

  instance_id = huaweicloud_esw_instance.test.id
}

locals{
  connection_id = huaweicloud_esw_connection.test.id
}
data "huaweicloud_esw_connections" "connection_id_filter" {
  depends_on = [huaweicloud_esw_connection.test]

  instance_id   = huaweicloud_esw_instance.test.id
  connection_id = huaweicloud_esw_connection.test.id
}
output "connection_id_filter_is_useful" {
  value = length(data.huaweicloud_esw_connections.connection_id_filter.connections) > 0 && alltrue(
    [for v in data.huaweicloud_esw_connections.connection_id_filter.connections[*].id : v == local.connection_id]
  )
}

locals{
  name = huaweicloud_esw_connection.test.name
}
data "huaweicloud_esw_connections" "name_filter" {
  depends_on = [huaweicloud_esw_connection.test]

  instance_id = huaweicloud_esw_instance.test.id
  name        = huaweicloud_esw_connection.test.name
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_esw_connections.name_filter.connections) > 0 && alltrue(
    [for v in data.huaweicloud_esw_connections.name_filter.connections[*].name : v == local.name]
  )
}
`, testAccEswConnection_basic(name))
}
