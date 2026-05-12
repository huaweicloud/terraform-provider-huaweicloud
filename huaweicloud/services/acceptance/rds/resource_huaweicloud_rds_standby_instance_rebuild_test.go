package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRdsStandbyInstanceRebuild_basic(t *testing.T) {
	rName := "huaweicloud_rds_standby_instance_rebuild.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsStandbyInstanceRebuild_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "instance_id", acceptance.HW_RDS_INSTANCE_ID),
					resource.TestCheckResourceAttrSet(rName, "workflow_id"),
					resource.TestCheckResourceAttrSet(rName, "last_rebuild_time"),
					resource.TestCheckResourceAttrSet(rName, "next_rebuild_time"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRdsStandbyInstanceRebuild_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_standby_instance_rebuild" "test" {
  instance_id = "%s"
}
`, acceptance.HW_RDS_INSTANCE_ID)
}
