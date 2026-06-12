package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDbDrInstanceToPrimary_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBDrInstanceToPrimary_basic(),
			},
		},
	})
}

func testAccGaussDBDrInstanceToPrimary_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_dr_instance_to_primary" "test" {
  instance_id        = "%s"
  disaster_type      = "stream"
  is_support_restore = true
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
