package taurusdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBAuditLogs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_audit_logs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTaurusDBAuditLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.0.begin_time"),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.0.end_time"),
				),
			},
		},
	})
}

func testDataSourceTaurusDBAuditLogs_basic() string {
	cst := time.FixedZone("CST", 8*3600)
	endTime := time.Now().In(cst).Format("2006-01-02T15:04:05-0700")
	startTime := time.Now().In(cst).Add(-24 * time.Hour).Format("2006-01-02T15:04:05-0700")
	return fmt.Sprintf(`
data "huaweicloud_taurusdb_audit_logs" "test" {
  instance_id = "%s"
  start_time  = "%s"
  end_time    = "%s"
}
`, acceptance.HW_TAURUSDB_INSTANCE_ID, startTime, endTime)
}
