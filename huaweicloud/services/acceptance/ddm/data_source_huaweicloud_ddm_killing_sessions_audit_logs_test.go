package ddm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDdmKillingSessionsAuditLogs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ddm_killing_sessions_audit_logs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDMInstanceID(t)
			acceptance.TestAccPreCheckDDMTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDdmKillingSessionsAuditLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "process_audit_logs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "process_audit_logs.0.process_id"),
					resource.TestCheckResourceAttrSet(dataSource, "process_audit_logs.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "process_audit_logs.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "process_audit_logs.0.execute_time"),
					resource.TestCheckResourceAttrSet(dataSource, "process_audit_logs.0.execute_user_name"),
				),
			},
		},
	})
}

func testDataSourceDdmKillingSessionsAuditLogs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_ddm_killing_sessions_audit_logs" "test" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
}
`, acceptance.HW_DDM_INSTANCE_ID, acceptance.HW_DDM_START_TIME, acceptance.HW_DDM_END_TIME)
}
