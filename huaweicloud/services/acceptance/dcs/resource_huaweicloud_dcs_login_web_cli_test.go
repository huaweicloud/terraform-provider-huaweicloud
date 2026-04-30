package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDcsLoginWebCli_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsLoginWebCli_basic(name),
			},
		},
	})
}

func testAccDcsLoginWebCli_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dcs_login_web_cli" "test" {
  instance_id = huaweicloud_dcs_instance.instance_1.id
  password    = "Huawei_test"
}
`, testAccDcsInstance_base(name))
}
