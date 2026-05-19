package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceNodeSessionStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_geminidb_node_session_statistics.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGeminiDBNodeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNodeSessionStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "top_source_ips.#"),
					resource.TestCheckResourceAttrSet(dataSource, "top_source_ips.0.client_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "top_source_ips.0.connection_count"),
					resource.TestCheckResourceAttrSet(dataSource, "top_dbs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "top_dbs.0.db"),
					resource.TestCheckResourceAttrSet(dataSource, "top_dbs.0.connection_count"),
					resource.TestCheckResourceAttrSet(dataSource, "total_connection_count"),
					resource.TestCheckResourceAttrSet(dataSource, "active_connection_count"),
				),
			},
		},
	})
}

func testAccDataSourceNodeSessionStatistics_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_geminidb_node_session_statistics" "test" {
  node_id = "%[1]s"
}
`, acceptance.HW_GEMINIDB_NODE_ID)
}
