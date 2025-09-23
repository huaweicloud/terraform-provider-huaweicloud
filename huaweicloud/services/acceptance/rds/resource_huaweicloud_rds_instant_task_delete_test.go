package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRdsInstantTaskDelete_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstantJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstantTaskDelete_basic(),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccRdsInstantTaskDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_instant_task_delete" "test" {
  job_id = "%[1]s"
}`, acceptance.HW_RDS_INSTANCE_ID)
}
