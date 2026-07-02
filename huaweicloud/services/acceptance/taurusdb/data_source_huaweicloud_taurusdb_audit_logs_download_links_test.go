package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBAuditLogsDownloadLinks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_audit_logs_download_links.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTaurusDBAuditLogsDownloadLinks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "links.#"),
					resource.TestCheckResourceAttrSet(dataSource, "links.0"),
				),
			},
		},
	})
}

func testAccDataSourceTaurusDBAuditLogsDownloadLinks_basic() string {
	return fmt.Sprintf(`
%s

locals {
  ids = data.huaweicloud_taurusdb_audit_logs.test.audit_logs.*.id
}

data "huaweicloud_taurusdb_audit_logs_download_links" "test" {
  instance_id = "%s"
  ids         = local.ids
}
`, testDataSourceTaurusDBAuditLogs_basic(), acceptance.HW_TAURUSDB_INSTANCE_ID)
}
