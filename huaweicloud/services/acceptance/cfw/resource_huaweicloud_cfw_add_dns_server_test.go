package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAddDNSServer_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
			acceptance.TestAccPreCheckCfwServerIp(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAddDNSServer_basic(),
			},
		},
	})
}

func testAddDNSServer_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cfw_add_dns_server" "test" {
  fw_instance_id = "%[1]s"
  server_ip      = "%[2]s"
}
`, acceptance.HW_CFW_INSTANCE_ID, acceptance.HW_CFW_SERVER_IP)
}
