package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEventUnblockIp_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID with host protection enabled,
			// and the host is under the default enterprise project.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testEventUnblockIp_basic(),
			},
		},
	})
}

func testEventUnblockIp_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_event_unblock_ip" "test" {
  data_list {
    host_id    = "%[1]s"
    src_ip     = "127.0.0.1"
    login_type = "mysql"
  }

  data_list {
    host_id    = "%[1]s"
    src_ip     = "127.0.0.2"
    login_type = "mysql"
  }

  enterprise_project_id = "0"
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
