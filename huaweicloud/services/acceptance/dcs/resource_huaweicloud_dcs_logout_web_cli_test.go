package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDcsLogoutWebCli_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	loginResourceName := "huaweicloud_dcs_login_web_cli.test"
	logoutResourceName := "huaweicloud_dcs_logout_web_cli.test"

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
				Config: testAccDcsLogoutWebCli_basic(name),
				Check:  resource.TestCheckResourceAttrSet(logoutResourceName, "id"),
			},
		},
	})
}

func testAccDcsLogoutWebCli_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dcs_login_web_cli" "test" {
  instance_id = huaweicloud_dcs_instance.instance_1.id
  password    = "Huawei_test"
}

resource "huaweicloud_dcs_logout_web_cli" "test" {
  instance_id = huaweicloud_dcs_instance.instance_1.id
  client_id   = huaweicloud_dcs_login_web_cli.test.client_id
}
`, testAccDcsInstance_base(name))
}
