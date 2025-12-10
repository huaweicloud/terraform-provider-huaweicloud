package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCollectorConnections_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_collector_connections.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running test, prepare a collector channel group.
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCollectorConnections_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.channel_refer_count"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.connection_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.connection_type"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.info"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.module_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.template_title"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.title"),

					resource.TestCheckOutput("is_connection_type_filter_useful", "true"),
					resource.TestCheckOutput("is_title_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceCollectorConnections_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_collector_connections" "test" {
  workspace_id = "%[1]s"
}

# Filter by connection_type
locals {
  connection_type = data.huaweicloud_secmaster_collector_connections.test.records[0].connection_type
}

data "huaweicloud_secmaster_collector_connections" "filter_by_connection_type" {
  workspace_id    = "%[1]s"
  connection_type = local.connection_type
}

output "is_connection_type_filter_useful" {
  value = length(data.huaweicloud_secmaster_collector_connections.filter_by_connection_type.records) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_collector_connections.filter_by_connection_type.records[*].connection_type :
    v == local.connection_type]
  )
}

# Filter by title
locals {
  title = data.huaweicloud_secmaster_collector_connections.test.records[0].title
}

data "huaweicloud_secmaster_collector_connections" "filter_by_title" {
  workspace_id = "%[1]s"
  title        = local.title
}

output "is_title_filter_useful" {
  value = length(data.huaweicloud_secmaster_collector_connections.filter_by_title.records) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_collector_connections.filter_by_title.records[*].title : v == local.title]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
