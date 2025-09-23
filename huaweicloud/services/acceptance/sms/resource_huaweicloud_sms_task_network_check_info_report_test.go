package sms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceSmsTaskNetworkCheckInfoReport_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSms(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testSmsTaskNetworkCheckInfoReport_basic(name),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testSmsTaskNetworkCheckInfoReport_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_sms_task_network_check_info_report" "test" {
  task_id                  = huaweicloud_sms_task.migration.id
  network_delay            = 20.0
  network_jitter           = 2.0
  migration_speed          = 100.0
  loss_percentage          = 0.0
  cpu_usage                = 20.0
  mem_usage                = 20.0
  evaluation_result        = "success"
  domain_connectivity      = true
  destination_connectivity = true
}
`, testAccMigrationTask_basic(name))
}
