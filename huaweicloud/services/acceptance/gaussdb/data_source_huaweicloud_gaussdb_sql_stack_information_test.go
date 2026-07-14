package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbSqlStackInformation_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_sql_stack_information.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbSqlStackInformation_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "gs_stack"),
				),
			},
		},
	})
}

func testDataSourceGaussdbSqlStackInformation_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_key_view_nodes_deliver" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_gaussdb_instance_real_time_sessions" "test" {
  instance_id  = "%[1]s"
  node_id      = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.node_id
  component_id = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.component_id
}

data "huaweicloud_gaussdb_sql_stack_information" "test" {
  instance_id = "%[1]s"
  pid         = data.huaweicloud_gaussdb_instance_real_time_sessions.test.sessions.0.pid
  node_id     = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.node_id
}

data "huaweicloud_gaussdb_sql_stack_information" "comp_id_useful" {
  instance_id = "%[1]s"
  pid         = data.huaweicloud_gaussdb_instance_real_time_sessions.test.sessions.0.pid
  node_id     = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.node_id
  comp_id     = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.component_id
}

output "comp_id_useful" {
  value = length(data.huaweicloud_gaussdb_sql_stack_information.comp_id_useful.gs_stack) > 0
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
