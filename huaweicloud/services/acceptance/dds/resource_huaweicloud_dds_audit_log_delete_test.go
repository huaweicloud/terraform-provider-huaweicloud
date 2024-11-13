package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDDSAuditLogDelete_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSInstanceID(t)
			acceptance.TestAccPreCheckDDSTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSAuditLogDelete_basic(),
			},
		},
	})
}

func testAccDDSAuditLogDelete_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dds_audit_logs" "test" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
}

resource "huaweicloud_dds_audit_log_delete" "test" {
  instance_id = "%[1]s"
  file_names  = [try(data.huaweicloud_dds_audit_logs.test.audit_logs[0].name, "")]

  lifecycle {
    ignore_changes = [
      file_names,
    ]
  }
}`, acceptance.HW_DDS_INSTANCE_ID, acceptance.HW_DDS_START_TIME, acceptance.HW_DDS_END_TIME)
}
