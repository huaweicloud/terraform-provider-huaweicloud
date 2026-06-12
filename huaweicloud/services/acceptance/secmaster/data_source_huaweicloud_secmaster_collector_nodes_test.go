package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCollectorNodes_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_collector_nodes.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCollectorNodes_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.channel_instance_refer_count"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.monitor.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.monitor.0.cpu_idle"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.monitor.0.cpu_usage"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.node_name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.ip_address"),

					resource.TestCheckOutput("is_node_id_filter_useful", "true"),
					resource.TestCheckOutput("is_node_name_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testDataSourceCollectorNodes_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_collector_nodes" "test" {
  workspace_id  = "%[1]s"
  health_status = "NORMAL"
  sort_key      = "node_id"
  sort_dir      = "desc"
}

# Filter using node_id.
locals {
  node_id = data.huaweicloud_secmaster_collector_nodes.test.records[0].node_id
}

data "huaweicloud_secmaster_collector_nodes" "node_id_filter" {
  workspace_id = "%[1]s"
  node_id      = local.node_id
}

output "is_node_id_filter_useful" {
  value = length(data.huaweicloud_secmaster_collector_nodes.node_id_filter.records) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_collector_nodes.node_id_filter.records[*].node_id : v == local.node_id]
  )
}

# Filter using node_name.
locals {
  node_name = data.huaweicloud_secmaster_collector_nodes.test.records[0].node_name
}

data "huaweicloud_secmaster_collector_nodes" "node_name_filter" {
  workspace_id = "%[1]s"
  node_name    = local.node_name
}

output "is_node_name_filter_useful" {
  value = length(data.huaweicloud_secmaster_collector_nodes.node_name_filter.records) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_collector_nodes.node_name_filter.records[*].node_name : v == local.node_name]
  )
}

# Filter using non existent name.
data "huaweicloud_secmaster_collector_nodes" "not_found" {
  workspace_id = "%[1]s"
  node_name    = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_secmaster_collector_nodes.not_found.records) == 0
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
