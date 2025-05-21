package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRdsDrInstanceToPrimarySwitch_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsDrInstanceToPrimary_basic(),
			},
		},
	})
}

func testAccRdsDrInstanceToPrimary_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_dr_instance_to_primary" "test" {
  instance_id = "%[1]s"
}
`, acceptance.HW_RDS_INSTANCE_ID)
}
