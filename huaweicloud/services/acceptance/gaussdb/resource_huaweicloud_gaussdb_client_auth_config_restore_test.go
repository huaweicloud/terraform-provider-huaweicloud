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
data "huaweicloud_gaussdb_client_auth_config_history" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_gaussdb_client_auth_config_restore" "test" {
  instance_id    = "%[1]s"
  hba_history_id = data.huaweicloud_gaussdb_client_auth_config_history.test.hba_histories.0.id
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
