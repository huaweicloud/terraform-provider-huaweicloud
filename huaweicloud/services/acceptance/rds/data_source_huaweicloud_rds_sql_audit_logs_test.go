package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsSqlAuditLogs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_sql_audit_logs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
			acceptance.TestAccPreCheckRdsTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsSqlAuditLogs_basic(),
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

func testDataSourceRdsSqlAuditLogs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_sql_audit_logs" "test" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
}
`, acceptance.HW_RDS_INSTANCE_ID, acceptance.HW_RDS_START_TIME, acceptance.HW_RDS_END_TIME)
}
