package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDdsAuditLogs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dds_audit_logs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSInstanceID(t)
			acceptance.TestAccPreCheckDDSTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDdsAuditLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.0.end_time"),

					resource.TestCheckOutput("node_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDdsAuditLogs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dds_audit_logs" "test" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
}

locals {
  node_id  = data.huaweicloud_dds_audit_logs.test.audit_logs.0.node_id
}

data "huaweicloud_dds_audit_logs" "filter" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
  node_id     = local.node_id
}

output "node_id_filter_is_useful" {
  value = length(data.huaweicloud_dds_audit_logs.filter.audit_logs) > 0 && alltrue(
    [for v in data.huaweicloud_dds_audit_logs.filter.audit_logs[*].node_id : v == local.node_id]
  )
}
`, acceptance.HW_DDS_INSTANCE_ID, acceptance.HW_DDS_START_TIME, acceptance.HW_DDS_END_TIME)
}
