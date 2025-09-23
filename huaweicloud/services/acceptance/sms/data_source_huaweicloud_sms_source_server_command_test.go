package sms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSmsSourceServerCommand_basic(t *testing.T) {
	dataSource := "data.huaweicloud_sms_source_server_command.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSms(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceSmsSourceServerCommand_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "command_name"),
				),
			},
		},
	})
}

func testDataSourceDataSourceSmsSourceServerCommand_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_sms_source_servers" "source" {
  name = "%s"
}

data "huaweicloud_sms_source_server_command" "test" {
  server_id = data.huaweicloud_sms_source_servers.source.servers[0].id
}
`, acceptance.HW_SMS_SOURCE_SERVER)
}
