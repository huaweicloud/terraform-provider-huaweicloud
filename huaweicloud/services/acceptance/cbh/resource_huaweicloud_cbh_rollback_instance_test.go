package cbh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRollbackInstance_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCbhServerID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testRollbackInstance_basic(),
			},
		},
	})
}

func testRollbackInstance_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cbh_rollback_instance" "test" {
  server_id = "%s"
}
`, acceptance.HW_CBH_SERVER_ID)
}
