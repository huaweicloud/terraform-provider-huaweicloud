package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDbDrRecordDelete_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDbDrRecordDelete_basic(),
			},
		},
	})
}

func testAccGaussDbDrRecordDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_dr_record_delete" "test" {
  task_id = "%[1]s"
}
`, acceptance.HW_GAUSSDB_JOB_ID)
}
