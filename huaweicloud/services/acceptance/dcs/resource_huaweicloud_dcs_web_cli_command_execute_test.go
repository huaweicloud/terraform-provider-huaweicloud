package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDcsWebCliCommandExecute_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	loginResourceName := "huaweicloud_dcs_login_web_cli.test"
	commandExecuteResourceName := "huaweicloud_dcs_web_cli_command_execute.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsLoginWebCli_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(loginResourceName, "client_id"),
				),
			},
			{
				Config: testAccDcsWebCliCommandExecute_basic(name),
				Check:  resource.TestCheckResourceAttrSet(commandExecuteResourceName, "id"),
			},
		},
	})
}

func testAccDcsWebCliCommandExecute_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dcs_login_web_cli" "test" {
  instance_id = huaweicloud_dcs_instance.instance_1.id
  password    = "Huawei_test"
}

resource "huaweicloud_dcs_web_cli_command_execute" "test" {
  instance_id = huaweicloud_dcs_instance.instance_1.id
  client_id   = huaweicloud_dcs_login_web_cli.test.client_id
  command     = "scan 0"
  database    = 0
}
`, testAccDcsInstance_base(name))
}
