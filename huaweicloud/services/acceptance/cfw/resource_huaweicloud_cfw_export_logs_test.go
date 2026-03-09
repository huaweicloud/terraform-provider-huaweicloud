package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccExportLogs_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
			acceptance.TestAccPreCheckCfwTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testExportLogs_basic(),
			},
		},
	})
}

func testExportLogs_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cfw_export_logs" "test" {
  fw_instance_id = "%[1]s"
  start_time     = %[2]s
  end_time       = %[3]s
  log_type       = "internet"
  type           = "flow"
  time_zone      = "GMT+08:00"

  filters {
    field    = "src_ip"
    operator = "contain"
    values   = ["ANY", "TCP"]
  }

  export_file_name = "cfw-tf-test-log"
}
`, acceptance.HW_CFW_INSTANCE_ID, acceptance.HW_CFW_START_TIME, acceptance.HW_CFW_END_TIME)
}
