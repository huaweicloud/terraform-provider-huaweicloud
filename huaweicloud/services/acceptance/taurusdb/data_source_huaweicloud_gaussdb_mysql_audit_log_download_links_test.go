package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDBMysqlAuditLogDownloadLinks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_mysql_audit_log_download_links.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBMysqlInstanceId(t)
			acceptance.TestAccPreCheckGaussDBMysqlNodeId(t)
			acceptance.TestAccPreCheckGaussDBMysqlTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussDBMysqlAuditLogDownloadLinks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "links.#"),
					resource.TestCheckResourceAttrSet(dataSource, "links.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "links.0.full_name"),
					resource.TestCheckResourceAttrSet(dataSource, "links.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "links.0.updated_time"),
					resource.TestCheckResourceAttrSet(dataSource, "links.0.download_link"),
					resource.TestCheckResourceAttrSet(dataSource, "links.0.link_expired_time"),

					resource.TestCheckOutput("node_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussDBMysqlAuditLogDownloadLinks_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_mysql_audit_log_download_links" "test" {
  instance_id = "%[1]s"
  start_time  = "%[3]s"
  end_time    = "%[4]s"
}

locals {
  node_id = "%[2]s"
}
data "huaweicloud_gaussdb_mysql_audit_log_download_links" "node_id_filter" {
  instance_id  = "%[1]s"
  node_id      = "%[2]s"
  start_time   = "%[3]s"
  end_time     = "%[4]s"
}
output "node_id_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_audit_log_download_links.node_id_filter.links) > 0
}
`, acceptance.HW_GAUSSDB_MYSQL_INSTANCE_ID, acceptance.HW_GAUSSDB_MYSQL_NODE_ID, acceptance.HW_GAUSSDB_MYSQL_START_TIME,
		acceptance.HW_GAUSSDB_MYSQL_END_TIME)
}
