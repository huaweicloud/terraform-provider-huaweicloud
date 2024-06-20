package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsSqlAuditLogLinks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_sql_audit_log_links.test"
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
				Config: testDataSourceDataSourceRdsSqlAuditLogLinks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "links.#"),
					resource.TestCheckResourceAttr(dataSource, "links.#", "2"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRdsSqlAuditLogLinks_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_sql_audit_log_links" "test" {
  instance_id = "%[2]s"

  ids = [
    data.huaweicloud_rds_sql_audit_logs.test.audit_logs[0].id,
    data.huaweicloud_rds_sql_audit_logs.test.audit_logs[1].id,
  ]
}
`, testDataSourceRdsSqlAuditLogs_basic(), acceptance.HW_RDS_INSTANCE_ID)
}
