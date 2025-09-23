package sms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceSmsSourceServerCommandResultReport_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSms(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testSmsSourceServerCommandResultReport_basic(),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testSmsSourceServerCommandResultReport_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_sms_source_servers" "source" {
  name = "%s"
}

resource "huaweicloud_sms_source_server_command_result_report" "test" {
  server_id     = data.huaweicloud_sms_source_servers.source.servers[0].id
  command_name  = "START"
  result        = "success"
  result_detail = jsonencode({
    "msg": "test"
  })
}
`, acceptance.HW_SMS_SOURCE_SERVER)
}
