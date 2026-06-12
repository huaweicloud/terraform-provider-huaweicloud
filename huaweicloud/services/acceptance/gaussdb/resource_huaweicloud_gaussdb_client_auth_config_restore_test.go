package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceClientAuthConfigRestore_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
			acceptance.TestAccPreCheckGaussDBHbaHistoryId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccClientAuthConfigRestore_basic(),
			},
		},
	})
}

func testAccClientAuthConfigRestore_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_client_auth_config_restore" "test" {
  instance_id    = "%[1]s"
  hba_history_id = "%[2]s"
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID, acceptance.HW_GAUSSDB_HBA_HISTORY_ID)
}
