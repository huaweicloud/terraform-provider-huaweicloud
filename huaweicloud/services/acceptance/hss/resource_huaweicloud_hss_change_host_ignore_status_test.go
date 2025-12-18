package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccChangeHostIgnoreStatus_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test requires preparing a host under the default enterprise project that
			// has not enabled host protection.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testChangeHostIgnoreStatus_basic(),
			},
		},
	})
}

func testChangeHostIgnoreStatus_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_change_host_ignore_status" "test" {
  operate_type          = "ignore"
  host_id_list          = ["%s"]
  enterprise_project_id = "0"
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
