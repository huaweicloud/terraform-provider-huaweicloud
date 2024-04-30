package dataarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDataConnections_basic(t *testing.T) {
	var (
		name             = acceptance.RandomAccResourceName()
		byConnectionId   = "data.huaweicloud_dataarts_studio_data_connections.by_connection_id"
		dcByConnectionId = acceptance.InitDataSourceCheck(byConnectionId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDataConnections_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dcByConnectionId.CheckResourceExists(),
					resource.TestCheckResourceAttr(byConnectionId, "connections.#", "1"),
					resource.TestCheckResourceAttrSet(byConnectionId, "connections.0.qualified_name"),
					resource.TestCheckResourceAttrSet(byConnectionId, "connections.0.created_by"),
					resource.TestCheckResourceAttrSet(byConnectionId, "connections.0.created_at"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceDataConnections_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dataarts_studio_data_connections" "by_connection_id" {
  workspace_id  = "%[2]s"
  connection_id = huaweicloud_dataarts_studio_data_connection.test.id
}

data "huaweicloud_dataarts_studio_data_connections" "by_name_filter" {
  depends_on   = [huaweicloud_dataarts_studio_data_connection.test]

  workspace_id = "%[2]s"
  name         = "%[3]s"
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_dataarts_studio_data_connections.by_name_filter.connections) > 0 && alltrue(
    [for v in data.huaweicloud_dataarts_studio_data_connections.by_name_filter.connections[*].name : 
    v == "%[3]s"]
  )  
}

data "huaweicloud_dataarts_studio_data_connections" "by_type_filter" {
  depends_on   = [huaweicloud_dataarts_studio_data_connection.test]

  workspace_id = "%[2]s"
  type         = "DLI"
}

output "type_filter_is_useful" {
  value = length(data.huaweicloud_dataarts_studio_data_connections.by_type_filter.connections) > 0 && alltrue(
    [for v in data.huaweicloud_dataarts_studio_data_connections.by_type_filter.connections[*].type : 
    v == "DLI"]
  )  
}
`, testAccDataConnection_basic(name), acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}
