package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDbDrRelationshipReestablish_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDbDrRelationshipReestablish_basic(),
			},
		},
	})
}

func testAccGaussDbDrRelationshipReestablish_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_dr_relationship_reestablish" "test" {
  instance_id   = "%[1]s"
  disaster_type = "stream"
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
