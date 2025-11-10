package cbh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResetLoginMode_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCbhServerID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testResetLoginMode_basic(),
			},
		},
	})
}

func testResetLoginMode_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cbh_reset_login_mode" "test" {
  server_id = "%s"
}
`, acceptance.HW_CBH_SERVER_ID)
}
