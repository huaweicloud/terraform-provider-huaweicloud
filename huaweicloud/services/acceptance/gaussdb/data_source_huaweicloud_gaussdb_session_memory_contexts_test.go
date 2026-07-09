package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDBSessionMemoryContexts_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_session_memory_contexts.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGaussDBSessionMemoryContexts_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "memory_context_info.#"),
					resource.TestCheckResourceAttrSet(dataSource, "memory_context_info.0.context_name"),
					resource.TestCheckResourceAttrSet(dataSource, "memory_context_info.0.amount"),
					resource.TestCheckResourceAttrSet(dataSource, "memory_context_info.0.size"),
				),
			},
		},
	})
}

func testAccDataSourceGaussDBSessionMemoryContexts_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_key_view_nodes_deliver" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_gaussdb_instance_real_time_sessions" "test" {
  instance_id  = "%[1]s"
  node_id      = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.node_id
  component_id = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.component_id
}

data "huaweicloud_gaussdb_session_memory_contexts" "test" {
  instance_id = "%[1]s"
  node_id     = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.node_id
  session_id  = data.huaweicloud_gaussdb_instance_real_time_sessions.test.sessions.0.session_id
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
