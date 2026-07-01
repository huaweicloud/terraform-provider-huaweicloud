package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBSlowLogsDownloadLinks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_slow_logs_download_links.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
			acceptance.TestAccPreCheckTaurusDBNodeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTaurusDBSlowLogsDownloadLinks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.file_name"),
					resource.TestCheckResourceAttr(dataSource, "list.0.status", "SUCCESS"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.file_size"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.file_link"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.create_at"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.updated_at"),
				),
			},
		},
	})
}

func testAccDataSourceTaurusDBSlowLogsDownloadLinks_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_taurusdb_slow_logs_download_links" "test" {
  instance_id = "%[1]s"
  node_id     = "%[2]s"
}
`, acceptance.HW_TAURUSDB_INSTANCE_ID, acceptance.HW_TAURUSDB_NODE_ID)
}
